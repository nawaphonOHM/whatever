// Package client provides managed cloud object storage client implementation
// wrapping the AWS SDK for Go v2.
package client

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/contracts"
)

func formatRequestChecksumCalc(calc aws.RequestChecksumCalculation) string {
	if calc == aws.RequestChecksumCalculationWhenRequired {
		return "WHEN_REQUIRED"
	}
	return "WHEN_SUPPORTED"
}

func formatResponseChecksumVal(val aws.ResponseChecksumValidation) string {
	if val == aws.ResponseChecksumValidationWhenRequired {
		return "WHEN_REQUIRED"
	}
	return "WHEN_SUPPORTED"
}

func resolveNonAWSChecksumBehavior(
	calc aws.RequestChecksumCalculation,
	provider contracts.Provider,
) string {
	if calc != aws.RequestChecksumCalculationWhenRequired {
		return BehaviorChecksumExplicit
	}
	if provider == contracts.ProviderGCS {
		return BehaviorChecksumSuppressed
	}
	return BehaviorChecksumEnforced
}

func resolveChecksumBehaviorChange(
	calc aws.RequestChecksumCalculation,
	provider contracts.Provider,
) string {
	if provider != contracts.ProviderAWS {
		return resolveNonAWSChecksumBehavior(calc, provider)
	}
	if calc == aws.RequestChecksumCalculationWhenRequired {
		return BehaviorChecksumExplicit
	}
	return BehaviorChecksumStandard
}

func resolvePathStyleBehaviorChange(opts s3.Options, provider contracts.Provider) string {
	if !opts.UsePathStyle {
		return BehaviorPathStyleStandard
	}
	if provider == contracts.ProviderAWS {
		return BehaviorPathStyleConfigured
	}
	return BehaviorPathStyleEnforced
}
