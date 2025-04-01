package main

import (
	"fmt"
	"graphraggo/internal/bootstrap"
	"graphraggo/internal/global"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

const (
	condaEnvName = "graphrag-go"
)

func init() {
	// Host
	global.Host = "0.0.0.0"

	// Port
	global.Port = 8088

	// PythonServerPort
	global.PythonServerPort = 8089

	// WorkDir
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("fail to get work directory, err: %s", err.Error()))
	}
	global.WorkDir = dir

	// ExampleSettingFile
	global.ExampleSettingFile = fmt.Sprintf("%s/%s/settings-example.yaml", dir, global.KBDir)

	// PythonPath
	envName := "graphrag-go"
	cmd := exec.Command("conda", "run", "-n", envName, "which", "python")
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("fail to get python path, err: %s", err.Error()))
	}
	global.PythonPath = strings.TrimRight(string(out), "\n")

	// Print
	fmt.Printf("Host: %s\n", global.Host)
	fmt.Printf("Port: %d\n", global.Port)
	fmt.Printf("PythonServerPort: %d\n", global.PythonServerPort)
	fmt.Printf("ExampleSettingFile: %s\n", global.ExampleSettingFile)
	fmt.Printf("WorkDir: %s\n", global.WorkDir)
	fmt.Printf("PythonPath: %s\n", global.PythonPath)
}

func main() {
	// 启动 Python 服务
	go func() {
		bootstrap.MustInitPythonServer()
	}()

	r := bootstrap.MustInitRouter()

	if err := r.Run(fmt.Sprintf("%s:%d", global.Host, global.Port)); err != nil {
		slog.Error("failed to run server",
			slog.String("host", global.Host),
			slog.Int("port", global.Port),
			slog.String("err", err.Error()))
	}
}
