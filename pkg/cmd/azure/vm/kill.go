package vm

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"math/rand"
	"strings"
	"time"
)

// Flags exposes the vm reboot command arguments/flags
type Flags struct {
	Mode          string
	Scope         string
	ResourceGroup string
	ResourceName  string
}

func Kill(subID string, flags Flags) error {
	c, err := NewVMsClient(subID)
	if err != nil {
		return err
	}
	vmsList, err := c.list(context.TODO(), flags.ResourceGroup, flags.ResourceName)
	if err != nil {
		return err
	}
	instanceToReboot := selectInstanceToReboot(vmsList, flags)
	_, err = c.Restart(context.TODO(), flags.ResourceGroup, flags.ResourceName, instanceToReboot)
	if err != nil {
		return err
	}
	return err
}

func selectInstanceToReboot(setVMS []compute.VirtualMachineScaleSetVM, flags Flags) string {
	if strings.Compare("random", flags.Mode) == 0 {
		rand.Seed(time.Now().UnixNano())
		return *setVMS[rand.Intn(len(setVMS))].InstanceID
	} else {
		return *setVMS[0].InstanceID
	}
}

func (c *vmssClient) list(ctx context.Context, group, name string) ([]compute.VirtualMachineScaleSetVM, error) {
	result, err := c.VirtualMachineScaleSetVMsClient.List(ctx, group, name, "", "", string(compute.InstanceView))
	if err != nil {
		return nil, err
	}
	vms := make([]compute.VirtualMachineScaleSetVM, 0, len(result.Values()))
	vms = append(vms, result.Values()...)
	return vms, nil
}
