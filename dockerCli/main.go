package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/brinkpku/training_center/dockerCli/worker"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types"
	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	networktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var dockerCli *command.DockerCli

var wName string

func init() {
	flag.StringVar(&wName, "wname", "tms", "")
	flag.Parse()
	var err error
	dockerCli, err = command.NewDockerCli() // different with official api init method
	if err != nil {
		panic(err)
	}
	if err = dockerCli.Initialize(flags.NewClientOptions(), command.WithInitializeClient(func(cli *command.DockerCli) (client.APIClient, error) {
		return client.NewClientWithOpts()
	})); err != nil {
		panic(err)
	}
}

func main() {
	c := Config{Image: "golang:1.13.14", Command: "bash"}
	InitManager(c)
	app := &Applet{
		Name:       "test",
		Version:    1,
		Path:       "/Users/youmengyuan",
		AlgoConfig: "/Users/yumengyuan",
	}

	// list and stop remove
	if containers, err := ListAll(); err != nil {
		fmt.Println("list error", err)
	} else {
		for _, c := range containers {
			fmt.Println(c.Names, c.State, c.Command)
			if c.Names[0] == "/atm-test.v1" {
				var timeout = new(time.Duration)
				*timeout = time.Second * 10
				if err := dockerCli.Client().ContainerStop(context.Background(), "atm-test.v1", timeout); err != nil {
					fmt.Println("stop error:", err)
				}
				dockerCli.Client().ContainerRemove(context.Background(), "atm-test.v1", types.ContainerRemoveOptions{})
			}
		}
	}
	// two methods of running a container
	// WorkerManager.Start(app, 1)
	_, err := dockerCli.Client().ContainerCreate(context.Background(),
		&containertypes.Config{Image: c.Image, Cmd: strslice.StrSlice{c.Command}, Tty: true},
		// nil, nil,
		&containertypes.HostConfig{
			// NetworkMode: containertypes.NetworkMode("viper-lite"),
			Mounts: []mount.Mount{{Type: mount.TypeBind, Source: "/Users/youmengyuan", Target: "/viper-lite/test"}},
			// Resources: containertypes.Resources{DeviceRequests: []containertypes.DeviceRequest{{
			// 	Driver:    "nivida",
			// 	Count:     1,
			// 	DeviceIDs: []string{"0"},
			// }}},
		},
		&networktypes.NetworkingConfig{},
		nil, "atm-test.v1")
	if err != nil {
		fmt.Println("create error:", err)
		return
	}
	err = dockerCli.Client().ContainerStart(context.Background(), "atm-test.v1", types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("start error:", err)
		return
	}
	// inspect is running
	if running, err := WorkerManager.IsRunning("atm-test.v1"); err != nil {
		fmt.Println("inspect error:", err)
	} else {
		fmt.Println("running status:", running)
	}
	WorkerManager.Stop(app)
	WorkerManager.Start(app, 1)
}

// WorkerManager vps worker manager
var WorkerManager *Manager

// Config ...
type Config struct {
	Network string `toml:"network"`
	Image   string `toml:"image"`
	Restart bool   `toml:"restart"`
	License string `toml:"license"`
	Command string `toml:"command"`
}

// Manager ...
type Manager struct {
	config Config
}

// InitManager ...
func InitManager(config Config) {
	WorkerManager = &Manager{config: config}
}

// Applet ...
type Applet struct {
	Name           string
	Version        int32
	Path           string
	ZKConfig       string
	AlgoConfig     string
	PipelineConfig string
	RenderConfig   string
}

// Start ...
func (m *Manager) Start(app *Applet, gpuID int) error {
	tmpl := worker.VPSWorker{
		Network: m.config.Network,
		Name:    fmt.Sprintf("atm-%s.v%d", app.Name, app.Version),
		// MemoryLimit: app.AlgoAppInstance.GetResource().GetMemory().GetLimit(),
		// CPULimit:    app.AlgoAppInstance.GetResource().GetCpu().GetLimit(),
		Volume: map[string]string{
			app.Path:       fmt.Sprintf("/viper-lite/apps/%s.v%d", app.Name, app.Version),
			app.AlgoConfig: "/viper-lite/license/client.lic",
		},
		// GPUS:    fmt.Sprintf("device=%d", gpuID),
		Image:   m.config.Image,
		Restart: m.config.Restart,
		Command: m.config.Command,
		Env:     []string{"Y=nb"},
	}

	backOff := 3

	log.Println(tmpl.StringSlice())
	var err error
	for retry := 0; retry < 1; retry++ {
		// err = Execute(container.NewRunCommand, tmpl.StringSlice())
		err = runContainerByCmd(tmpl.StringSlice())
		if err != nil {
			fmt.Println("run docker error", err)
			time.Sleep(time.Duration(backOff) * time.Second)
			backOff = backOff * 2
		} else {
			break
		}
	}

	return err
}

// Stop ...
func (m *Manager) Stop(app *Applet) error {
	name := fmt.Sprintf("atm-%s.v%d", app.Name, app.Version)
	return Execute(container.NewStopCommand, []string{name})
}

// Remove ...
func (m *Manager) Remove(app *Applet) error {
	name := fmt.Sprintf("atm-%s.v%d", app.Name, app.Version)
	return Execute(container.NewRmCommand, []string{name})
}

// IsRunning ...
func (m *Manager) IsRunning(cname string) (bool, error) {
	name := fmt.Sprintf(cname)
	c, err := Inspect(name)
	if err != nil {
		return false, err
	}
	return c.State.Running, nil
}

func runContainerByCmd(args []string) (err error) {
	cmd := exec.Command("docker", args...)
	data, err := cmd.Output()
	if err != nil {
		err = fmt.Errorf("failed to call Output(): %v", err)
		return
	}
	fmt.Println("hhhh", string(data))
	return
}

// Execute ...
func Execute(newCommandFunc func(dockerCli command.Cli) *cobra.Command, args []string) error {
	cmd := newCommandFunc(dockerCli)
	cmd.SetArgs(args)
	return cmd.Execute()
}

// Inspect ...
func Inspect(ref string) (types.ContainerJSON, error) {
	return dockerCli.Client().ContainerInspect(context.Background(), ref)
}

func ListAll() ([]types.Container, error) {
	return dockerCli.Client().ContainerList(context.Background(), types.ContainerListOptions{All: true})
}
