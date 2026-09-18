package server

import "fmt"

// registrationScan holds mutable validation state.
type registrationScan struct {
	seen map[string]bool
	out  []*validatedRoute
}

// appendRoute validates one API and appends it when valid.
func (s *registrationScan) appendRoute(
	reg *RestAPIRegistration,
	api *ExportableAPI,
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
	reg *RestAPIRegistration,
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
func (s *registrationScan) scanAPIs(reg *RestAPIRegistration) error {
	for apiIdx, api := range reg.Apis {
		if err := s.appendRoute(reg, api, apiIdx); err != nil {
			return err
		}
	}
	return nil
}

// scanRegistration validates one registration group.
func (s *registrationScan) scanRegistration(
	reg *RestAPIRegistration,
	regIdx int,
) error {
	if err := requireRegistration(reg, regIdx); err != nil {
		return err
	}
	return s.scanAPIs(reg)
}

// validateRegistrations validates registrations and returns routes.
func validateRegistrations(
	registrations []*RestAPIRegistration,
) ([]*validatedRoute, error) {
	scan := &registrationScan{
		seen: make(map[string]bool),
		out:  make([]*validatedRoute, 0),
	}
	for regIdx, reg := range registrations {
		if err := scan.scanRegistration(reg, regIdx); err != nil {
			return nil, err
		}
	}
	return scan.out, nil
}
