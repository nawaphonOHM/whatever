// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/logging/central"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

// Storage client compatibility behavior description constants.
const (
	BehaviorChecksumSuppressed  = "suppress_unsolicited_trailers_for_gcs"
	BehaviorChecksumEnforced    = "enforce_when_required_for_compat"
	BehaviorChecksumExplicit    = "explicit_checksum_configuration"
	BehaviorChecksumStandard    = "standard_aws_defaults"
	BehaviorPathStyleEnforced   = "enforced_path_style_for_provider"
	BehaviorPathStyleConfigured = "configured_path_style"
	BehaviorPathStyleStandard   = "standard_virtual_hosted"
)

// Storage client authentication mode constants.
const (
	AuthModeStaticCredentials = "static_credentials"
	AuthModeDefaultChain      = "default_credential_chain"
)

func resolveAuthMode(cfg *config.Config) string {
	if hasStaticCredentials(cfg) {
		return AuthModeStaticCredentials
	}
	return AuthModeDefaultChain
}

func hasDriverOverrides(opts *Options) bool {
	if opts == nil {
		return false
	}
	return len(opts.DriverOptions) > 0
}

type compatLogParams struct {
	cfg      *config.Config
	provider contracts.Provider
	opts     *Options
	s3Opts   s3.Options
}

func logCompatibilityDiagnostics(ctx context.Context, p compatLogParams) {
	logger := central.DefaultWorker().Logger()
	if !isDebugLoggingEnabled(ctx, p.cfg, logger) {
		return
	}
	logger.DebugContext(ctx, "storage client compatibility configuration",
		"provider", string(p.provider),
		"request_checksum_calculation", formatRequestChecksumCalc(p.s3Opts.RequestChecksumCalculation),
		"response_checksum_validation", formatResponseChecksumVal(p.s3Opts.ResponseChecksumValidation),
		"checksum_behavior_change", resolveChecksumBehaviorChange(p.s3Opts.RequestChecksumCalculation, p.provider),
		"use_path_style", p.s3Opts.UsePathStyle,
		"path_style_behavior_change", resolvePathStyleBehaviorChange(p.s3Opts, p.provider),
		"disable_log_output_checksum_validation_skipped", p.s3Opts.DisableLogOutputChecksumValidationSkipped,
		"auth_mode", resolveAuthMode(p.cfg),
		"has_driver_overrides", hasDriverOverrides(p.opts),
	)
}

func logClientInitAndCompat(ctx context.Context, p compatLogParams) {
	logInitialization(ctx, lifecycleLogParams{
		provider: p.provider,
		bucket:   p.cfg.Bucket,
		s3Opts:   p.s3Opts,
	})
	logCompatibilityDiagnostics(ctx, p)
}
