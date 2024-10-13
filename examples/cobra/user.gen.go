//go:build user

// Code generated by Cobra annotation processor. DO NOT EDIT.
// versions:
//
//	go: go1.22.4
//	go-annotation: 0.1.0
//	Cobra: 1.0.0
//
package main

import (
	_imp0 "github.com/YReshetko/go-annotation/examples/cobra/commands"
	"github.com/spf13/cobra"
)

func init() {
	root = &cobra.Command{}
	root.Use = "cli"
	root.Example = "cli [-F file | -D dir] ... [-f format] profile"
	root.Short = "Root command of the application (short)"
	root.Long = "Root command of the application (long)"
	root.Flags().StringP("output", "o", "", "output")
	if err := root.MarkFlagRequired("output"); err != nil {
		fatal(err)
	}
	root.Flags().Int("num", 42, "some")
	root.PersistentFlags().BoolP("is-ok", "i", true, "some")
	root.RunE = func(cmd *cobra.Command, args []string) error {
		executor := _imp0.RootCommand{}
		if err := parseFlags(cmd, &executor); err != nil {
			return err
		}
		return executor.Run(cmd, args)
	}
	root.PreRunE = func(cmd *cobra.Command, args []string) error {
		executor := _imp0.RootCommand{}
		if err := parseFlags(cmd, &executor); err != nil {
			return err
		}
		return executor.PreRun(cmd, args)
	}
	root.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		executor := _imp0.RootCommand{}
		if err := parseFlags(cmd, &executor); err != nil {
			return err
		}
		return executor.PersistPreRun(cmd, args)
	}
	_cmd1 := &cobra.Command{}
	_cmd1.Use = "get"
	_cmd1.Example = "cli get [-F file | -D dir] ... [-f format] profile"
	_cmd1.Short = "Child command of the application (short)"
	_cmd1.Long = "Child command of the application (long)"
	_cmd1.Flags().StringP("output", "o", "", "output")
	if err := _cmd1.MarkFlagRequired("output"); err != nil {
		fatal(err)
	}
	_cmd1.Flags().IntP("num", "n", 42, "some")
	_cmd1.Flags().Int32P("num", "n", 42, "some")
	_cmd1.RunE = func(cmd *cobra.Command, args []string) error {
		executor := _imp0.ChildCommand{}
		if err := parseFlags(cmd, &executor); err != nil {
			return err
		}
		return executor.Run(cmd, args)
	}
	_cmd1.PostRunE = func(cmd *cobra.Command, args []string) error {
		executor := _imp0.ChildCommand{}
		if err := parseFlags(cmd, &executor); err != nil {
			return err
		}
		return executor.PostRun(cmd, args)
	}
	_cmd1.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		executor := _imp0.ChildCommand{}
		if err := parseFlags(cmd, &executor); err != nil {
			return err
		}
		return executor.PersistPostRun(cmd, args)
	}
	root.AddCommand(_cmd1)
}
