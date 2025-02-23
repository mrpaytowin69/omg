package imon

import "opensvc.com/opensvc/core/instance"

func (o *imon) orchestrateAborted() {
	o.log.Info().Msg("abort orchestration: unset global expect")
	o.change = true
	o.state.GlobalExpect = instance.MonitorGlobalExpectUnset
	o.state.GlobalExpectOptions = nil
	o.updateIfChange()
}
