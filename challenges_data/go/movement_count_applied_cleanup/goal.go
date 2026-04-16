//go:build ignore

package main

var configRows = []string{
	"alpha_name=worker",
	"alpha_port=9000",
	"alpha_env=prod",
	"beta_name=api",
	"beta_port=8080",
	"beta_env=stage",
	"gamma_name=cron",
	"gamma_mode=draft",
	"gamma_env=nightly",
}
