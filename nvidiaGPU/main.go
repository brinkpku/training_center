package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
)

var deviceID uint

func init() {
	flag.UintVar(&deviceID, "did", 0, "device id")
	flag.Parse()
}

func main() {
	nvml.Init()
	defer nvml.Shutdown()
	device, err := nvml.NewDevice(deviceID)
	if err != nil {
		fmt.Println("nvml new device error:", err)
		return
	}
	attr, err := device.GetAttributes()
	if err != nil {
		fmt.Println("get attr error:", err)
	}
	fmt.Printf("device %v attr: %v\n", deviceID, attr)
	for i := 0; i < 2; i++ {
		status, err := device.Status()
		if err != nil {
			fmt.Println("get status error:", err)
			return
		}
		// fmt.Printf("device %v status: %v\n", deviceID, status)
		fmt.Println("memory free", *status.Memory.Global.Free, "used", *status.Memory.Global.Used)
		fmt.Println("utilization", *status.Utilization.Memory)
		fmt.Println("temperature", *status.Temperature)
		fmt.Println("power", *status.Power)
		time.Sleep(time.Second * 5)
	}
}
