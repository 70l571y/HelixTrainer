//go:build ignore

package main

func main() {
	apiRoutes := []string{
		"/v2/users",
		"/v2/projects",
	}

	adminRoutes := []string{
		"/v1/admin/health",
		"/v1/admin/stats",
	}

	println(apiRoutes[0], adminRoutes[0])
}
