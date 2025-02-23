package rescontainerdocker

import (
	"os/exec"

	"opensvc.com/opensvc/util/capabilities"
)

func init() {
	capabilities.Register(capabilitiesScanner)
}

func capabilitiesScanner() ([]string, error) {
	l := make([]string, 0)
	drvCap := drvID.Cap()
	if _, err := exec.LookPath("docker"); err != nil {
		return l, nil
	}
	l = append(l, drvCap)
	l = append(l, drvCap+".registry_creds")
	l = append(l, drvCap+".signal")
	l = append(l, altDrvID.Cap())
	return l, nil
}
