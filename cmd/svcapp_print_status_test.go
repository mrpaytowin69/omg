package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opensvc.com/opensvc/testhelper"
)

func TestAppPrintStatusFlatJson(t *testing.T) {
	type logT struct {
		Level   string
		Message string
	}
	cases := [][]logT{
		{
			{"info", "FOO"},
		},
		{},

		{
			{"warn", "DeadlineExceeded"},
		},
		{
			{"warn", "line1"},
			{"warn", "line2"},
		},
	}
	env := testhelper.Setup(t)
	env.InstallFile("../testdata/svcapp_print_status_status_log.conf", "etc/svcapp.conf")
	args := []string{"svcapp", "print", "status", "-r", "--format", "flat_json"}
	t.Logf("run 'om %v'", strings.Join(args, " "))
	cmd := exec.Command(os.Args[0], args...)
	cmd.Env = append(os.Environ(), "GO_TEST_MODE=off", "OSVC_ROOT_PATH="+env.Root)
	out, err := cmd.CombinedOutput()
	require.Nil(t, err, "got: \n%v", string(out))

	outS := string(out)
	for rid, c := range cases {
		t.Logf("check rid %d, expected %v", rid, c)
		for i, log := range c {
			prefix := fmt.Sprintf("instances.node1.status.resources[%d].log[%d]", rid, i)
			searched := fmt.Sprintf("%s.message = %s%s%s", prefix, string('"'), log.Message, string('"'))
			assert.Containsf(t, outS, searched, "%s not found in \n%s", searched, string(outS))

			searched = fmt.Sprintf("%s.level = %s%s%s", prefix, string('"'), log.Level, string('"'))
			assert.Containsf(t, outS, searched, "%s not found in \n%s", searched, string(outS))
		}
		mustNotExist := fmt.Sprintf("instances.node1.status.resources[%d].log[%d]", rid, len(c)+1)
		assert.NotContainsf(t, outS, mustNotExist, "extra log has been found: '%s' in \n'%s'", mustNotExist, outS)
	}
}
