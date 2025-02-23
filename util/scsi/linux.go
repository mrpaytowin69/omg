//go:build linux

package scsi

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opensvc/fcntllock"
	"github.com/opensvc/flock"

	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/util/capabilities"
	"opensvc.com/opensvc/util/xerrors"
	"opensvc.com/opensvc/util/xsession"
)

func (t *PersistentReservationHandle) setup() error {
	if t.persistentReservationDriver != nil {
		return nil
	}
	if capabilities.Has(MpathPersistCapability) {
		t.persistentReservationDriver = MpathPersistDriver{
			Log: t.Log,
		}
	} else if capabilities.Has(SGPersistCapability) {
		t.persistentReservationDriver = SGPersistDriver{
			Log: t.Log,
		}
	} else {
		return ErrNotSupported
	}
	return nil
}

func doWithLock(timeout time.Duration, name, intent string, f func() error) error {
	p := filepath.Join(rawconfig.Paths.Lock, strings.Join([]string{"scsi", name}, "."))
	lock := flock.New(p, xsession.ID, fcntllock.New)
	err := lock.Lock(timeout, intent)
	if err != nil {
		return err
	}
	defer func() { _ = lock.UnLock() }()
	return f()
}

func ListHostDirs() ([]string, error) {
	dirs, err := filepath.Glob("/sys/class/scsi_host/host*")
	if err != nil {
		return []string{}, err
	}
	return dirs, nil
}

func LockedScanAll(timeout time.Duration) error {
	return doWithLock(timeout, "scan", "all", func() error {
		return ScanAll()
	})
}

func ScanAll() error {
	return ScanAllBusTargetLun("-", "-", "-")
}

func ScanAllBusTargetLun(b, t, l string) error {
	hosts, err := ListHostDirs()
	if err != nil {
		return err
	}
	for _, h := range hosts {
		if e := ScanHostDirBusTargetLun(h, b, t, l); err != nil {
			err = xerrors.Append(err, e)
		}
	}
	return err
}

func ScanHostDirBusTargetLun(h, b, t, l string) error {
	filename := fmt.Sprintf("%s/scan", h)
	s := fmt.Sprintf("%s %s %s", b, t, l)
	return ioutil.WriteFile(filename, []byte(s), os.ModePerm)
}

func ScanHostBusTargetLun(h, b, t, l string) error {
	filename := fmt.Sprintf("/sys/class/scsi_host/host%s/scan", h)
	s := fmt.Sprintf("%s %s %s", b, t, l)
	return ioutil.WriteFile(filename, []byte(s), os.ModePerm)
}
