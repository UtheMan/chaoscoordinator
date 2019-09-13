package vm

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure"
	"math/rand"
	"time"
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

func PickRandom(vms []compute.VirtualMachineScaleSetVM, amount int) []compute.VirtualMachineScaleSetVM {
	if amount > len(vms) {
		amount = len(vms)
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	result := make([]compute.VirtualMachineScaleSetVM, amount)
	perm := r.Perm(len(vms))
	for i := 0; i < amount; i++ {
		randIndex := perm[i]
		result[i] = vms[randIndex]
	}
	return result
}
