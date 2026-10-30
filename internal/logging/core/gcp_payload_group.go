package core

import "log/slog"

func cloneGroups(groups []string, next string) []string {
	if next == "" {
		res := make([]string, len(groups))
		copy(res, groups)
		return res
	}
	res := make([]string, len(groups)+1)
	copy(res, groups)
	res[len(groups)] = next
	return res
}

func getOrCreateGroup(curr map[string]any, g string) map[string]any {
	if next, ok := curr[g].(map[string]any); ok {
		return next
	}
	next := make(map[string]any)
	curr[g] = next
	return next
}

func navigateToGroup(m map[string]any, groups []string) map[string]any {
	curr := m
	for _, g := range groups {
		if g != "" {
			curr = getOrCreateGroup(curr, g)
		}
	}
	return curr
}

func applyReplace(groups []string, a slog.Attr, replace ReplaceAttrFunc) (slog.Attr, bool) {
	res := a
	if replace != nil {
		res = replace(groups, a)
	}
	return res, res.Key != ""
}

func addGroupAttrToMap(res map[string]any, groups []string, a slog.Attr, replace ReplaceAttrFunc) {
	val := a.Value.Resolve()
	if a.Key == "" {
		addAttrsToMap(res, groups, val.Group(), replace)
		return
	}
	subGroups := cloneGroups(groups, a.Key)
	subMap := groupToMap(subGroups, val.Group(), replace)
	if len(subMap) > 0 {
		mergeGroupIntoTarget(res, a.Key, subMap)
	}
}

func addLeafAttrToMap(res map[string]any, groups []string, a slog.Attr, replace ReplaceAttrFunc) {
	attr, ok := applyReplace(groups, a, replace)
	if !ok {
		return
	}
	subGroups := cloneGroups(groups, attr.Key)
	val := resolveAttrValue(subGroups, attr.Value, replace)
	insertDisambiguatedLeaf(res, attr.Key, val)
}

func mergeMapEntry(dst map[string]any, k string, v any) {
	if subSrc, ok := v.(map[string]any); ok {
		mergeGroupIntoTarget(dst, k, subSrc)
	} else if _, exists := dst[k]; !exists {
		dst[k] = v
	} else {
		cand := findDisambiguatedKey(dst, k)
		dst[cand] = v
	}
}

func mergeMaps(dst, src map[string]any) {
	for k, v := range src {
		mergeMapEntry(dst, k, v)
	}
}

func groupToMap(groups []string, attrs []slog.Attr, replace ReplaceAttrFunc) map[string]any {
	res := make(map[string]any, len(attrs))
	for _, a := range attrs {
		val := a.Value.Resolve()
		if val.Kind() == slog.KindGroup {
			addGroupAttrToMap(res, groups, a, replace)
		} else {
			addLeafAttrToMap(res, groups, a, replace)
		}
	}
	return res
}
