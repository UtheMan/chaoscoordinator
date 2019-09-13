package network

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"strconv"
)

func BeginLatencyIncrease(subID string, flags Flags) error {
	scriptContent, err := cmdutil.LoadScript("scripts/latency.sh")
	if err != nil {
		return err
	}
	r := NewCmdRequest(flags, subID, scriptContent)
	r = addValuesToRequest(r)
	return beginNetworkOperation(r)
}

func addValuesToRequest(r CmdRequest) CmdRequest {
	r.Flags = setDefaultValues(r.Flags)
	r = setUpParams(r)
	return r
}
func setDefaultValues(flags Flags) Flags {
	if flags.Duration == 0 {
		flags.Duration = 60
	}
	if flags.Latency == 0 {
		flags.Latency = 200
	}
	return flags
}

func setUpParams(request CmdRequest) CmdRequest {
	cmd := make([]string, 0)
	cmd = append(cmd, string(request.ScriptContent))
	durationParam := "duration"
	latencyParam := "latencyIncrease"
	cmdParams := make([]compute.RunCommandInputParameter, 0)
	durationValue := strconv.Itoa(request.Duration)
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
