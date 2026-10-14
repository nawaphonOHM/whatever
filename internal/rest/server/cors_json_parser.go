package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func checkValidationErrors(path string, result *gojsonschema.Result) error {
	if result.Valid() {
		return nil
	}
	errMsgs := make([]string, 0, len(result.Errors()))
	for _, desc := range result.Errors() {
		errMsgs = append(errMsgs, desc.String())
	}
	return fmt.Errorf("CORS configuration file %q failed schema validation: %s", path, strings.Join(errMsgs, "; "))
}

func validateCORSJSON(path string, jsonData []byte) error {
	schema, err := getCompiledCORSSchema()
	if err != nil {
		return fmt.Errorf("failed to load CORS schema: %w", err)
	}
	result, err := schema.Validate(gojsonschema.NewBytesLoader(jsonData))
	if err != nil {
		return fmt.Errorf("failed to validate CORS configuration %q: %w", path, err)
	}
	return checkValidationErrors(path, result)
}

func checkTrailingJSON(path string, dec *json.Decoder) error {
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return fmt.Errorf("multiple JSON values or trailing content detected in %q", path)
	}
	return nil
}

func decodeJSONDoc(path string, data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	var doc any
	if err := dec.Decode(&doc); err != nil {
		return fmt.Errorf("failed to parse JSON CORS configuration %q: %w", path, err)
	}
	return checkTrailingJSON(path, dec)
}

func unmarshalJSONCFG(path string, data []byte) (*CORSFileConfig, error) {
	if err := validateCORSJSON(path, data); err != nil {
		return nil, err
	}
	var fileCfg CORSFileConfig
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON CORS configuration %q: %w", path, err)
	}
	return &fileCfg, nil
}

func parseJSONCORS(path string, data []byte) ([]byte, *CORSFileConfig, error) {
	if err := decodeJSONDoc(path, data); err != nil {
		return nil, nil, err
	}
	fileCfg, err := unmarshalJSONCFG(path, data)
	return data, fileCfg, err
}
