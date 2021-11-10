package main

import (
	"flag"
	"fmt"
	"strconv"
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
	for i := 0; i < 2; i++ {
		deviceID = uint(i)
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
		status, err := device.Status()
		if err != nil {
			fmt.Println("get status error:", err)
			return
		}
		println(device.Model)
		// fmt.Printf("device %v status: %v\n", deviceID, status)
		fmt.Println("memory free", *status.Memory.Global.Free, "used", *status.Memory.Global.Used)
		fmt.Println("utilization", *status.Utilization.Memory)
		free, _ := strconv.ParseInt(fmt.Sprintf("%d", *status.Memory.Global.Free), 10, 64)
		used, _ := strconv.ParseInt(fmt.Sprintf("%d", *status.Memory.Global.Used), 10, 64)
		fmt.Println("free", free)
		fmt.Println("used", used)
		total := free + used
		if total == 0 {
			fmt.Println("get toal memory 0")
			return
		}
		fmt.Println(float64(used) / float64(total))
		fmt.Println("temperature", *status.Temperature)
		fmt.Println("power", *status.Power)
		time.Sleep(time.Second * 5)
	}
}
