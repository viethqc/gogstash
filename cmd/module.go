package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/viethqc/gogstash/KDGoLib/cliutil/cobrather"
)

// command line flags
var (
	flagConfig = &cobrather.StringFlag{
		Name:    "config",
		Default: "",
		Usage:   "Path to configuration file, default search path: config.json, config.yml",
		EnvVar:  "CONFIG",
	}
	flagDebug = &cobrather.BoolFlag{
		Name:    "debug",
		Default: false,
		Usage:   "Enable debug logging",
		EnvVar:  "DEBUG",
	}
	flagPProf = &cobrather.StringFlag{
		Name:    "pprof",
		Default: "",
		Usage:   "Enable golang pprof for listening address, ex: localhost:6060",
		EnvVar:  "PPROF",
	}
)

// modules
var (
	WorkerModule *cobrather.Module
	Module       *cobrather.Module
)

func init() {
	// WorkerModule info
	WorkerModule = &cobrather.Module{
		Use:   "worker",
		Short: "gogstash worker mode",
		RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return gogstash(ctx, flagConfig.String(), flagDebug.Bool(), flagPProf.String(), true)
		},
	}

	// Module info
	Module = &cobrather.Module{
		Use:   "gogstash",
		Short: "Logstash like, written in golang",
		Commands: []*cobrather.Module{
			cobrather.VersionModule,
			WorkerModule,
		},
		GlobalFlags: []cobrather.Flag{
			flagConfig,
			flagDebug,
			flagPProf,
		},
		RunE: func(ctx context.Context, cmd *cobra.Command, args []string) error {
			return gogstash(ctx, flagConfig.String(), flagDebug.Bool(), flagPProf.String(), false)
		},
	}
}
