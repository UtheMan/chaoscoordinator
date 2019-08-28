package cmdutil

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"io/ioutil"
	"time"
)

type GeneralRequest struct {
	Vm            compute.VirtualMachineScaleSet
	CmdID         string
	ScriptContent []byte
	Cmd           []string
	CmdParams     []compute.RunCommandInputParameter
}
type GeneralFlags struct {
	Duration      int
	Timeout       int
	ResourceGroup string
	ScaleSet      string
	Filter        string
}
type Request interface {
	HandleRequest() error
	PrepareCmd(client *vm.VmssClient) ([]compute.VirtualMachineScaleSetVM, error)
	BeginCmdExec(client *vm.VmssClient, instanceId string) error
	ExecCmd(cmdErrors chan error, results chan compute.RunCommandResult, client *vm.VmssClient, vmId string) (chan error, chan compute.RunCommandResult)
}

func ExecuteAzureScript(r Request) error {
	return r.HandleRequest()
}

func NewGeneralRequest(scriptContent []byte) GeneralRequest {
	return GeneralRequest{
		ScriptContent: scriptContent,
		CmdID:         "RunShellScript",
	}
}

func LoadScript(path string) ([]byte, error) {
	scriptContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return scriptContent, nil
}

func HandleCmdResult(cmdErrors chan error, results chan compute.RunCommandResult, flags GeneralFlags) error {
	select {
	case err := <-cmdErrors:
		return err
	case res := <-results:
		println(*(*res.Value)[0].Message)
		return nil
	case <-time.After(2 * time.Duration(flags.Timeout+flags.Duration) * time.Second):
		err := errors.New("operation timed out")
		return err
	}
}

func GetEligibleVms(client *vm.VmssClient, flags GeneralFlags) ([]compute.VirtualMachineScaleSetVM, error) {
	vmsList, err := client.VirtualMachineScaleSetVMsClient.List(context.TODO(), flags.ResourceGroup, flags.ScaleSet, flags.Filter, "", string(compute.InstanceView))
	if err != nil {
		return nil, err
	}
	vms := make([]compute.VirtualMachineScaleSetVM, 0, len(vmsList.Values()))
	vms = append(vms, vmsList.Values()...)
	return vms, nil
}
