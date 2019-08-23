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
	Amount        string
	ResourceGroup string
	ScaleSet      string
	Vm            compute.VirtualMachineScaleSet
	CmdID         string
	Cmd           []string
	CmdParams     []compute.RunCommandInputParameter
}

type Flags struct {
	Duration      int
	TimeOut       int
	Amount        string
	ResourceGroup string
	ScaleSetName  string
	Filter        string
}

func NewCmdRequest(amount string, duration int, timeOut int, resourceGroup string, scaleSet string, cmd []string, cmdParams []compute.RunCommandInputParameter) *CmdRequest {
	r := new(CmdRequest)
	r.Duration = time.Duration(duration)
	r.Timeout = time.Duration(timeOut)
	r.Amount = amount
	r.ResourceGroup = resourceGroup
	r.ScaleSet = scaleSet
	r.Cmd = cmd
	r.CmdID = "RunShellScript"
	r.CmdParams = cmdParams
	return r
}

func LoadScript(path string) ([]byte, error) {
	scriptContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return scriptContent, nil
}

func RunVmCmd(cmdErrors chan error, results chan compute.RunCommandResult, client *vm.VmssClient, request *CmdRequest, vmId string) (chan error, chan compute.RunCommandResult) {
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
