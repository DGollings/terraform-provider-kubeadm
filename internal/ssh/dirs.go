package ssh

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/communicator"
	"github.com/hashicorp/terraform/terraform"
)

// DoMkdir creates a remote directory
func DoMkdir(path string) ApplyFunc {
	return ApplyFunc(func(o terraform.UIOutput, comm communicator.Communicator, useSudo bool) error {
		o.Output(fmt.Sprintf("Creating directory %s", path))
		mkdirCmd := fmt.Sprintf("mkdir -p %s", path)
		return DoExec(mkdirCmd).Apply(o, comm, useSudo)
	})
}

// CheckDirExists checks that a directory exists
func CheckDirExists(path string) CheckerFunc {
	return CheckerFunc(func(o terraform.UIOutput, comm communicator.Communicator, useSudo bool) (bool, error) {
		const markFound = "YES_IT_IS_THERE"
		found := false
		var interceptor OutputFunc = func(s string) {
			if strings.Contains(s, markFound) {
				found = true
			}
		}

		command := fmt.Sprintf("[ -d '%s' ] && echo '%s'", path, markFound)
		err := DoExec(command).Apply(interceptor, comm, useSudo)
		if err != nil {
			return false, err
		}

		return found, nil
	})
}
