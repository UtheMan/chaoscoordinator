package network

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/vm"
	"strconv"
)

func BeginLatencyIncrease(subID string, flags cmdutil.Flags) error {
	client, err := vm.NewVMsClient(subID)
	if err != nil {
		return err
	}
	flags = setDefaultValues(flags)
	vms, cmdRequest, err := cmdutil.PrepareRequest("scripts/latency.sh", flags, client)
	if err != nil {
		return err
	}
	cmdRequest = setUpParams(cmdRequest)
	for i := range vms {
		println("Increasing latency on machine", *vms[i].Name)
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
	latencyParam := "latencyIncrease"
	cmdParams := make([]compute.RunCommandInputParameter, 0)
	durationValue := strconv.Itoa(int(request.Duration))
	latencyValue := strconv.Itoa(request.Latency)
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &durationParam,
		Value: &durationValue,
	})
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &latencyParam,
		Value: &latencyValue,
	})
	request.Cmd = cmd
	request.CmdParams = cmdParams
	return request
}

func setDefaultValues(flags cmdutil.Flags) cmdutil.Flags {
	if flags.Duration == 0 {
		flags.Duration = 60
	}
	if flags.Latency == 0 {
		flags.Latency = 200
	}
	return flags
}
