//go:build linux

package resipcni

import "opensvc.com/opensvc/util/capabilities"

func init() {
	capabilities.Register(capabilitiesScanner)
}

func capabilitiesScanner() ([]string, error) {
	return []string{drvID.Cap()}, nil
}
