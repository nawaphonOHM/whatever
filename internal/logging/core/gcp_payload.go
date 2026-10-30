package core

import (
	"fmt"
	"log/slog"
)

func cloneMap(m map[string]any) map[string]any {
	if m == nil {
		return make(map[string]any)
	}
	res := make(map[string]any, len(m))
	copyMapEntries(res, m)
	return res
}

func copyMapEntries(dst, src map[string]any) {
	for k, v := range src {
		if sub, ok := v.(map[string]any); ok {
			dst[k] = cloneMap(sub)
		} else {
			dst[k] = v
		}
	}
}

func findDisambiguatedKey(target map[string]any, key string) string {
	for idx := 1; ; idx++ {
		cand := fmt.Sprintf("%s_%d", key, idx)
		if _, exists := target[cand]; !exists {
			return cand
		}
	}
}

func insertDisambiguatedLeaf(target map[string]any, key string, val any) {
	if _, exists := target[key]; !exists {
		target[key] = val
		return
	}
	cand := findDisambiguatedKey(target, key)
	target[cand] = val
}

func insertLeafAttr(m map[string]any, groups []string, a slog.Attr, replace ReplaceAttrFunc) {
	attr, ok := applyReplace(groups, a, replace)
	if !ok {
		return
	}
	target := navigateToGroup(m, groups)
	subGroups := cloneGroups(groups, attr.Key)
	val := resolveAttrValue(subGroups, attr.Value, replace)
	insertDisambiguatedLeaf(target, attr.Key, val)
}

func mergeGroupIntoTarget(target map[string]any, key string, subMap map[string]any) {
	existing, exists := target[key]
	if !exists {
		target[key] = subMap
		return
	}
	if existingMap, ok := existing.(map[string]any); ok {
		mergeMaps(existingMap, subMap)
		return
	}
	cand := findDisambiguatedKey(target, key)
	target[cand] = subMap
}

func insertGroupAttr(m map[string]any, groups []string, a slog.Attr, replace ReplaceAttrFunc) {
	val := a.Value.Resolve()
	if a.Key == "" {
		addAttrsToMap(m, groups, val.Group(), replace)
		return
	}
	subGroups := cloneGroups(groups, a.Key)
	subMap := groupToMap(subGroups, val.Group(), replace)
	if len(subMap) > 0 {
		target := navigateToGroup(m, groups)
		mergeGroupIntoTarget(target, a.Key, subMap)
	}
}

func addAttrsToMap(m map[string]any, groups []string, attrs []slog.Attr, replace ReplaceAttrFunc) {
	for _, a := range attrs {
		val := a.Value.Resolve()
		if val.Kind() == slog.KindGroup {
			insertGroupAttr(m, groups, a, replace)
		} else {
			insertLeafAttr(m, groups, a, replace)
		}
	}
}
