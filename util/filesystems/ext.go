package filesystems

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog"
	"opensvc.com/opensvc/util/command"
)

func extCanFSCK() error {
	if _, err := exec.LookPath("e2fsck"); err != nil {
		return err
	}
	return nil
}

func extFSCK(s string) error {
	cmd := exec.Command("e2fsck", "-p", s)
	cmd.Start()
	cmd.Wait()
	exitCode := cmd.ProcessState.ExitCode()
	switch exitCode {
	case 0: // All good
		return nil
	case 1: // File system errors corrected
		return nil
	case 32: // E2fsck canceled by user request
		return nil
	case 33: // ?
		return nil
	default:
		return fmt.Errorf("%s exit code: %d", cmd, exitCode)
	}
}

func extIsFormated(s string) (bool, error) {
	if _, err := exec.LookPath("tune2fs"); err != nil {
		return false, errors.New("tune2fs not found")
	}
	cmd := exec.Command("tune2fs", "-l", s)
	cmd.Start()
	cmd.Wait()
	exitCode := cmd.ProcessState.ExitCode()
	switch exitCode {
	case 0: // All good
		return true, nil
	default:
		return false, nil
	}
}

func xMKFS(x string, s string, xargs []string, log *zerolog.Logger) error {
	if _, err := exec.LookPath(x); err != nil {
		return fmt.Errorf("%s not found", x)
	}
	args := []string{"-F", "-q", s}
	args = append(args, xargs...)
	cmd := command.New(
		command.WithName(x),
		command.WithArgs(args),
		command.WithLogger(log),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	return cmd.Run()
}
