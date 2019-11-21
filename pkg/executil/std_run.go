package executil

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fatih/color"
)

func StdRun(cmd *exec.Cmd) {
	cwd, _ := os.Getwd()

	envVars := EnvVars{}
	envVars.LoadFromEnviron()

	fmt.Fprintf(os.Stdout, "%s %s\n", color.CyanString(path.Join(cwd, cmd.Dir)), envVars.Parse(strings.Join(cmd.Args, " ")))
	{
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		go scanAndStdout(bufio.NewScanner(stdoutPipe))
	}
	{
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		go scanAndStderr(bufio.NewScanner(stderrPipe))
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(">> %s", err.Error()))
		os.Exit(1)
	}
}

func scanAndStdout(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Fprintln(os.Stdout, color.GreenString(">> %s", scanner.Text()))
	}
}

func scanAndStderr(scanner *bufio.Scanner) {
	for scanner.Scan() {
		fmt.Fprintln(os.Stderr, color.RedString(">> %s", scanner.Text()))
	}
}
