//go:build windows

package pluginhost

import (
	"os/exec"
)

func setSysProcAttr(cmd *exec.Cmd) {
	// Options for process isolation on Windows would go here
	// e.g. CREATE_NEW_PROCESS_GROUP
}
