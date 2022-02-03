package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	requests "github.com/mfojtik/oc-adm-apiserver/pkg/cmd/requests"
)

type Options struct {
	configFlags *genericclioptions.ConfigFlags

	genericclioptions.IOStreams
}

func NewCmdApiserver(streams genericclioptions.IOStreams) *cobra.Command {
	o := &Options{IOStreams: streams}

	cmd := &cobra.Command{
		Use:          "apiserver",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.AddCommand(requests.NewCmdTopRequests(streams))

	return cmd
}

func (o *Options) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run() error {
	return nil
}

func main() {
	flags := pflag.NewFlagSet("dev-helpers", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := NewCmdApiserver(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
