package smon

import (
	"opensvc.com/opensvc/core/status"
)

func (o *smon) orchestratePlacedAt(dst string) {
	dstNodes := o.parsePlacedAtDestination(dst)
	if dstNodes.Contains(o.localhost) {
		o.orchestratePlacedStart()
	} else {
		o.orchestratePlacedStop()
	}
}

func (o *smon) placedUnfreeze() {
	if o.instStatus[o.localhost].IsThawed() {
		o.transitionTo(statusThawed)
	} else {
		o.doUnfreeze()
	}

}

func (o *smon) doPlacedStart() {
	o.doAction(o.crmStart, statusStarting, statusStarted, statusStartFailed)
}

func (o *smon) placedStart() {
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

func (o *smon) placedStop() {
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

func (o *smon) doPlacedStop() {
	o.createPendingWithDuration(stopDuration)
	o.doAction(o.crmStop, statusStopping, statusStopped, statusStopFailed)
}

func (o *smon) skipPlacedStop() {
	o.loggerWithState().Info().Msg("instance is already down")
	o.change = true
	o.state.Status = statusStopped
	o.clearPending()
}

func (o *smon) skipPlacedStart() {
	o.loggerWithState().Info().Msg("instance is already up")
	o.change = true
	o.state.Status = statusStarted
	o.clearPending()
}

func (o *smon) clearStopFailedIfDown() {
	instStatus := o.instStatus[o.localhost]
	switch instStatus.Avail {
	case status.Down, status.StandbyDown:
		o.loggerWithState().Info().Msg("instance is down, clear stop failed")
		o.change = true
		o.state.Status = statusStopped
		o.clearPending()
	}
}

func (o *smon) clearStoppedIfAggUp() {
	switch o.svcAgg.Avail {
	case status.Up:
		o.clearStopped()
	}
}

func (o *smon) clearStopped() {
	o.loggerWithState().Info().Msg("agg status is up, unset global expect")
	o.change = true
	o.state.GlobalExpect = globalExpectUnset
	if o.state.LocalExpect != statusIdle {
		o.state.LocalExpect = statusIdle
	}
	if o.state.Status != statusIdle {
		o.state.Status = statusIdle
	}
	o.clearPending()
}

func (o *smon) orchestrateFailoverPlacedStart() {
	switch o.state.Status {
	case statusIdle:
		o.placedUnfreeze()
	case statusThawed:
		o.orchestrateFailoverPlacedStartFromThawed()
	case statusStarted:
		o.orchestrateFailoverPlacedStartFromStarted()
	case statusStopped:
		o.orchestrateFailoverPlacedStartFromStopped()
	}
}

func (o *smon) orchestrateFailoverPlacedStartFromThawed() {
	o.transitionTo(statusStopped)
}

func (o *smon) orchestrateFailoverPlacedStartFromStopped() {
	switch o.svcAgg.Avail {
	case status.Down:
	default:
		return
	}
	o.placedStart()
}

func (o *smon) orchestrateFailoverPlacedStartFromStarted() {
	o.startedClearIfReached()
}

func (o *smon) orchestrateFlexPlacedStart() {
	switch o.state.Status {
	case statusIdle:
		o.placedUnfreeze()
	case statusThawed:
		o.orchestrateFlexPlacedStartFromThawed()
	case statusStarted:
		o.orchestrateFlexPlacedStartFromStarted()
	case statusStopped:
		o.transitionTo(statusIdle)
	}
}

func (o *smon) orchestrateFlexPlacedStartFromThawed() {
	o.placedStart()
}

func (o *smon) orchestrateFlexPlacedStartFromStarted() {
	o.startedClearIfReached()
}

func (o *smon) orchestrateFailoverPlacedStop() {
	switch o.state.Status {
	case statusIdle:
		o.placedUnfreeze()
	case statusThawed:
		o.placedStop()
	case statusStopFailed:
		o.clearStopFailedIfDown()
	case statusStopped:
		o.clearStoppedIfAggUp()
	case statusReady:
		o.transitionTo(statusIdle)
	case statusStartFailed:
		o.transitionTo(statusIdle)
	case statusThawing:
	case statusFreezing:
	case statusStopping:
	default:
		o.log.Error().Msgf("don't know how to orchestrate placed stopped from %s", o.state.Status)
	}
}

func (o *smon) orchestrateFlexPlacedStop() {
	switch o.state.Status {
	case statusIdle:
		o.placedUnfreeze()
	case statusThawed:
		o.placedStop()
	case statusStopFailed:
		o.clearStopFailedIfDown()
	case statusStopped:
		o.clearStoppedIfAggUp()
	case statusReady:
		o.transitionTo(statusIdle)
	case statusStartFailed:
		o.transitionTo(statusIdle)
	case statusThawing:
	case statusFreezing:
	case statusStopping:
	default:
		o.log.Error().Msgf("don't know how to orchestrate placed stopped from %s", o.state.Status)
	}
}
