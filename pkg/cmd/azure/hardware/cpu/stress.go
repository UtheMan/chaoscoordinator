package cpu

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"strconv"
)

func BeginStress(subID string, flags cmdutil.Flags) error {
	client, err := vm.NewVMsClient(subID)
	if err != nil {
		return err
	}
	flags = setDefaultValues(flags)
	vms, cmdRequest, err := cmdutil.PrepareRequest("scripts/stressCpu.sh", flags, client)
	if err != nil {
		return err
	}
	cmdRequest = setUpParams(cmdRequest)
	for i := range vms {
		println("Stressing cpu on machine", *vms[i].Name)
		err := cmdutil.ExecutePreparedCmd(client, flags, *vms[i].InstanceID, cmdRequest)
		if err != nil {
			return err
		}
	}
	return err
}

func setUpParams(request *cmdutil.CmdRequest) *cmdutil.CmdRequest {
	cmd := make([]string, 0)
	cmd = append(cmd, string(request.ScriptContent))
	durationParam := "duration"
	cmdParams := make([]compute.RunCommandInputParameter, 0)
	durationValue := strconv.Itoa(int(request.Duration))
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &durationParam,
		Value: &durationValue,
	})
	request.Cmd = cmd
	request.CmdParams = cmdParams
	return request
}

func setDefaultValues(flags cmdutil.Flags) cmdutil.Flags {
	if flags.TimeOut == 0 {
		flags.TimeOut = 45
	}
	if flags.Duration == 0 {
		flags.Duration = 60
	}
	return flags
}