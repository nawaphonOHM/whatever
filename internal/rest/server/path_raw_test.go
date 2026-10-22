package server

import "github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"

// rawPathCases returns raw group coverage rows.
func rawPathCases() []*pathCase {
	raw := contracts.APIGroupRaw
	return []*pathCase{
		{name: "raw empty all", group: raw, expected: testRootPath},
		{name: "raw empty prefix with slash path", group: raw, path: testRootPath, expected: testRootPath},
		{name: "raw slash prefix with empty path", group: raw, prefix: testRootPath, expected: testRootPath},
		{
			name: "raw unversioned simple", group: raw,
			prefix: testItemsPrefix, path: testListPath, expected: "/items/list",
		},
		{
			name: "raw missing slashes", group: raw,
			prefix: "items", path: "list", expected: "/items/list",
		},
		{
			name: "raw trailing slash", group: raw,
			prefix: "/items/", path: "/list/", expected: "/items/list",
		},
		{
			name: "raw multiple slashes", group: raw,
			prefix: "///custom///", path: "///path///", expected: "/custom/path",
		},
		{
			name: "raw prefix with version in path", group: raw,
			prefix: "/v1/custom", path: "/hook", expected: "/v1/custom/hook",
		},
		{
			name: "raw empty prefix with custom metrics", group: raw,
			path: "/metrics-custom", expected: "/metrics-custom",
		},
		{
			name: "raw prefix with api", group: raw,
			prefix: "/api/raw", path: "/test", expected: "/api/raw/test",
		},
		{
			name: "raw ignores version number", group: raw,
			version: 1, prefix: "/users", path: "/:id", expected: "/users/:id",
		},
		{
			name: "raw ignores version 2", group: raw,
			version: 2, prefix: "/orders", path: "/status", expected: "/orders/status",
		},
		{
			name: "raw reserved health exact", group: raw,
			path: ReservedHealthPath, expected: ReservedHealthPath,
		},
		{
			name: "raw reserved health prefix", group: raw,
			prefix: ReservedHealthPath, expected: ReservedHealthPath,
		},
		{
			name: "raw reserved ready exact", group: raw,
			path: ReservedReadyPath, expected: ReservedReadyPath,
		},
		{
			name: "raw reserved ready prefix", group: raw,
			prefix: ReservedReadyPath, expected: ReservedReadyPath,
		},
	}
}
