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
	Port = 8080
)

func init() {
	// WorkDir
	dir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("fail to get work directory, err: %s", err.Error()))
	}

	global.WorkDir = dir
	global.ExampleSettingFile = fmt.Sprintf("%s/%s/settings-example.yaml", dir, global.KBDir)

	// PythonPath
	envName := "graphrag-go"
	cmd := exec.Command("conda", "run", "-n", envName, "which", "python")
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Sprintf("fail to get python path, err: %s", err.Error()))
	}

	global.PythonPath = strings.TrimRight(string(out), "\n")
}

func main() {
	r := bootstrap.MustInitRouter()

	if err := r.Run(fmt.Sprintf("0.0.0.0:%d", Port)); err != nil {
		slog.Error("failed to run server",
			slog.String("host", "0.0.0.0"),
			slog.Int("port", Port),
			slog.String("err", err.Error()))
	}
}
