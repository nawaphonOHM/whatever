package server

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// isKnownGroup checks if group is standard, raw, or empty default.
func isKnownGroup(g contracts.APIGroup) bool {
	return g == "" || g == contracts.APIGroupStandard || g == contracts.APIGroupRaw
}

// errInvalidRegistrationGroupAt formats an invalid API registration group validation error.
func errInvalidRegistrationGroupAt(regIdx int, group contracts.APIGroup) error {
	return fmt.Errorf(
		"%w: registration at index %d has invalid group %q",
		ErrInvalidRegistrationGroup,
		regIdx,
		group,
	)
}

// resolveAPIGroup validates and returns the effective API group.
// Empty group defaults to standard; invalid values return ErrInvalidRegistrationGroup.
func resolveAPIGroup(group contracts.APIGroup, regIdx int) (contracts.APIGroup, error) {
	if !isKnownGroup(group) {
		return "", errInvalidRegistrationGroupAt(regIdx, group)
	}
	if group == contracts.APIGroupRaw {
		return contracts.APIGroupRaw, nil
	}
	return contracts.APIGroupStandard, nil
}

// errMixedAPIGroupsAt formats a mixed API groups validation error.
func errMixedAPIGroupsAt(
	currIdx int,
	currGroup contracts.APIGroup,
	firstIdx int,
	firstGroup contracts.APIGroup,
) error {
	return fmt.Errorf(
		"%w: registration at index %d has group %q, but registration at index %d has group %q",
		ErrMixedAPIGroups,
		currIdx,
		currGroup,
		firstIdx,
		firstGroup,
	)
}

// recordFirstGroup saves the first encountered API group in the scan.
func (s *registrationScan) recordFirstGroup(grp contracts.APIGroup, regIdx int) {
	s.firstGroup = grp
	s.firstIdx = regIdx
	s.hasGroup = true
}

// matchExistingGroup verifies group consistency with the first group in the scan.
func (s *registrationScan) matchExistingGroup(grp contracts.APIGroup, regIdx int) error {
	if grp != s.firstGroup {
		return errMixedAPIGroupsAt(regIdx, grp, s.firstIdx, s.firstGroup)
	}
	return nil
}

// effectiveGroup returns the resolved API group of the scan, defaulting to standard.
func (s *registrationScan) effectiveGroup() contracts.APIGroup {
	if s.hasGroup {
		return s.firstGroup
	}
	return contracts.APIGroupStandard
}

// checkAPIGroup verifies that all registrations share the same effective API group.
func (s *registrationScan) checkAPIGroup(
	reg *contracts.RRestAPIRegistration,
	regIdx int,
) error {
	grp, err := resolveAPIGroup(reg.Group, regIdx)
	if err != nil {
		return err
	}
	if !s.hasGroup {
		s.recordFirstGroup(grp, regIdx)
		return nil
	}
	return s.matchExistingGroup(grp, regIdx)
}
