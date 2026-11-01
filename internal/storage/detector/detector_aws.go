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

var awsRegionLabel = regexp.MustCompile(`^[a-z]{2}(-[a-z]+)+-[0-9]+$`)

func isAWSS3Label(label string) bool {
	return label == "s3" || strings.HasPrefix(label, "s3-") || strings.HasPrefix(label, "s3express-")
}

// isAWSTailLabel reports whether label may sit between the S3 service label and the AWS suffix.
func isAWSTailLabel(label string) bool {
	switch label {
	case "dualstack", "fips", "vpce":
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

// IsAWSHost determines if hostname matches native AWS S3 endpoints.
func IsAWSHost(host string) bool {
	h := strings.ToLower(host)
	suffix, ok := awsSuffixOf(h)
	return ok && endsWithAWSS3Label(prefixLabels(h, suffix))
}
