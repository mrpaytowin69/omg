package imon

import (
	"opensvc.com/opensvc/core/instance"
	"opensvc.com/opensvc/core/status"
	"opensvc.com/opensvc/util/stringslice"
)

func (o *imon) orchestrateFailoverPlacedStart() {
	switch o.state.State {
	case instance.MonitorStateIdle:
		o.placedUnfreeze()
	case instance.MonitorStateThawed:
		o.orchestrateFailoverPlacedStartFromThawed()
	case instance.MonitorStateStarted:
		o.orchestrateFailoverPlacedStartFromStarted()
	case instance.MonitorStateStopped:
		o.orchestrateFailoverPlacedStartFromStopped()
	case instance.MonitorStateStartFailed:
		o.orchestratePlacedFromStartFailed()
	case instance.MonitorStateThawing:
	case instance.MonitorStateFreezing:
	case instance.MonitorStateStopping:
	case instance.MonitorStateStarting:
	default:
		o.log.Error().Msgf("don't know how to orchestrate placed start from %s", o.state.State)
	}
}

func (o *imon) orchestrateFlexPlacedStart() {
	switch o.state.State {
	case instance.MonitorStateIdle:
		o.placedUnfreeze()
	case instance.MonitorStateThawed:
		o.orchestrateFlexPlacedStartFromThawed()
	case instance.MonitorStateStarted:
		o.orchestrateFlexPlacedStartFromStarted()
	case instance.MonitorStateStopped:
		o.transitionTo(instance.MonitorStateIdle)
	case instance.MonitorStateStartFailed:
		o.orchestratePlacedFromStartFailed()
	case instance.MonitorStateThawing:
	case instance.MonitorStateFreezing:
	case instance.MonitorStateStopping:
	case instance.MonitorStateStarting:
	default:
		o.log.Error().Msgf("don't know how to orchestrate placed start from %s", o.state.State)
	}
}

func (o *imon) orchestrateFailoverPlacedStop() {
	switch o.state.State {
	case instance.MonitorStateIdle:
		o.placedUnfreeze()
	case instance.MonitorStateThawed:
		o.placedStop()
	case instance.MonitorStateStopFailed:
		o.clearStopFailedIfDown()
	case instance.MonitorStateStopped:
		o.clearStoppedIfObjectStatusAvailUp()
	case instance.MonitorStateReady:
		o.transitionTo(instance.MonitorStateIdle)
	case instance.MonitorStateStartFailed:
		o.orchestratePlacedFromStartFailed()
	case instance.MonitorStateThawing:
	case instance.MonitorStateFreezing:
	case instance.MonitorStateStopping:
	case instance.MonitorStateStarting:
	default:
		o.log.Error().Msgf("don't know how to orchestrate placed stop from %s", o.state.State)
	}
}

func (o *imon) orchestrateFlexPlacedStop() {
	switch o.state.State {
	case instance.MonitorStateIdle:
		o.placedUnfreeze()
	case instance.MonitorStateThawed:
		o.placedStop()
	case instance.MonitorStateStopFailed:
		o.clearStopFailedIfDown()
	case instance.MonitorStateStopped:
		o.clearStoppedIfObjectStatusAvailUp()
	case instance.MonitorStateReady:
		o.transitionTo(instance.MonitorStateIdle)
	case instance.MonitorStateStartFailed:
		o.orchestratePlacedFromStartFailed()
	case instance.MonitorStateThawing:
	case instance.MonitorStateFreezing:
	case instance.MonitorStateStopping:
	case instance.MonitorStateStarting:
	default:
		o.log.Error().Msgf("don't know how to orchestrate placed stop from %s", o.state.State)
	}
}

func (o *imon) orchestratePlacedAt() {
	var dstNodes []string
	if options, ok := o.state.GlobalExpectOptions.(instance.MonitorGlobalExpectOptionsPlacedAt); ok {
		dstNodes = options.Destination
	} else {
		o.log.Error().Msgf("missing placed@ destination")
		return
	}
	if stringslice.Has(o.localhost, dstNodes) {
		o.orchestratePlacedStart()
	} else {
		o.orchestratePlacedStop()
	}
}

func (o *imon) placedUnfreeze() {
	if o.instStatus[o.localhost].IsThawed() {
		o.transitionTo(instance.MonitorStateThawed)
	} else {
		o.doUnfreeze()
	}
}

func (o *imon) doPlacedStart() {
	o.doAction(o.crmStart, instance.MonitorStateStarting, instance.MonitorStateStarted, instance.MonitorStateStartFailed)
}

func (o *imon) placedStart() {
	instStatus := o.instStatus[o.localhost]
	switch instStatus.Avail {
	case status.Down, status.StandbyDown, status.StandbyUp:
		o.doPlacedStart()
	case status.Up, status.Warn:
		o.skipPlacedStart()
	default:
		return
	}
}

func (o *imon) placedStop() {
	instStatus := o.instStatus[o.localhost]
	switch instStatus.Avail {
	case status.Down, status.StandbyDown, status.StandbyUp:
		o.skipPlacedStop()
	case status.Up, status.Warn:
		o.doPlacedStop()
	default:
		return
	}
}

func (o *imon) doPlacedStop() {
	o.createPendingWithDuration(stopDuration)
	o.doAction(o.crmStop, instance.MonitorStateStopping, instance.MonitorStateStopped, instance.MonitorStateStopFailed)
}

func (o *imon) skipPlacedStop() {
	o.loggerWithState().Info().Msg("instance is already down")
	o.change = true
	o.state.State = instance.MonitorStateStopped
	o.clearPending()
}

func (o *imon) skipPlacedStart() {
	o.loggerWithState().Info().Msg("instance is already up")
	o.change = true
	o.state.State = instance.MonitorStateStarted
	o.clearPending()
}

func (o *imon) clearStopFailedIfDown() {
	instStatus := o.instStatus[o.localhost]
	switch instStatus.Avail {
	case status.Down, status.StandbyDown:
		o.loggerWithState().Info().Msg("instance is down, clear stop failed")
		o.change = true
		o.state.State = instance.MonitorStateStopped
		o.clearPending()
	}
}

func (o *imon) clearStoppedIfObjectStatusAvailUp() {
	switch o.objStatus.Avail {
	case status.Up:
		o.clearStopped()
	}
}

func (o *imon) clearStopped() {
	o.loggerWithState().Info().Msg("object avail status is up, unset global expect")
	o.change = true
	o.state.GlobalExpect = instance.MonitorGlobalExpectUnset
	o.state.LocalExpect = instance.MonitorLocalExpectUnset
	o.state.State = instance.MonitorStateIdle
	o.clearPending()
}

func (o *imon) orchestrateFailoverPlacedStartFromThawed() {
	instStatus := o.instStatus[o.localhost]
	switch instStatus.Avail {
	case status.Up:
		o.transitionTo(instance.MonitorStateStarted)
	default:
		o.transitionTo(instance.MonitorStateStopped)
	}
}

func (o *imon) orchestrateFailoverPlacedStartFromStopped() {
	switch o.objStatus.Avail {
	case status.NotApplicable, status.Undef:
		o.startedClearIfReached()
	case status.Down:
		o.placedStart()
	default:
		return
	}
}

func (o *imon) orchestrateFailoverPlacedStartFromStarted() {
	o.startedClearIfReached()
}

func (o *imon) orchestrateFlexPlacedStartFromThawed() {
	o.placedStart()
}

func (o *imon) orchestrateFlexPlacedStartFromStarted() {
	o.startedClearIfReached()
}

func (o *imon) orchestratePlacedFromStartFailed() {
	switch {
	case o.AllInstanceMonitorState(instance.MonitorStateStartFailed):
		o.loggerWithState().Info().Msg("all instances are start failed, unset global expect")
		o.change = true
		o.state.GlobalExpect = instance.MonitorGlobalExpectUnset
		o.clearPending()
	case o.objStatus.Avail == status.Up:
		o.startedClearIfReached()
	}
}
