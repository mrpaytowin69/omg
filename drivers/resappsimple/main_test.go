package resappsimple

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"opensvc.com/opensvc/core/actionrollback"
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/drivers/resapp"
	"opensvc.com/opensvc/util/file"
	"opensvc.com/opensvc/util/pg"
)

var (
	log = zerolog.New(os.Stdout).With().Timestamp().Logger()
)

func prepareConfig(t *testing.T) (td string, cleanup func()) {
	td = t.TempDir()
	rawconfig.Load(map[string]string{"osvc_root_path": td})
	cleanup = func() {
		rawconfig.Load(map[string]string{})
	}
	return
}

func getActionContext() (ctx context.Context, cancel context.CancelFunc) {
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	ctx = actionrollback.NewContext(ctx)
	return
}
func WithLoggerAndPgApp(app T) T {
	app.SetLoggerForTest(log)
	app.SetRID("foo")
	app.SetPG(&pg.Config{})
	return app
}

func TestStart(t *testing.T) {
	startReturnMsg := "Start(...) returned value"

	t.Run("execute start command", func(t *testing.T) {
		td, cleanup := prepareConfig(t)
		defer cleanup()

		filename := filepath.Join(td, "trace")
		app := WithLoggerAndPgApp(T{T: resapp.T{StartCmd: "touch " + filename}})

		ctx, cancel := getActionContext()
		defer cancel()
		assert.Nil(t, app.Start(ctx), startReturnMsg)
		time.Sleep(20 * time.Millisecond) // give time for file start cmd does its job
		assert.True(t, file.Exists(filename), "missing start cmd !")
	})

	t.Run("does not execute start command if status is already up", func(t *testing.T) {
		if os.Getpid() != 0 {
			t.Skip("skipped for non root user")
		}
		td, cleanup := prepareConfig(t)
		defer cleanup()
		createdFileFromStart := filepath.Join(td, "trace")
		app := WithLoggerAndPgApp(T{T: resapp.T{StartCmd: "touch " + createdFileFromStart, CheckCmd: "echo"}})
		ctx, cancel := getActionContext()
		defer cancel()
		assert.Nil(t, app.Start(ctx), startReturnMsg)
		assert.False(t, file.Exists(createdFileFromStart), "start cmd called !")
	})

	t.Run("when start succeed stop is added to rollback stack", func(t *testing.T) {
		if os.Getpid() != 0 {
			t.Skip("skipped for non root user")
		}
		td, cleanup := prepareConfig(t)
		defer cleanup()

		filename := filepath.Join(td, "trace")
		app := WithLoggerAndPgApp(
			T{T: resapp.T{
				StartCmd: "echo",
				CheckCmd: "exit 2",
				StopCmd:  "touch " + filename,
			}})
		ctx, cancel := getActionContext()
		defer cancel()
		assert.Nil(t, app.Start(ctx), startReturnMsg)
		assert.Nil(t, actionrollback.Rollback(ctx))
		assert.True(t, file.Exists(filename), "missing rollback stop cmd !")
	})

	t.Run("when start fails stop is not added to rollback stack", func(t *testing.T) {
		td, cleanup := prepareConfig(t)
		defer cleanup()

		filename := filepath.Join(td, "trace")
		app := WithLoggerAndPgApp(
			T{T: resapp.T{
				StartCmd: "noSuchAppTest",
				StopCmd:  "touch " + filename,
			}})
		ctx, cancel := getActionContext()
		defer cancel()
		assert.NotNil(t, app.Start(ctx), startReturnMsg)
		assert.Nil(t, actionrollback.Rollback(ctx))
		assert.False(t, file.Exists(filename), "rollback cmd called !")
	})

	t.Run("when already started stop is not added to rollback stack", func(t *testing.T) {
		if os.Getpid() != 0 {
			t.Skip("skipped for non root user")
		}
		td, cleanup := prepareConfig(t)
		defer cleanup()

		filename := filepath.Join(td, "trace")
		app := WithLoggerAndPgApp(
			T{T: resapp.T{
				StartCmd: "echo",
				CheckCmd: "echo",
				StopCmd:  "touch " + filename,
			}})
		ctx, cancel := getActionContext()
		defer cancel()
		assert.Nil(t, app.Start(ctx), startReturnMsg)
		assert.Nil(t, actionrollback.Rollback(ctx))
		assert.False(t, file.Exists(filename), "rollback cmd called !")
	})
}

func TestStop(t *testing.T) {
	if os.Getpid() != 0 {
		t.Skip("skipped for non root user")
	}
	t.Run("execute stop command", func(t *testing.T) {
		td, cleanup := prepareConfig(t)
		defer cleanup()

		filename := filepath.Join(td, "trace")
		app := WithLoggerAndPgApp(T{T: resapp.T{
			CheckCmd: "exit 2",
			StopCmd:  "touch " + filename,
		}})
		ctx, cancel := getActionContext()
		defer cancel()
		assert.Nil(t, app.Stop(ctx), "Stop(...) returned value")
		assert.True(t, file.Exists(filename))
	})

	t.Run("does not execute stop command if status is already down", func(t *testing.T) {
		td, cleanup := prepareConfig(t)
		defer cleanup()
		filename := filepath.Join(td, "succeed")
		app := WithLoggerAndPgApp(T{T: resapp.T{StopCmd: "touch " + filename, CheckCmd: "bash -c false"}})
		ctx, cancel := getActionContext()
		defer cancel()
		assert.Nil(t, app.Stop(ctx), "Stop(...) returned value")
		assert.False(t, file.Exists(filename), "stop cmd called !")
	})
}
