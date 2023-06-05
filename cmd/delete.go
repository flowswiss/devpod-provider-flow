package cmd

import (
	"context"
	"github.com/flowswiss/devpod-provider-flow/pkg/flow"
	"github.com/flowswiss/devpod-provider-flow/pkg/options"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		return flow.NewFlow(options.Token).DeleteInstanceByName(context.Background(), options.MachineID)
	},
}
