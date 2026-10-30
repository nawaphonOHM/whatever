package core

import (
	"fmt"
	"slices"
)

func isReservedKey(key string) bool {
	switch key {
	case "severity", "timestamp", "message", "httpRequest":
		return true
	}
	return false
}

func sortedPayloadKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func isKeyOccupied(m, payload map[string]any, k string) bool {
	if _, ok := m[k]; ok {
		return true
	}
	_, ok := payload[k]
	return ok
}

func findIndexedFreeKey(m, payload map[string]any, key string) string {
	for idx := 1; ; idx++ {
		cand := fmt.Sprintf("custom_%s_%d", key, idx)
		if !isKeyOccupied(m, payload, cand) {
			return cand
		}
	}
}

func findFreeKey(m, payload map[string]any, key string) string {
	cand := "custom_" + key
	if !isKeyOccupied(m, payload, cand) {
		return cand
	}
	return findIndexedFreeKey(m, payload, key)
}

func mergeReservedPayloadKey(m, payload map[string]any, k string) {
	freeKey := findFreeKey(m, payload, k)
	m[freeKey] = payload[k]
}

func mergeNonReservedKeys(m, payload map[string]any, keys []string) {
	for _, k := range keys {
		if !isReservedKey(k) {
			m[k] = payload[k]
		}
	}
}

func mergeReservedKeys(m, payload map[string]any, keys []string) {
	for _, k := range keys {
		if isReservedKey(k) {
			mergeReservedPayloadKey(m, payload, k)
		}
	}
}

func mergePayload(m, payload map[string]any) {
	keys := sortedPayloadKeys(payload)
	mergeNonReservedKeys(m, payload, keys)
	mergeReservedKeys(m, payload, keys)
}
