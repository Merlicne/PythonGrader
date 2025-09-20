package executer

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

type PythonExecutor struct{}

func (p *PythonExecutor) Execute(ctx context.Context, code string, stdin string) (string, error) {
	cmd := exec.CommandContext(ctx, "python3", "-c", code)

	var outMsgBytes bytes.Buffer
	var errMsgBytes bytes.Buffer
	cmd.Stdout = &outMsgBytes
	cmd.Stderr = &errMsgBytes

	if stdin != "" {
		cmd.Stdin = bytes.NewBufferString(stdin)
	}
	if err := cmd.Start(); err != nil {
		return "", parseCodeError(&errMsgBytes, err)
	}

	if err := cmd.Wait(); err != nil {
		return "", parseCodeError(&errMsgBytes, err)
	}
	return outMsgBytes.String(), nil
}

func parseCodeError(outputBytes *bytes.Buffer, err error) error {
	output := outputBytes.String()
	if len(output) == 0 {
		return err
	}

	return errors.New(output)
}

// check version of python
func (p *PythonExecutor) Version() (string, error) {
	cmd := exec.Command("python3", "--version")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}
