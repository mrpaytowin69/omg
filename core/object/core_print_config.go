package object

import (
	"opensvc.com/opensvc/core/rawconfig"
)

// PrintConfig gets a keyword value
func (t *core) PrintConfig() (rawconfig.T, error) {
	return t.config.Raw(), nil
}

func (t *core) EvalConfig() (rawconfig.T, error) {
	if actor, ok := t.config.Referrer.(Actor); ok {
		// required to eval references like {<rid>.exposed_devs}
		actor.ConfigureResources()
	}
	return t.config.RawEvaluated()
}

func (t *core) EvalConfigAs(nodename string) (rawconfig.T, error) {
	if actor, ok := t.config.Referrer.(Actor); ok {
		// required to eval references like {<rid>.exposed_devs}
		actor.ConfigureResources()
	}
	return t.config.RawEvaluatedAs(nodename)
}
