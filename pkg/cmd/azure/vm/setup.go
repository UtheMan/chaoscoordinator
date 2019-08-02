package vm

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
)

const userAgent string = "genesys"

type vmssClient struct {
	compute.VirtualMachineScaleSetVMsClient
}

func NewVMsClient(subId string) (*vmssClient, error) {
	vmClient, err := newScaleSetVMsClient(subId)
	if err != nil {
		return nil, err
	}
	return &vmssClient{vmClient}, nil
}

func newScaleSetVMsClient(subID string) (compute.VirtualMachineScaleSetVMsClient, error) {
	a, err := injectAuthorizer()

	if err != nil {
		return compute.VirtualMachineScaleSetVMsClient{}, err
	}
	client := compute.NewVirtualMachineScaleSetVMsClient(subID)
	client.Authorizer = a
	if err := client.AddToUserAgent(userAgent); err != nil {
		return compute.VirtualMachineScaleSetVMsClient{}, err
	}
	return client, nil
}
