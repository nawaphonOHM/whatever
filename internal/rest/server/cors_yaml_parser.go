package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

func checkTrailingContent(path string, dec *yaml.Decoder) error {
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return fmt.Errorf("multiple documents or trailing content detected in %q", path)
	}
	return nil
}

func unmarshalYAMLCFG(path string, data []byte) (*CORSFileConfig, error) {
	var fileCfg CORSFileConfig
	if err := yaml.UnmarshalWithOptions(data, &fileCfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML CORS configuration %q: %w", path, err)
	}
	return &fileCfg, nil
}

func decodeYAMLCORS(path string, data []byte) (any, error) {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	var doc any
	if err := dec.Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to parse YAML CORS configuration %q: %w", path, err)
	}
	if err := checkTrailingContent(path, dec); err != nil {
		return nil, err
	}
	return doc, nil
}

func convertYAMLDocToJSON(path string, doc any) ([]byte, error) {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert YAML CORS configuration %q to JSON: %w", path, err)
	}
	return jsonData, validateCORSJSON(path, jsonData)
}

func parseYAMLCORS(path string, data []byte) ([]byte, *CORSFileConfig, error) {
	doc, err := decodeYAMLCORS(path, data)
	if err != nil {
		return nil, nil, err
	}
	jsonData, err := convertYAMLDocToJSON(path, doc)
	if err != nil {
		return nil, nil, err
	}
	fileCfg, err := unmarshalYAMLCFG(path, data)
	return jsonData, fileCfg, err
}

func isYAMLExtension(ext string) bool {
	return ext == ".yaml" || ext == ".yml"
}

func parseCORSFile(path string, data []byte) ([]byte, *CORSFileConfig, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".json" {
		return parseJSONCORS(path, data)
	}
	if isYAMLExtension(ext) {
		return parseYAMLCORS(path, data)
	}
	return nil, nil, fmt.Errorf("unsupported CORS configuration file format %q (expected .json or .yaml)", path)
}
