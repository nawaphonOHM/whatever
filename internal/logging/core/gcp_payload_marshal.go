package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

const panicPrefix = "!PANIC:"

func isJSONMarshaler(v any) bool {
	_, ok := v.(json.Marshaler)
	return ok
}

func isTypedNil(v any) bool {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func marshalLeaf(v any) (value any, err error) {
	defer func() {
		if r := recover(); r != nil {
			value, err = nil, fmt.Errorf("%s %v", panicPrefix, r)
		}
	}()
	data, marshalErr := json.Marshal(v)
	if marshalErr != nil {
		return nil, marshalErr
	}
	return json.RawMessage(data), nil
}

func safeErrorText(err error) (text string) {
	defer func() {
		if r := recover(); r != nil {
			text = fmt.Sprintf("%s %v", panicPrefix, r)
		}
	}()
	return err.Error()
}

func formatErrorObject(err error, raw any) (any, bool) {
	if isJSONMarshaler(raw) {
		if value, mErr := marshalLeaf(raw); mErr == nil {
			return value, true
		}
	}
	return safeErrorText(err), true
}

func formatMarshalError(err error) string {
	if strings.HasPrefix(err.Error(), panicPrefix) {
		return err.Error()
	}
	return "!ERROR:" + err.Error()
}

func safeMarshalValue(v any) any {
	if v == nil {
		return nil
	}
	value, err := marshalLeaf(v)
	if err != nil {
		return formatMarshalError(err)
	}
	return value
}
