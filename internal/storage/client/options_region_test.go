package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/client"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func TestExtractAWSRegion_Patterns(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		host     string
		expected string
	}{
		{name: "standard_s3_region", host: "s3.ap-southeast-1.amazonaws.com", expected: "ap-southeast-1"},
		{name: "virtual_hosted_s3_region", host: "mybucket.s3.ap-southeast-1.amazonaws.com", expected: "ap-southeast-1"},
		{name: "trailing_dot_region", host: "s3.ap-southeast-1.amazonaws.com.", expected: "ap-southeast-1"},
		{name: "api_aws_suffix", host: "s3.ap-southeast-1.api.aws", expected: "ap-southeast-1"},
		{name: "china_region", host: "s3.cn-north-1.amazonaws.com.cn", expected: "cn-north-1"},
		{name: "external_1_alias", host: "s3-external-1.amazonaws.com", expected: "us-east-1"},
		{name: "fips_dotted", host: "s3-fips.us-east-1.amazonaws.com", expected: "us-east-1"},
		{name: "fips_dashed", host: "s3-fips-us-east-1.amazonaws.com", expected: "us-east-1"},
		{name: "dualstack_dotted", host: "s3.dualstack.us-east-2.amazonaws.com", expected: "us-east-2"},
		{name: "dualstack_dashed", host: "s3-dualstack-us-east-2.amazonaws.com", expected: "us-east-2"},
		{name: "vpce_endpoint", host: "bucket.vpce-1a2b-3c4d.s3.us-east-1.vpce.amazonaws.com", expected: "us-east-1"},
		{name: "vpce_china", host: "vpce-1a2b-3c4d.s3.cn-north-1.vpce.amazonaws.com.cn", expected: "cn-north-1"},
		{name: "vpce_api_aws", host: "bucket.vpce-1a2b-3c4d.s3.us-west-2.vpce.api.aws", expected: "us-west-2"},
		{name: "global_s3_no_region", host: "s3.amazonaws.com", expected: ""},
		{name: "gcs_endpoint", host: "storage.googleapis.com", expected: ""},
		{name: "minio_custom", host: "minio.internal.net", expected: ""},
		{name: "localhost_custom", host: "localhost:9000", expected: ""},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			opts, err := client.BuildS3Options(&config.Config{
				Endpoint: "https://" + tc.host,
			}, contracts.ProviderAWS, nil)
			if tc.expected != "" {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, opts.Region)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "us-east-1", opts.Region)
			}
		})
	}
}

func TestResolveEffectiveRegion_Scenarios(t *testing.T) {
	t.Parallel()

	t.Run("nil_config_defaults_to_us_east_1", func(t *testing.T) {
		t.Parallel()
		opts, err := client.BuildS3Options(nil, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		assert.Equal(t, "us-east-1", opts.Region)
	})

	t.Run("empty_config_defaults_to_us_east_1", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		assert.Equal(t, "us-east-1", opts.Region)
	})

	t.Run("custom_region_preserved_for_custom_s3_compatible", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Endpoint: "https://minio.internal.net",
			Region:   "us-west-2",
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderS3Compatible, nil)
		require.NoError(t, err)
		assert.Equal(t, "us-west-2", opts.Region)
	})

	t.Run("custom_region_preserved_for_gcs", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Endpoint: "https://storage.googleapis.com",
			Region:   "auto",
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderGCS, nil)
		require.NoError(t, err)
		assert.Equal(t, "auto", opts.Region)
	})

	t.Run("default_region_adopts_detected_endpoint_region", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Endpoint: "https://s3.ap-southeast-1.amazonaws.com",
			Region:   "us-east-1",
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		assert.Equal(t, "ap-southeast-1", opts.Region)
	})

	t.Run("matching_configured_and_detected_region", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Endpoint: "https://s3.eu-west-1.amazonaws.com",
			Region:   "eu-west-1",
		}
		opts, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.NoError(t, err)
		assert.Equal(t, "eu-west-1", opts.Region)
	})

	t.Run("region_conflict_returns_error", func(t *testing.T) {
		t.Parallel()
		cfg := &config.Config{
			Endpoint: "https://s3.ap-southeast-1.amazonaws.com",
			Region:   "us-west-2",
		}
		_, err := client.BuildS3Options(cfg, contracts.ProviderAWS, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `region mismatch: config region "us-west-2" conflicts with endpoint region "ap-southeast-1"`)
	})

	t.Run("bucket_url_region_detection_and_conflict", func(t *testing.T) {
		t.Parallel()
		cfgMatch := &config.Config{
			BucketURL: "https://mybucket.s3.eu-central-1.amazonaws.com",
			Region:    "us-east-1",
		}
		optsMatch, errMatch := client.BuildS3Options(cfgMatch, contracts.ProviderAWS, nil)
		require.NoError(t, errMatch)
		assert.Equal(t, "eu-central-1", optsMatch.Region)

		cfgConflict := &config.Config{
			BucketURL: "https://mybucket.s3.eu-central-1.amazonaws.com",
			Region:    "ap-southeast-1",
		}
		_, errConflict := client.BuildS3Options(cfgConflict, contracts.ProviderAWS, nil)
		require.Error(t, errConflict)
		assert.Contains(t, errConflict.Error(), `region mismatch: config region "ap-southeast-1" conflicts with endpoint region "eu-central-1"`)
	})
}
