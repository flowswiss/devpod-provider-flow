package cmd

import (
	"context"
	"github.com/flowswiss/devpod-provider-flow/pkg/flow"
	"github.com/flowswiss/devpod-provider-flow/pkg/options"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init an instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		options, err := options.FromEnv(true)
		if err != nil {
			return err
		}

		return flow.NewFlow(options.Token).Init(context.Background())
	},
}
