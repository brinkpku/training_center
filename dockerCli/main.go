package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/brinkpku/training_center/dockerCli/worker"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/container"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var dockerCli *command.DockerCli

var wName string

func init() {
	flag.StringVar(&wName, "wname", "tms", "")
	flag.Parse()
	var err error
	dockerCli, err = command.NewDockerCli()
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
		Name:         "test",
		Version:      1,
		Path:         "pathxxx",
		ZKConfig:     "zkxxx",
		AlgoConfig:   "algo.xxx",
		RenderConfig: "renderxxx",
	}
	if containers, err := ListAll(); err != nil {
		fmt.Println("list error", err)
	} else {
		for _, c := range containers {
			fmt.Println(c.Names, c.State, c.Command)
		}
	}
	if err := WorkerManager.Remove(app); err != nil {
		fmt.Println("remove container error", err)
	} else {
		fmt.Println("container removed")
	}
	WorkerManager.Start(app, 1)
	fmt.Println(wName)
	if running, err := WorkerManager.IsRunning("atm-test.v1"); err != nil {
		fmt.Println("inspect error:", err)
	} else {
		fmt.Println("running status:", running)
	}
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

const dockerName = "docker"

// Name ...
func (m *Manager) Name() string {
	return dockerName
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
	for retry := 0; retry < 3; retry++ {
		err = Execute(container.NewRunCommand, tmpl.StringSlice())
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
