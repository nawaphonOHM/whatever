package detector

import (
	"regexp"
	"slices"
	"strings"
)

var awsSuffixPatterns = []string{
	"amazonaws.com",
	"amazonaws.com.cn",
	"api.aws",
}

const (
	vpceLabel   = "vpce"
	bucketAlias = "bucket"
)

var awsRegionLabel = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-[0-9]+$`)

func isAWSS3Label(label string) bool {
	return label == "s3" || strings.HasPrefix(label, "s3-") || strings.HasPrefix(label, "s3express-")
}

// isAWSTailLabel reports whether label may sit between the S3 service label and the AWS suffix.
func isAWSTailLabel(label string) bool {
	switch label {
	case "dualstack", "fips", vpceLabel:
		return true
	default:
		return awsRegionLabel.MatchString(label)
	}
}

func trimAWSTail(labels []string) []string {
	n := len(labels)
	for n > 0 && isAWSTailLabel(labels[n-1]) {
		n--
	}
	return labels[:n]
}

// endsWithAWSS3Label reports whether the S3 service label is followed only by
// region, dualstack, fips or vpce labels.
func endsWithAWSS3Label(labels []string) bool {
	rest := trimAWSTail(labels)
	return len(rest) > 0 && isAWSS3Label(rest[len(rest)-1])
}

func awsSuffixOf(host string) (string, bool) {
	idx := slices.IndexFunc(awsSuffixPatterns, func(suffix string) bool {
		return matchesDomainOrSubdomain(host, suffix)
	})
	if idx < 0 {
		return "", false
	}
	return awsSuffixPatterns[idx], true
}

func dropBucketAlias(prefix []string) []string {
	if len(prefix) == 2 && prefix[0] == bucketAlias {
		return prefix[1:]
	}
	return prefix
}

// isVPCEBasePrefix matches the service alias labels of a VPCE endpoint: vpce-x or bucket.vpce-x.
func isVPCEBasePrefix(prefix []string) bool {
	p := dropBucketAlias(prefix)
	return len(p) == 1 && strings.HasPrefix(p[0], vpceLabel)
}

// isVPCEBaseHost reports whether the labels describe a VPCE service endpoint rather than a bucket.
func isVPCEBaseHost(prefix, tail []string) bool {
	return slices.Contains(tail, vpceLabel) && isVPCEBasePrefix(prefix)
}

// splitAWSS3Labels splits labels into the bucket prefix labels and the tail after the S3 service label.
func splitAWSS3Labels(labels []string) (prefix, tail []string, ok bool) {
	rest := trimAWSTail(labels)
	if len(rest) < 2 || !isAWSS3Label(rest[len(rest)-1]) {
		return nil, nil, false
	}
	return rest[:len(rest)-1], labels[len(rest):], true
}

func awsLabelsBucket(labels []string) (string, bool) {
	prefix, tail, ok := splitAWSS3Labels(labels)
	if !ok || isVPCEBaseHost(prefix, tail) {
		return "", false
	}
	return prefix[0], true
}

func awsBucketSubdomain(host string) (string, bool) {
	suffix, ok := awsSuffixOf(host)
	if !ok {
		return "", false
	}
	return awsLabelsBucket(prefixLabels(host, suffix))
}

// IsAWSHost determines if hostname matches native AWS S3 endpoints.
func IsAWSHost(host string) bool {
	h := normalizeHost(host)
	suffix, ok := awsSuffixOf(h)
	return ok && endsWithAWSS3Label(prefixLabels(h, suffix))
}
