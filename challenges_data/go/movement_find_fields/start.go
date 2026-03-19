//go:build ignore

package main

func activeColumns() []string {
	return []string{
		"user_id,status,created_at",
		"order_id,total,created_at",
		"team_id,region,created_at",
	}
}
