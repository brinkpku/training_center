package worker

import (
	"fmt"
	"strings"
)

// VPSWorker ...
type VPSWorker struct {
	Name        string
	Network     string
	MemoryLimit string
	CPULimit    string
	Volume      map[string]string
	// https://stackoverflow.com/questions/25185405/using-gpu-from-a-docker-container
	GPUS    string
	Image   string
	Command string
	Restart bool
	Env     []string
}

// StringSlice ...
func (v VPSWorker) StringSlice() []string {
	ret := []string{"-dit"}
	if v.Name != "" {
		ret = append(ret, []string{"--name", v.Name}...)
	}
	if v.Network != "" {
		ret = append(ret, []string{"--network", v.Network}...)
	}
	if v.MemoryLimit != "" {
		ret = append(ret, []string{"--memory", v.MemoryLimit}...)
	}
	if v.CPULimit != "" {
		ret = append(ret, []string{"--cpus", v.CPULimit}...)
	}
	if v.GPUS != "" {
		ret = append(ret, []string{"--gpus", v.GPUS}...)
	}
	if v.Restart {
		ret = append(ret, []string{"--restart", "always"}...)
	}
	if len(v.Volume) != 0 {
		for src, dst := range v.Volume {
			ret = append(ret, []string{"--volume", fmt.Sprintf("%s:%s", src, dst)}...)
		}
	}
	if len(v.Env) != 0 {
		for _, e := range v.Env {
			ret = append(ret, "-e", e)
		}
	}
	if len(v.Image) != 0 {
		ret = append(ret, v.Image)
	}
	if len(v.Command) != 0 {
		ret = append(ret, strings.Split(v.Command, " ")...)
	}
	return ret
}
