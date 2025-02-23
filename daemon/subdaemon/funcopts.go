package subdaemon

import (
	"opensvc.com/opensvc/daemon/routinehelper"
	"opensvc.com/opensvc/util/funcopt"
)

func WithName(name string) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.name = name
		return nil
	})
}

func WithMainManager(mgr Manager) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.main = mgr
		return nil
	})
}

func WithRoutineTracer(o routinehelper.Tracer) funcopt.O {
	return funcopt.F(func(i interface{}) error {
		t := i.(*T)
		t.SetTracer(o)
		return nil
	})
}
