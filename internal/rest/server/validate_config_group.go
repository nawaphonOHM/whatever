package server

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// isKnownConfigGroup checks if the configuration group is standard, raw, or empty default.
func isKnownConfigGroup(g string) bool {
	return g == "" || g == string(contracts.APIGroupStandard) || g == string(contracts.APIGroupRaw)
}

// resolveConfigAPIGroup validates and returns the effective server configuration API group.
// Empty group defaults to standard; unrecognized values return ErrInvalidAPIGroup.
func resolveConfigAPIGroup(group string) (contracts.APIGroup, error) {
	if !isKnownConfigGroup(group) {
		return "", fmt.Errorf(
			"%w: %q is not a valid API group (allowed: %q, %q)",
			ErrInvalidAPIGroup,
			group,
			contracts.APIGroupStandard,
			contracts.APIGroupRaw,
		)
	}
	if group == string(contracts.APIGroupRaw) {
		return contracts.APIGroupRaw, nil
	}
	return contracts.APIGroupStandard, nil
}

// errAPIGroupMismatch formats an API group mismatch validation error.
func errAPIGroupMismatch(regGroup contracts.APIGroup, cfgGroup string) error {
	return fmt.Errorf(
		"%w: registrations have group %q, but server is configured with %q",
		ErrAPIGroupMismatch,
		regGroup,
		cfgGroup,
	)
}

// matchGroupOrError verifies that the first registration group matches the expected group.
func matchGroupOrError(
	regGroup, expectedGroup contracts.APIGroup,
	cfgGroup string,
) error {
	if regGroup != expectedGroup {
		return errAPIGroupMismatch(regGroup, cfgGroup)
	}
	return nil
}

// checkGroupMatch resolves the registration group and validates it against expectedGroup.
func checkGroupMatch(
	reg *contracts.RRestAPIRegistration,
	expectedGroup contracts.APIGroup,
	cfgGroup string,
) error {
	regGroup, err := resolveAPIGroup(reg.Group, 0)
	if err != nil {
		return err
	}
	return matchGroupOrError(regGroup, expectedGroup, cfgGroup)
}

// checkEffectiveGroupMatch verifies that the registrations' group matches the expected server API group.
func checkEffectiveGroupMatch(
	registrations []*contracts.RRestAPIRegistration,
	expectedGroup contracts.APIGroup,
	cfgGroup string,
) error {
	if len(registrations) == 0 {
		return nil
	}
	return checkGroupMatch(registrations[0], expectedGroup, cfgGroup)
}

// validateConfigAndRoutes validates the registrations and checks group match.
func validateConfigAndRoutes(
	registrations []*contracts.RRestAPIRegistration,
	expectedGroup contracts.APIGroup,
	cfgGroup string,
) ([]*validatedRoute, error) {
	validated, err := validateRegistrations(registrations)
	if err != nil {
		return nil, err
	}
	if err := checkEffectiveGroupMatch(registrations, expectedGroup, cfgGroup); err != nil {
		return nil, err
	}
	return validated, nil
}

// validateRegistrationsWithConfig validates registrations against server configuration.
func validateRegistrationsWithConfig(
	registrations []*contracts.RRestAPIRegistration,
	cfg *Config,
) ([]*validatedRoute, error) {
	c := resolveConfig(cfg)
	expectedGroup, err := resolveConfigAPIGroup(c.APIGroup)
	if err != nil {
		return nil, err
	}
	return validateConfigAndRoutes(registrations, expectedGroup, c.APIGroup)
}
