package cmd

import (
	"context"
	"encoding/base64"
	"github.com/flowswiss/devpod-provider-flow/pkg/flow"
	"time"

	"github.com/flowswiss/devpod-provider-flow/pkg/options"
	"github.com/flowswiss/goclient/compute"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an instance",
	RunE: func(_ *cobra.Command, args []string) error {
		options, err := options.FromEnv(false)
		if err != nil {
			return err
		}

		req, err := buildCreateInstanceRequest(options)
		if err != nil {
			return err
		}

		flowClient := flow.NewFlow(options.Token)
		err = flowClient.CreateInstance(context.Background(), *req)
		if err != nil {
			return err
		}

		// wait until instance is available
		for {
			_, err = flowClient.GetStatusByInstanceName(context.Background(), options.MachineID)
			if err == nil {
				break
			}

			// make sure we don't spam
			time.Sleep(time.Second)
		}

		return nil
	},
}

func buildCreateInstanceRequest(opts *options.Options) (*compute.ServerCreate, error) {
	flowClient := flow.NewFlow(opts.Token)

	loc, err := flowClient.GetLocationByName(context.Background(), opts.Location)
	if err != nil {
		return nil, err
	}

	image, err := flowClient.GetImageByKey(context.Background(), opts.Image)
	if err != nil {
		return nil, err
	}

	product, err := flowClient.GetProductByName(context.Background(), opts.Product)
	if err != nil {
		return nil, err
	}

	network, err := flowClient.GetNetworkByName(context.Background(), opts.Network)
	if err != nil {
		return nil, err
	}

	keyPair, err := flowClient.CreateKeyPair(context.Background(), opts.MachineID, opts.MachineFolder)
	if err != nil {
		return nil, err
	}

	initScript, err := getInjectKeypairScript(opts.MachineFolder)
	if err != nil {
		return nil, err
	}

	serverCreateReq := &compute.ServerCreate{
		Name:             opts.MachineID,
		LocationID:       loc.ID,
		ImageID:          image.ID,
		ProductID:        product.ID,
		AttachExternalIP: true,
		NetworkID:        network.ID,
		KeyPairID:        keyPair.ID,
		CloudInit:        initScript,
	}
	return serverCreateReq, nil
}

func getInjectKeypairScript(dir string) (string, error) {
	publicKey, err := flow.GetMachinePublicKey(dir)
	if err != nil {
		return "", err
	}

	resultScript := `#!/bin/sh
useradd devpod -d /home/devpod
mkdir -p /home/devpod
if grep -q sudo /etc/groups; then
	usermod -aG sudo devpod
elif grep -q wheel /etc/groups; then
	usermod -aG wheel devpod
fi
echo "devpod ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/91-devpod
mkdir -p /home/devpod/.ssh
echo "` + string(publicKey) + `" >> /home/devpod/.ssh/authorized_keys
chmod 0700 /home/devpod/.ssh
chmod 0600 /home/devpod/.ssh/authorized_keys
chown -R devpod:devpod /home/devpod`

	return base64.StdEncoding.EncodeToString([]byte(resultScript)), nil
}
