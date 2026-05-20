package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davaba86/agentfiles/internal/initcmd"
	"github.com/davaba86/agentfiles/internal/migrate"
	"github.com/davaba86/agentfiles/internal/validate"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	switch os.Args[1] {
	case "-h", "--help", "help":
		usage()
	case "-v", "--version", "version":
		fmt.Fprintf(os.Stdout, "agentfiles %s\n", version)
	case "init":
		if _, err := initcmd.Run(cwd, os.Stdout); err != nil {
			fatal(err)
		}
	case "check":
		result, err := validate.Run(cwd, os.Stdout)
		if err != nil {
			fatal(err)
		}
		if !result.OK() {
			os.Exit(1)
		}
	case "migrate":
		fs := flag.NewFlagSet("migrate", flag.ExitOnError)
		dryRun := fs.Bool("dry-run", false, "show what would happen without changing files")
		backup := fs.Bool("backup", true, "create backups when modifying files")
		if err := fs.Parse(os.Args[2:]); err != nil {
			fatal(err)
		}
		if err := migrate.Run(cwd, migrate.Options{DryRun: *dryRun, Backup: *backup}, os.Stdout); err != nil {
			fatal(err)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: agentfiles <command>")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  init                 create starter instruction files")
	fmt.Fprintln(os.Stderr, "  check                validate instruction files")
	fmt.Fprintln(os.Stderr, "  migrate [--dry-run]  migrate CLAUDE.md rules into AGENTS.md")
	fmt.Fprintln(os.Stderr, "  version              print version")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
