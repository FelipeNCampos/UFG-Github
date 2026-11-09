package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/felip/runner/simulador-cli/internal/runner"
	"github.com/spf13/cobra"
)

const (
	defaultPort    = 8443
	defaultJavaBin = "java"
)

func defaultJarPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "simulador.jar"
	}
	return filepath.Join(home, ".hubsaude", "simulador.jar")
}

type runtimeFlags struct {
	JavaBin string
	JarPath string
}

func bindRuntimeFlags(command *cobra.Command, target *runtimeFlags) {
	flags := command.Flags()
	flags.StringVar(&target.JavaBin, "java-bin", defaultJavaBin, "Executavel Java a ser usado para invocar o JAR.")
	flags.StringVar(&target.JarPath, "jar", defaultJarPath(), "Caminho para o arquivo simulador.jar.")
}

func newRunnerConfig(rt runtimeFlags) runner.Config {
	return runner.Config{
		JavaBin: rt.JavaBin,
		JarPath: rt.JarPath,
	}
}

func validationError(message string, args ...any) error {
	return &ExitError{
		Code:    2,
		Message: fmt.Sprintf(message, args...),
	}
}

func emitResult(stdout io.Writer, stderr io.Writer, result runner.Result) error {
	if result.Stdout != "" {
		if _, err := io.WriteString(stdout, result.Stdout); err != nil {
			return err
		}
		if !strings.HasSuffix(result.Stdout, "\n") {
			if _, err := io.WriteString(stdout, "\n"); err != nil {
				return err
			}
		}
	}

	if result.Stderr != "" {
		if _, err := io.WriteString(stderr, result.Stderr); err != nil {
			return err
		}
		if !strings.HasSuffix(result.Stderr, "\n") {
			if _, err := io.WriteString(stderr, "\n"); err != nil {
				return err
			}
		}
	}

	if result.ExitCode != 0 {
		return &ExitError{Code: result.ExitCode}
	}

	return nil
}

func wrapRuntimeError(err error) error {
	if err == nil {
		return nil
	}
	return &ExitError{
		Code:    1,
		Message: err.Error(),
	}
}
