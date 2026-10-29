package server

import (
	"fmt"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// validatedRoute is a fully resolved route ready for engine mount.
type validatedRoute struct {
	Handler     contracts.Handler
	Path        string
	Method      string
	Middlewares []contracts.Middleware
}

// registrationScan holds mutable validation state.
type registrationScan struct {
	firstGroup contracts.APIGroup
	seen       map[string]bool
	out        []*validatedRoute
	firstIdx   int
	hasGroup   bool
}

// appendRoute validates one API and appends it when valid.
func (s *registrationScan) appendRoute(
	reg *contracts.RRestAPIRegistration,
	api *contracts.ExportableAPI,
	apiIdx int,
) error {
	if err := validateAPIEntry(api, apiIdx, reg.Prefix); err != nil {
		return err
	}
	route, err := buildValidatedRoute(reg, api, s.seen)
	if err != nil {
		return err
	}
	s.out = append(s.out, route)
	return nil
}

// requireRegistration rejects a nil registration group.
func requireRegistration(
	reg *contracts.RRestAPIRegistration,
	regIdx int,
) error {
	if reg != nil {
		return nil
	}
	return fmt.Errorf(
		"%w: registration at index %d is nil",
		ErrNilRegistration,
		regIdx,
	)
}

// scanAPIs validates every API entry in a registration group.
func (s *registrationScan) scanAPIs(reg *contracts.RRestAPIRegistration) error {
	for apiIdx, api := range reg.Apis {
		if err := s.appendRoute(reg, api, apiIdx); err != nil {
			return err
		}
	}
	return nil
}

// scanRegistration validates one registration group.
func (s *registrationScan) scanRegistration(
	reg *contracts.RRestAPIRegistration,
	regIdx int,
) error {
	if err := requireRegistration(reg, regIdx); err != nil {
		return err
	}
	if err := s.checkAPIGroup(reg, regIdx); err != nil {
		return err
	}
	return s.scanAPIs(reg)
}

// validationResult encapsulates validated routes and resolved API group.
type validationResult struct {
	group  contracts.APIGroup
	routes []*validatedRoute
}

// validateRegistrationsWithResult validates registrations and returns routes and resolved group.
func validateRegistrationsWithResult(
	registrations []*contracts.RRestAPIRegistration,
) (*validationResult, error) {
	scan := &registrationScan{
		seen: make(map[string]bool),
		out:  make([]*validatedRoute, 0),
	}
	for regIdx, reg := range registrations {
		if err := scan.scanRegistration(reg, regIdx); err != nil {
			return nil, err
		}
	}
	return &validationResult{
		routes: scan.out,
		group:  scan.effectiveGroup(),
	}, nil
}

// validateRegistrations validates registrations and returns routes.
func validateRegistrations(
	registrations []*contracts.RRestAPIRegistration,
) ([]*validatedRoute, error) {
	res, err := validateRegistrationsWithResult(registrations)
	if err != nil {
		return nil, err
	}
	return res.routes, nil
}
