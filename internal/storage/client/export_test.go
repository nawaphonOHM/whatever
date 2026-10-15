package client

// MapErrorForTest exposes mapError for testing error mapping directly.
func MapErrorForTest(err error) error {
	return mapError(err)
}
