package cli

import (
	"d2rinfo/config"
	"d2rinfo/server"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Server ServerCmd `cmd:"" help:"Starts the d2rinfo API server." aliases:"serve"`
}

type ServerCmd struct {
	ConfigFile string `short:"c" help:"Path to the config file"`
	Host       string `short:"h" help:"Host to listen on."`
	Port       int    `short:"p" help:"Port to listen on."`
}

func (s *ServerCmd) Run() error {
	var cfg *config.Config
	if s.ConfigFile != "" {
		cfg = config.LoadConfig(s.ConfigFile)
	}
	if s.Host != "" {
		cfg.Host = s.Host
	}
	if s.Port != 0 {
		cfg.Port = s.Port
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}

	srv := server.New(cfg)
	srv.StartServer()
	return nil
}

func Execute() {
	var cli CLI
	ctx := kong.Parse(
		&cli,
		kong.Name("d2rinfo"),
		kong.Description(`A simple REST API middleman for D2Emu.`),
		kong.UsageOnError(),
	)

	if err := ctx.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
