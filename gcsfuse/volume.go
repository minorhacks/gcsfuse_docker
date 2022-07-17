package gcsfuse

import (
	"fmt"
	"os/exec"
	"sync"
)

type Volume struct {
	// Bucket to mount
	bucket string
	// Path at which to mount
	hostPath string

	// Protects the vars in this block
	mu sync.Mutex
	// If true, the volume is currently mounted
	mounted bool
}

func (v *Volume) Mount(credsPath string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Check to see if the volume is already mounted
	if v.mounted {
		return fmt.Errorf("volume is alredy mounted")
	}

	// Run gcsfuse to mount the volume
	// gcsfuse is ran in daemon mode (without --foreground) and it is assumed
	// the volume is mounted if it returns successfully.
	gcsfuseFlags := []string{
		"-o=allow_other",
		"--key-file",
		credsPath,
	}
	cmd := exec.Command("gcsfuse", gcsfuseFlags...)
	cmd.Args = append(cmd.Args, v.bucket, v.hostPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gcsfuse failed. output:\n%s", output)
	}

	v.mounted = true
	return nil
}
