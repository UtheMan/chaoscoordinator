package cpu

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/utheman/chaoscoordinator/pkg/cmd/azure/hardware/cmdutil"
	"strconv"
)

func BeginStress(subID string, flags Flags) error {
	scriptContent, err := cmdutil.LoadScript("scripts/stressCpu.sh")
	if err != nil {
		return err
	}
	r := NewCmdRequest(flags, subID, scriptContent)
	r = addValuesToRequest(r)
	return beginCpuOperation(r)
}

func addValuesToRequest(r CmdRequest) CmdRequest {
	r.Flags = setDefaultValues(r.Flags)
	r = setUpParams(r)
	return r
}

func setDefaultValues(flags Flags) Flags {
	if flags.Timeout == 0 {
		flags.Timeout = 45
	}
	if flags.Duration == 0 {
		flags.Duration = 60
	}
	return flags
}

func setUpParams(request CmdRequest) CmdRequest {
	cmd := make([]string, 0)
	cmd = append(cmd, string(request.ScriptContent))
	durationParam := "duration"
	cmdParams := make([]compute.RunCommandInputParameter, 0)
	durationValue := strconv.Itoa(request.Duration)
	cmdParams = append(cmdParams, compute.RunCommandInputParameter{
		Name:  &durationParam,
		Value: &durationValue,
	})
	request.Cmd = cmd
	request.CmdParams = cmdParams
	return request
}
