package script

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

type Flags struct {
	Timeout       int
	Duration      int
	Multiple      int
	Kind          string
	Filter        string
	ResourceGroup string
	ScaleSet      string
	Path          string
}

type CmdRequest struct {
	Flags
	GeneralRequest
	subID string
}

func newGeneralRequest(scriptContent []byte) GeneralRequest {
	return GeneralRequest{
		ScriptContent: scriptContent,
		CmdID:         "RunShellScript",
	}
}

func newCmdRequest(flags Flags, subID string, scriptContent []byte) CmdRequest {
	return CmdRequest{
		Flags:          flags,
		GeneralRequest: newGeneralRequest(scriptContent),
		subID:          subID,
	}
}

func loadScript(path string) ([]byte, error) {
	scriptContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return scriptContent, nil
}

func handleCmdResult(cmdErrors chan error, results chan compute.RunCommandResult, flags Flags) error {
	select {
	case err := <-cmdErrors:
		return err
	case res := <-results:
		println(*(*res.Value)[0].Message)
		return nil
	case <-time.After(2 * (time.Duration(flags.Duration) + time.Duration(flags.Timeout)) * time.Second):
		err := errors.New("operation timed out")
		return err
	}
}

func getEligibleVms(client *vm.VmssClient, flags Flags) ([]compute.VirtualMachineScaleSetVM, error) {
	vmsList, err := client.VirtualMachineScaleSetVMsClient.List(context.TODO(), flags.ResourceGroup, flags.ScaleSet, flags.Filter, "", string(compute.InstanceView))
	if err != nil {
		return nil, err
	}
	vms := make([]compute.VirtualMachineScaleSetVM, 0, len(vmsList.Values()))
	vms = append(vms, vmsList.Values()...)
	return vms, nil
}

func loadArgs(args []string) map[string]string {
	argsMap := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		argsMap[args[i]] = args[i+1]
	}
	return argsMap
}

func addParams(r CmdRequest, args map[string]string) CmdRequest {
	cmd := make([]string, 0)
	cmd = append(cmd, string(r.ScriptContent))
	cmdParams := setParams(args)
	r.Cmd = cmd
	r.CmdParams = cmdParams
	return r
}

func setParams(args map[string]string) []compute.RunCommandInputParameter {
	params := make([]compute.RunCommandInputParameter, 0)
	for k := range args {
		paramName := k
		paramVal := args[k]
		params = append(params, compute.RunCommandInputParameter{
			Name:  &paramName,
			Value: &paramVal,
		})
	}
	return params
}
