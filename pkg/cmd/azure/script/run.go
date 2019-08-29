package script

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"time"
)

func Run(subID string, flags Flags, args []string) error {
	scriptContent, err := loadScript(flags.Path)
	if err != nil {
		return err
	}
	loadedArgs := loadArgs(args)
	r := newCmdRequest(flags, subID, scriptContent)
	r = addParams(r, loadedArgs)
	return executeRequest(r)
}

func executeRequest(r CmdRequest) error {
	return r.HandleRequest()
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
		println("Executing script:", r.Flags.Kind, "on machine", *vms[i].Name)
		err = r.BeginCmdExec(client, *vms[i].InstanceID)
		if err != nil {
			return err
		}
	}
	return err
}

func (r CmdRequest) PrepareCmd(client *vm.VmssClient) ([]compute.VirtualMachineScaleSetVM, error) {
	vms, err := getEligibleVms(client, r.Flags)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

func (r CmdRequest) BeginCmdExec(client *vm.VmssClient, instanceId string) error {
	cmdErrors := make(chan error)
	results := make(chan compute.RunCommandResult)
	go r.ExecCmd(cmdErrors, results, client, instanceId)
	return handleCmdResult(cmdErrors, results, r.Flags)
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
	time.Sleep(2 * time.Duration(r.Duration) * time.Second)
	result, err := future.Result(client.VirtualMachineScaleSetVMsClient)
	if err != nil {
		cmdErrors <- err
		return cmdErrors, results
	}
	results <- result
	return cmdErrors, results
}
