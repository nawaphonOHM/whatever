package core

import (
	"log/slog"
	"math"
	"strconv"
	"time"
)

const floatPrecisionBits = 64

func resolveErrorValue(val slog.Value) (any, bool) {
	err, ok := val.Any().(error)
	if !ok {
		return nil, false
	}
	if isTypedNil(err) {
		return nil, true
	}
	return formatErrorObject(err, val.Any())
}

func formatFloatKind(val slog.Value) (any, bool) {
	f := val.Float64()
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "!ERROR:json: unsupported value: " + strconv.FormatFloat(f, 'g', -1, floatPrecisionBits), true
	}
	return nil, false
}

func formatTemporalKind(val slog.Value) (any, bool) {
	if val.Kind() == slog.KindTime {
		return val.Time().Format(time.RFC3339Nano), true
	}
	if val.Kind() == slog.KindDuration {
		return val.Duration().String(), true
	}
	return nil, false
}

func formatSpecialKind(val slog.Value) (any, bool) {
	if res, ok := formatTemporalKind(val); ok {
		return res, true
	}
	if val.Kind() == slog.KindFloat64 {
		return formatFloatKind(val)
	}
	return nil, false
}

func resolveLeafValue(val slog.Value) any {
	if errVal, ok := resolveErrorValue(val); ok {
		return errVal
	}
	if formatted, ok := formatSpecialKind(val); ok {
		return formatted
	}
	return resolveOrdinaryValue(val)
}

func resolveOrdinaryValue(val slog.Value) any {
	if val.Kind() == slog.KindAny {
		return safeMarshalValue(val.Any())
	}
	return val.Any()
}

func resolveAttrValue(groups []string, v slog.Value, replace ReplaceAttrFunc) any {
	val := v.Resolve()
	if val.Kind() == slog.KindGroup {
		return groupToMap(groups, val.Group(), replace)
	}
	return resolveLeafValue(val)
}
