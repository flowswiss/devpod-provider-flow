package cmd

import (
	"context"
	"github.com/flowswiss/devpod-provider-flow/pkg/flow"
	"github.com/flowswiss/devpod-provider-flow/pkg/options"
	"github.com/loft-sh/devpod/pkg/client"
	"github.com/loft-sh/devpod/pkg/log"
	"github.com/spf13/cobra"
	"time"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		flowClient := flow.NewFlow(options.Token)
		err = flowClient.StartInstanceByName(context.Background(), options.MachineID)
		if err != nil {
			return err
		}

		// wait until running
		for {
			status, err := flowClient.GetStatusByInstanceName(context.Background(), options.MachineID)
			if err != nil {
				log.Default.Errorf("Error retrieving instance status: %v", err)
				break
			}

			if status == client.StatusRunning {
				break
			}

			// make sure we don't spam
			time.Sleep(time.Second)
		}

		return nil
	},
}
