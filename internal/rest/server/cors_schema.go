package server

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	jsonschema "github.com/nawaphonOHM/whatever/v3/internal/json_schema"
	"github.com/xeipuuv/gojsonschema"
)

// corsSchemaDraftURL is the dialect gojsonschema supports; the checked-in
// 2020-12 schema is rewritten to it in memory before compilation.
const corsSchemaDraftURL = "https://json-schema.org/draft-07/schema#"

var (
	compiledCORSSchema     *gojsonschema.Schema
	compiledCORSSchemaOnce sync.Once
	compiledCORSSchemaErr  error
)

func extractEnumVal(httpHeaders map[string]any) any {
	props, ok := httpHeaders["properties"].(map[string]any)
	if !ok {
		return nil
	}
	return props["enum"]
}

func fixHTTPHeadersDef(defs map[string]any) {
	httpHeaders, ok := defs["http-headers"].(map[string]any)
	if !ok {
		return
	}
	enumVal := extractEnumVal(httpHeaders)
	if enumVal != nil {
		httpHeaders["items"] = map[string]any{"type": "string", "enum": enumVal}
		delete(httpHeaders, "properties")
	}
}

func adaptDefinitions(doc map[string]any) {
	defs, ok := doc["$defs"].(map[string]any)
	if !ok {
		return
	}
	doc["definitions"] = defs
	delete(doc, "$defs")
	fixHTTPHeadersDef(defs)
}

func adaptPropertyRef(propMap map[string]any) {
	if ref, ok := propMap["$ref"].(string); ok {
		propMap["$ref"] = strings.ReplaceAll(ref, "#/$defs/", "#/definitions/")
	}
}

func adaptPropertyPattern(propMap map[string]any) {
	items, ok := propMap["items"].(map[string]any)
	if !ok {
		return
	}
	if pattern, ok := items["pattern"].(string); ok {
		items["pattern"] = strings.ReplaceAll(pattern, `\\.`, `\.`)
	}
}

func adaptProperty(propMap map[string]any) {
	adaptPropertyRef(propMap)
	adaptPropertyPattern(propMap)
}

func adaptPropItem(val any) {
	if propMap, ok := val.(map[string]any); ok {
		adaptProperty(propMap)
	}
}

func adaptProperties(doc map[string]any) {
	props, ok := doc["properties"].(map[string]any)
	if !ok {
		return
	}
	for _, val := range props {
		adaptPropItem(val)
	}
}

func adaptCORSSchema(raw []byte) ([]byte, error) {
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	doc["$schema"] = corsSchemaDraftURL
	adaptDefinitions(doc)
	adaptProperties(doc)
	return json.Marshal(doc)
}

func getCompiledCORSSchema() (*gojsonschema.Schema, error) {
	compiledCORSSchemaOnce.Do(func() {
		adapted, err := adaptCORSSchema(jsonschema.CORSSchema)
		if err != nil {
			compiledCORSSchemaErr = fmt.Errorf("failed to adapt CORS schema: %w", err)
			return
		}
		loader := gojsonschema.NewBytesLoader(adapted)
		compiledCORSSchema, compiledCORSSchemaErr = gojsonschema.NewSchema(loader)
	})
	return compiledCORSSchema, compiledCORSSchemaErr
}
