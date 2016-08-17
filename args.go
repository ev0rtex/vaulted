package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

var (
	ErrUnknownShell       = errors.New("Unknown shell")
	ErrTooManyArguments   = errors.New("too many arguments provided")
	ErrNotEnoughArguments = errors.New("not enough arguments provided")
)

func ParseArgs(args []string) (Command, error) {
	if len(args) == 0 {
		return nil, nil
	}

	switch args[0] {
	case "cp", "copy":
		return parseCopyArgs(args[1:])

	case "dump":
		return parseDumpArgs(args[1:])

	case "env":
		return parseEnvArgs(args[1:])

	case "ls", "list":
		return parseListArgs(args[1:])

	case "load":
		return parseLoadArgs(args[1:])

	case "rm":
		return parseRemoveArgs(args[1:])

	case "shell":
		return parseShellArgs(args[1:])

	case "upgrade":
		return parseUpgradeArgs(args[1:])

	default:
		return nil, nil
	}
}

func parseCopyArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("copy", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 2 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 2 {
		return nil, ErrTooManyArguments
	}

	c := &Copy{}
	c.OldVaultName = flag.Arg(0)
	c.NewVaultName = flag.Arg(1)
	return c, nil
}

func parseDumpArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("dump", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	d := &Dump{}
	d.VaultName = flag.Arg(0)
	return d, nil
}

func parseEnvArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("env", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	shell, err := detectShell()
	if err == ErrUnknownShell {
		shell = "sh"
	}

	usageHint := true
	fi, err := os.Stdout.Stat()
	if err == nil {
		if fi.Mode()&os.ModeCharDevice == 0 {
			usageHint = false
		}
	}

	e := &Env{}
	e.VaultName = flag.Arg(0)
	e.Shell = shell
	e.UsageHint = usageHint
	return e, nil
}

func parseListArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("list", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() > 0 {
		return nil, ErrTooManyArguments
	}

	return &List{}, nil
}

func parseLoadArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("load", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	l := &Load{}
	l.VaultName = flag.Arg(0)
	return l, nil
}

func parseRemoveArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("remove", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	r := &Remove{}
	r.VaultNames = flag.Args()
	return r, nil
}

func parseShellArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("shell", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	currentVaultedEnv := os.Getenv("VAULTED_ENV")
	if currentVaultedEnv != "" {
		return nil, fmt.Errorf("Refusing to spawn a new shell when already in environment '%s'.", currentVaultedEnv)
	}

	if flag.NArg() < 1 {
		return nil, ErrNotEnoughArguments
	}

	if flag.NArg() > 1 {
		return nil, ErrTooManyArguments
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	s := &Shell{}
	s.VaultName = flag.Arg(0)
	s.Command = []string{shell, "--login"}
	return s, nil
}

func parseUpgradeArgs(args []string) (Command, error) {
	flag := pflag.NewFlagSet("upgrade", pflag.ContinueOnError)
	err := flag.Parse(args)
	if err != nil {
		return nil, err
	}

	if flag.NArg() > 0 {
		return nil, ErrTooManyArguments
	}

	return &Upgrade{}, nil
}

func detectShell() (string, error) {
	shell := os.Getenv("SHELL")
	if shell != "" {
		return filepath.Base(shell), nil
	}

	return "", ErrUnknownShell
}