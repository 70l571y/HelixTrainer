//go:build ignore

package main

var configRows = []string{
	"alpha_name=worker",
	"alpha_port=7000",
	"alpha_env=prod",
	"beta_name=api",
	"beta_port=8080",
	"beta_env=stage",
	"gamma_name=cron",
	"gamma_mode=drafr",
	"gamma_env=nightly",
}
