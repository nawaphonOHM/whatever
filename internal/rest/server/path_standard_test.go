package server

import "github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"

// standardPathCases returns standard and default group coverage rows.
func standardPathCases() []*pathCase {
	std := contracts.APIGroupStandard
	return []*pathCase{
		{name: "standard empty all", group: std, expected: testAPIPrefix},
		{name: "default empty all", expected: testAPIPrefix},
		{name: "standard empty all v1", group: std, version: 1, expected: "/api/v1"},
		{name: "default empty all v1", version: 1, expected: "/api/v1"},
		{
			name: "standard simple", group: std,
			prefix: testItemsPrefix, path: testListPath, expected: testAPIItemsList,
		},
		{name: "default missing slashes", prefix: "items", path: "list", expected: testAPIItemsList},
		{
			name: "standard trailing slash", group: std,
			prefix: "/items/", path: "/list/", expected: testAPIItemsList,
		},
		{
			name: "standard multiple slashes", group: std,
			prefix: "///items///", path: "///list///", expected: testAPIItemsList,
		},
		{
			name: "standard api in prefix", group: std,
			prefix: "/api/items", path: testListPath, expected: testAPIItemsList,
		},
		{
			name: "standard prefix api only", group: std,
			prefix: testAPIPrefix, path: testItemsPrefix, expected: "/api/items",
		},
		{
			name: "standard versioned prefix", group: std,
			version: 1, prefix: "/users", path: "/:id", expected: "/api/v1/users/:id",
		},
		{
			name: "standard explicit api deduplicated", group: std,
			version: 1, prefix: "/api/users", path: "/:id", expected: "/api/v1/users/:id",
		},
		{
			name: "standard prefix contains v1", group: std,
			version: 1, prefix: "/v1/users", expected: testAPIV1Users,
		},
		{
			name: "standard prefix contains api/v1", group: std,
			version: 1, prefix: "/v1/api/users", expected: testAPIV1Users,
		},
		{
			name: "standard prefix contains v2", group: std,
			version: 2, prefix: "/v2/items", path: testListPath, expected: "/api/v2/items/list",
		},
		{
			name: "standard prefix contains api/v2", group: std,
			version: 2, prefix: "/api/v2/items", path: testListPath, expected: "/api/v2/items/list",
		},
		{
			name: "standard versioned empty prefix", group: std,
			version: 1, path: testItemsPrefix, expected: "/api/v1/items",
		},
		{
			name: "standard root ping path", group: std,
			version: 1, path: "/ping", expected: "/api/v1/ping",
		},
		{
			name: "standard prefix api ping", group: std,
			version: 1, prefix: testAPIPrefix, path: "/ping", expected: "/api/v1/ping",
		},
		{
			name: "standard orders pre-prefixed", group: std,
			version: 1, prefix: "/api/orders", path: "/status", expected: "/api/v1/orders/status",
		},
		{
			name: "standard reserved health exact", group: std,
			path: ReservedHealthPath, expected: ReservedHealthPath,
		},
		{
			name: "standard reserved health prefix", group: std,
			prefix: ReservedHealthPath, expected: ReservedHealthPath,
		},
		{
			name: "standard reserved ready exact", group: std,
			path: ReservedReadyPath, expected: ReservedReadyPath,
		},
		{
			name: "standard reserved ready prefix", group: std,
			prefix: ReservedReadyPath, expected: ReservedReadyPath,
		},
	}
}
