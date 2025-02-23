package object

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"opensvc.com/opensvc/util/file"
)

// lockName is the path of the file to use as an action lock.
func (t *Node) frozenFile() string {
	return filepath.Join(t.VarDir(), "frozen")
}

// Frozen returns the unix timestamp of the last freeze.
func (t *Node) Frozen() time.Time {
	p := t.frozenFile()
	fi, err := os.Stat(p)
	if err != nil {
		return time.Time{}
	}
	return fi.ModTime()
}

//
// Freeze creates a persistant flag file that prevents orchestration
// of the object instance.
//
func (t *Node) Freeze() error {
	p := t.frozenFile()
	if file.Exists(p) {
		return nil
	}
	d := filepath.Dir(p)
	if !file.Exists(d) {
		if err := os.MkdirAll(d, os.ModePerm); err != nil {
			return err
		}
	}
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	f.Close()
	log.Info().Msg("now frozen")
	return nil
}

//
// Unfreeze removes the persistant flag file that prevents orchestration
// of the object instance.
//
func (t *Node) Unfreeze() error {
	p := t.frozenFile()
	if !file.Exists(p) {
		return nil
	}
	err := os.Remove(p)
	if err != nil {
		return err
	}
	log.Info().Msg("now unfrozen")
	return nil
}
