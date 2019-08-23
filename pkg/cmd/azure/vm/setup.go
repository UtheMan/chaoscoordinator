package vm

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure"
)

const userAgent string = "genesys"

type VmssClient struct {
	compute.VirtualMachineScaleSetVMsClient
}

func NewVMsClient(subId string) (*VmssClient, error) {
	vmClient, err := newScaleSetVMsClient(subId)
	if err != nil {
		return nil, err
	}
	return &VmssClient{vmClient}, nil
}

func newScaleSetVMsClient(subID string) (compute.VirtualMachineScaleSetVMsClient, error) {
	a, err := azure.InjectAuthorizer()
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
