// vpath is a helper package easing the expansion of a virtual path like
// vol1/etc/nginx.conf to a host path like
// /srv/svc1data.ns1.vol.clu1/etc/nginx.conf
package vpath

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"opensvc.com/opensvc/core/kind"
	"opensvc.com/opensvc/core/object"
	"opensvc.com/opensvc/core/path"
	"opensvc.com/opensvc/core/status"
	"opensvc.com/opensvc/util/file"
	"opensvc.com/opensvc/util/loop"
)

var (
	ErrAccess = errors.New("vol is not accessible")
)

//
// Realpath expand a volume-relative path to a host full path.
//
// Example:
//
// INPUT        VOL     OUPUT            COMMENT
// /path                /path            host full path
// myvol/path   myvol   /srv/myvol/path  vol head relative path
//
func HostPath(s string, namespace string) (string, error) {
	var volRelativeSourcePath string
	l := strings.SplitN(s, "/", 2)
	if len(l[0]) == 0 {
		return s, nil
	}
	if len(l) == 2 {
		volRelativeSourcePath = l[1]
	}
	volPath := path.T{
		Name:      l[0],
		Namespace: namespace,
		Kind:      kind.Vol,
	}
	vol, err := object.NewVol(volPath)
	if err != nil {
		return s, err
	}
	if !vol.Path().Exists() {
		return s, errors.Errorf("%s does not exist", vol.Path())
	}
	st, err := vol.Status(context.Background())
	if err != nil {
		return s, err
	}
	switch st.Avail {
	case status.Up, status.NotApplicable, status.StandbyUp:
	default:
		return s, errors.Wrapf(ErrAccess, "%s(%s)", volPath, st.Avail)
	}
	return vol.Head() + "/" + volRelativeSourcePath, nil
}

// HostPaths applies the HostPath function to each path of the input list
func HostPaths(l []string, namespace string) ([]string, error) {
	for i, s := range l {
		if s2, err := HostPath(s, namespace); err != nil {
			return l, err
		} else {
			l[i] = s2
		}
	}
	return l, nil
}

//
// translation rules:
// INPUT        VOL     OUPUT       COMMENT
// /path                /dev/sda1   loop dev
// /dev/sda1            /dev/sda1   host full path
// myvol        myvol   /dev/sda1   vol dev path in host
//
func HostDevpath(s string, namespace string) (string, error) {
	if strings.HasPrefix(s, "/dev/") {
		return s, nil
	} else if file.ExistsAndRegular(s) {
		if lo, err := loop.New().FileGet(s); err != nil {
			return "", err
		} else {
			return lo.Name, nil
		}
		return s, nil
	} else {
		// volume device
		volPath := path.T{
			Name:      s,
			Namespace: namespace,
			Kind:      kind.Vol,
		}
		vol, err := object.NewVol(volPath)
		if err != nil {
			return s, err
		}
		st, err := vol.Status(context.Background())
		if err != nil {
			return s, err
		}
		switch st.Avail {
		case status.Up, status.NotApplicable, status.StandbyUp:
		default:
			return s, errors.Wrapf(ErrAccess, "%s(%s)", volPath, st.Avail)
		}
		dev := vol.Device()
		if dev == nil {
			return s, errors.Errorf("%s is not a device-capable vol", s)
		}
		return dev.Path(), nil
	}
}

// HostDevpaths applies the HostDevpath function to each path of the input list
func HostDevpaths(l []string, namespace string) ([]string, error) {
	for i, s := range l {
		if s2, err := HostDevpath(s, namespace); err != nil {
			return l, err
		} else {
			l[i] = s2
		}
	}
	return l, nil
}
