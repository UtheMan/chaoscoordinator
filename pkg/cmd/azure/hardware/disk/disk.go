package disk

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"time"
)

type Flags struct {
	cmdutil.GeneralFlags
	Amount string
}

type CmdRequest struct {
	Flags
	cmdutil.GeneralRequest
	subID string
}

func NewCmdRequest(flags Flags, subID string, scriptContent []byte) CmdRequest {
	return CmdRequest{
		Flags:          flags,
		GeneralRequest: cmdutil.NewGeneralRequest(scriptContent),
		subID:          subID,
	}
}

func beginDiskOperation(r CmdRequest) error {
	return cmdutil.ExecuteAzureScript(r)
}

func (r CmdRequest) HandleRequest() error {
	client, err := vm.NewVMsClient(r.subID)
	if err != nil {
		return err
	}
	vms, err := r.PrepareCmd(client)
	if err != nil {
		return err
	}
	for i := range vms {
		println("Executing disk cmd on machine", *vms[i].Name)
		err = r.BeginCmdExec(client, *vms[i].InstanceID)
		if err != nil {
			return err
		}
	}
	return err
}

func (r CmdRequest) PrepareCmd(client *vm.VmssClient) ([]compute.VirtualMachineScaleSetVM, error) {
	vms, err := cmdutil.GetEligibleVms(client, r.GeneralFlags)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

func (r CmdRequest) BeginCmdExec(client *vm.VmssClient, instanceId string) error {
	cmdErrors := make(chan error)
	results := make(chan compute.RunCommandResult)
	go r.ExecCmd(cmdErrors, results, client, instanceId)
	return cmdutil.HandleCmdResult(cmdErrors, results, r.GeneralFlags)
}

func (r CmdRequest) ExecCmd(cmdErrors chan error, results chan compute.RunCommandResult, client *vm.VmssClient, vmId string) (chan error, chan compute.RunCommandResult) {
	future, err := client.RunCommand(context.TODO(), r.ResourceGroup, r.ScaleSet, vmId, compute.RunCommandInput{
		CommandID:  &r.CmdID,
		Script:     &r.Cmd,
		Parameters: &r.CmdParams,
	})
	if err != nil {
		cmdErrors <- err
		return cmdErrors, results
	}
	time.Sleep(time.Duration(r.Duration)*time.Second + time.Duration(r.Timeout)*time.Second)
	result, err := future.Result(client.VirtualMachineScaleSetVMsClient)
	if err != nil {
		cmdErrors <- err
		return cmdErrors, results
	}
	results <- result
	return cmdErrors, results
}
