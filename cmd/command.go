package cmd

import (
	"context"
	"fmt"
	"github.com/flowswiss/devpod-provider-flow/pkg/flow"
	"github.com/flowswiss/devpod-provider-flow/pkg/options"
	"github.com/loft-sh/devpod/pkg/ssh"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
)

var commandCmd = &cobra.Command{
	Use:   "command",
	Short: "Run a command on the instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		command := os.Getenv("COMMAND")
		if command == "" {
			return fmt.Errorf("command environment variable is missing")
		}

		// get private key
		privateKey, err := ssh.GetPrivateKeyRawBase(options.MachineFolder)
		if err != nil {
			return fmt.Errorf("load private key: %w", err)
		}

		elasticIP, err := flow.NewFlow(options.Token).GetElasticIPByInstanceName(context.Background(), options.MachineID)
		if err != nil {
			return err
		}

		// dial external address
		sshClient, err := ssh.NewSSHClient("devpod", elasticIP.PublicIP+":22", privateKey)
		if err != nil {
			return errors.Wrap(err, "create ssh client")
		}
		defer sshClient.Close()

		return ssh.Run(context.Background(), sshClient, command, os.Stdin, os.Stdout, os.Stderr)
	},
}
