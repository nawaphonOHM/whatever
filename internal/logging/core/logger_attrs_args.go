package core

import (
	"fmt"
	"log/slog"
)

func argsToAttrs(args []any) []slog.Attr {
	if len(args) == 0 {
		return nil
	}
	attrs := make([]slog.Attr, 0, (len(args)+1)/2)
	for i := 0; i < len(args); {
		attr, step := attrAt(args, i)
		attrs = append(attrs, attr)
		i += step
	}
	return attrs
}

func attrAt(args []any, index int) (slog.Attr, int) {
	if attr, ok := args[index].(slog.Attr); ok {
		return attr, 1
	}
	if key, ok := args[index].(string); ok {
		return stringAttr(args, index, key)
	}
	return genericAttr(args, index)
}

func stringAttr(args []any, index int, key string) (slog.Attr, int) {
	if index+1 < len(args) {
		return slog.Any(key, args[index+1]), 2
	}
	return slog.String("!BADKEY", key), 1
}

func genericAttr(args []any, index int) (slog.Attr, int) {
	if index+1 < len(args) {
		return slog.Any(fmt.Sprintf("%v", args[index]), args[index+1]), 2
	}
	return slog.Any("!EXTRA", args[index]), 1
}
