package object

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"opensvc.com/opensvc/core/instance"
	"opensvc.com/opensvc/testhelper"
)

func Test_Instance_States_Render(t *testing.T) {
	testhelper.Setup(t)
	cases := []string{"instanceStatus"}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {

			b, err := os.ReadFile(filepath.Join("testdata", name+".json"))
			require.Nil(t, err)

			var instanceStatus instance.Status
			err = json.Unmarshal(b, &instanceStatus)
			require.Nil(t, err)
			var timeZero time.Time
			instanceState := instance.States{
				Node:    instance.Node{Name: "node1", Frozen: timeZero},
				Status:  instanceStatus,
				Monitor: instance.Monitor{State: instance.MonitorStateIdle, StateUpdated: time.Now()},
				Config:  instance.Config{Priority: 50},
			}
			goldenFile := filepath.Join("testdata", name+".render")
			s := instanceState.Render()

			if *update {
				//
				t.Logf("updating golden file %s with current result", goldenFile)
				err = os.WriteFile(goldenFile, []byte(s), 0644)
				require.Nil(t, err)
			}
			expected, err := os.ReadFile(goldenFile)
			require.Nil(t, err)

			require.Equalf(t, string(expected), s, "found: \n%s", s)
		})
	}
}
