package disk

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"strconv"
)

func BeginFill(subID string, flags cmdutil.Flags) error {
	client, err := vm.NewVMsClient(subID)
	if err != nil {
		return err
	}
	vms, cmdRequest, err := prepareRequest(flags, err, client)
	if err != nil {
		return err
	}
	for i := range vms {
		println("Executing fill on machine", *vms[i].Name)
		err := fillUpDisk(client, flags, *vms[i].InstanceID, cmdRequest)
		if err != nil {
			return err
		}
	}
	return err
}

func prepareRequest(flags cmdutil.Flags, err error, client *vm.VmssClient) ([]compute.VirtualMachineScaleSetVM, *cmdutil.CmdRequest, error) {
	setDefaultValues(flags)
	vms, err := getEligibleVms(client, flags)
	if err != nil {
		return nil, nil, err
	}
	scriptContent, err := cmdutil.LoadScript("scripts/fillDisk.sh")
	if err != nil {
		return nil, nil, err
	}
	cmd, cmdParams := setUpCmd(scriptContent, flags)
	cmdRequest := cmdutil.NewCmdRequest(flags.Amount, flags.Duration, flags.TimeOut, flags.ResourceGroup, flags.ScaleSetName, cmd, cmdParams)
	return vms, cmdRequest, nil
}

func fillUpDisk(client *vm.VmssClient, flags cmdutil.Flags, instanceId string, request *cmdutil.CmdRequest) error {
	cmdErrors := make(chan error)
	results := make(chan compute.RunCommandResult)
	go cmdutil.RunVmCmd(cmdErrors, results, client, request, instanceId)
	return cmdutil.HandleCmdResult(cmdErrors, results, flags)
}

func setDefaultValues(flags cmdutil.Flags) {
	if flags.Amount == "" {
		flags.Amount = "1000"
	}
	if flags.TimeOut == 0 {
		flags.TimeOut = 45
	}
	if flags.Duration == 0 {
		flags.Duration = 60
	}
}

func getEligibleVms(client *vm.VmssClient, flags cmdutil.Flags) ([]compute.VirtualMachineScaleSetVM, error) {
	vmsList, err := client.VirtualMachineScaleSetVMsClient.List(context.TODO(), flags.ResourceGroup, flags.ScaleSetName, flags.Filter, "", string(compute.InstanceView))
	if err != nil {
		return nil, err
	}
	vms := make([]compute.VirtualMachineScaleSetVM, 0, len(vmsList.Values()))
	vms = append(vms, vmsList.Values()...)
	return vms, nil
}

func setUpCmd(scriptContent []byte, flags cmdutil.Flags) ([]string, []compute.RunCommandInputParameter) {
	cmd := make([]string, 0)
	cmd = append(cmd, string(scriptContent))
	durationParam := "duration"
	timeOutParam := "timeout"
	amountParam := "amount"
	cmdParams := make([]compute.RunCommandInputParameter, 0)
	durationValue := strconv.Itoa(flags.Duration)
	timeOutValue := strconv.Itoa(flags.TimeOut)
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &durationParam,
		Value: &durationValue,
	})
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &timeOutParam,
		Value: &timeOutValue,
	})
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &amountParam,
		Value: &flags.Amount,
	})
	return cmd, cmdParams
}
