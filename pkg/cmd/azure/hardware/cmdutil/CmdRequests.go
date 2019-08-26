package cmdutil

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-03-01/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"io/ioutil"
	"time"
)

type CmdRequest struct {
	Duration      time.Duration
	Timeout       time.Duration
	Latency       int
	Amount        string
	ResourceGroup string
	ScaleSet      string
	Vm            compute.VirtualMachineScaleSet
	CmdID         string
	ScriptContent []byte
	Cmd           []string
	CmdParams     []compute.RunCommandInputParameter
}

func (request CmdRequest) RunVmCmd(cmdErrors chan error, results chan compute.RunCommandResult, client *vm.VmssClient, vmId string) (chan error, chan compute.RunCommandResult) {
	future, err := client.RunCommand(context.TODO(), request.ResourceGroup, request.ScaleSet, vmId, compute.RunCommandInput{
		CommandID:  &request.CmdID,
		Script:     &request.Cmd,
		Parameters: &request.CmdParams,
	})
	if err != nil {
		cmdErrors <- err
		return cmdErrors, results
	}
	time.Sleep(request.Duration*time.Second + request.Timeout*time.Second)
	result, err := future.Result(client.VirtualMachineScaleSetVMsClient)
	if err != nil {
		cmdErrors <- err
		return cmdErrors, results
	}
	results <- result
	return cmdErrors, results
}

type Flags struct {
	Duration      int
	Latency       int
	TimeOut       int
	Amount        string
	ResourceGroup string
	ScaleSetName  string
	Filter        string
}

func NewCmdRequest(amount string, duration int, timeOut int, latency int, resourceGroup string, scaleSet string, scriptContent []byte) *CmdRequest {
	r := new(CmdRequest)
	r.Amount = amount
	r.Duration = time.Duration(duration)
	r.Timeout = time.Duration(timeOut)
	r.Latency = latency
	r.ResourceGroup = resourceGroup
	r.ScaleSet = scaleSet
	r.CmdID = "RunShellScript"
	r.ScriptContent = scriptContent
	return r
}

func LoadScript(path string) ([]byte, error) {
	scriptContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return scriptContent, nil
}

func ExecutePreparedCmd(client *vm.VmssClient, flags Flags, instanceId string, request *CmdRequest) error {
	cmdErrors := make(chan error)
	results := make(chan compute.RunCommandResult)
	go request.RunVmCmd(cmdErrors, results, client, instanceId)
	return HandleCmdResult(cmdErrors, results, flags)
}

func PrepareRequest(path string, flags Flags, client *vm.VmssClient) ([]compute.VirtualMachineScaleSetVM, *CmdRequest, error) {
	vms, err := getEligibleVms(client, flags)
	if err != nil {
		return nil, nil, err
	}
	scriptContent, err := LoadScript(path)
	if err != nil {
		return nil, nil, err
	}
	cmdRequest := NewCmdRequest(flags.Amount, flags.Duration, flags.TimeOut, flags.Latency, flags.ResourceGroup, flags.ScaleSetName, scriptContent)
	return vms, cmdRequest, nil
}

func HandleCmdResult(cmdErrors chan error, results chan compute.RunCommandResult, flags Flags) error {
	select {
	case err := <-cmdErrors:
		return err
	case res := <-results:
		println(*(*res.Value)[0].Message)
		return nil
	case <-time.After(2 * time.Duration(flags.TimeOut+flags.Duration) * time.Second):
		err := errors.New("operation timed out")
		return err
	}
}

func getEligibleVms(client *vm.VmssClient, flags Flags) ([]compute.VirtualMachineScaleSetVM, error) {
	vmsList, err := client.VirtualMachineScaleSetVMsClient.List(context.TODO(), flags.ResourceGroup, flags.ScaleSetName, flags.Filter, "", string(compute.InstanceView))
	if err != nil {
		return nil, err
	}
	vms := make([]compute.VirtualMachineScaleSetVM, 0, len(vmsList.Values()))
	vms = append(vms, vmsList.Values()...)
	return vms, nil
}
