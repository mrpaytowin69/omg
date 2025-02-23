package status

import (
	"bytes"
	"encoding/json"
	"strings"
)

type (
	// T representing a Resource, Object Instance or Object status
	T int

	// L is a list of T
	L []T
)

const (
	// Undef means Undefined
	Undef T = 0
)

const (
	// NotApplicable means Not Applicable
	NotApplicable T = 1 << iota
	// Up means Configured or Active
	Up
	// Down means Unconfigured or Inactive
	Down
	// Warn means Partially configured or active
	Warn
	// StandbyUp means Instance with standby resources Configured or Active and no other resources
	StandbyUp
	// StandbyDown means Instance with standby resources Unconfigured or Inactive and no other resources
	StandbyDown
	// StandbyUpWithUp means Instance with standby resources Configured or Active and other resources Up
	StandbyUpWithUp
	// StandbyUpWithDown means Instance with standby resources Configured or Active and other resources Down
	StandbyUpWithDown
)

var toString = map[T]string{
	Up:                "up",
	Down:              "down",
	Warn:              "warn",
	NotApplicable:     "n/a",
	Undef:             "undef",
	StandbyUp:         "stdby up",
	StandbyDown:       "stdby down",
	StandbyUpWithUp:   "up",
	StandbyUpWithDown: "stdby up",
}

var toID = map[string]T{
	"up":         Up,
	"down":       Down,
	"warn":       Warn,
	"n/a":        NotApplicable,
	"undef":      Undef,
	"stdby up":   StandbyUp,
	"stdby down": StandbyDown,
}

func Parse(s string) T {
	if st, ok := toID[strings.TrimSpace(s)]; ok {
		return st
	}
	return Undef
}

func (t T) String() string {
	return toString[t]
}

func (t T) Is(l ...T) bool {
	for _, s := range l {
		if t == s {
			return true
		}
	}
	return false
}

// MarshalJSON marshals the enum as a quoted json string
func (t T) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *T) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*t = toID[j]
	return nil
}

// Add merges two states and returns the aggregate state.
func (t *T) Add(o T) {
	// handle invariants
	if o == Undef {
		return
	}
	if *t == Undef {
		*t = o
		return
	}
	if o == NotApplicable {
		return
	}
	if *t == NotApplicable {
		*t = o
		return
	}

	// other merges
	switch *t | o {
	case Up | Up:
		*t = Up
	case Up | Down:
		*t = Warn
	case Up | Warn:
		*t = Warn
	case Up | StandbyUp:
		*t = StandbyUpWithUp
	case Up | StandbyDown:
		*t = Warn
	case Up | StandbyUpWithUp:
		*t = StandbyUpWithUp
	case Up | StandbyUpWithDown:
		*t = Warn
	case Down | Down:
		*t = Down
	case Down | Warn:
		*t = Warn
	case Down | StandbyUp:
		*t = StandbyUpWithDown
	case Down | StandbyDown:
		*t = StandbyDown
	case Down | StandbyUpWithUp:
		*t = Warn
	case Down | StandbyUpWithDown:
		*t = StandbyUpWithDown
	case Warn | Warn:
		*t = Warn
	case Warn | StandbyUp:
		*t = Warn
	case Warn | StandbyDown:
		*t = Warn
	case Warn | StandbyUpWithUp:
		*t = Warn
	case Warn | StandbyUpWithDown:
		*t = Warn
	case StandbyUp | StandbyUp:
		*t = StandbyUp
	case StandbyUp | StandbyDown:
		*t = Warn
	case StandbyUp | StandbyUpWithUp:
		*t = StandbyUpWithUp
	case StandbyUp | StandbyUpWithDown:
		*t = StandbyUpWithDown
	case StandbyDown | StandbyDown:
		*t = StandbyDown
	case StandbyDown | StandbyUpWithUp:
		*t = Warn
	case StandbyDown | StandbyUpWithDown:
		*t = Warn
	case StandbyUpWithUp | StandbyUpWithDown:
		*t = Warn
	case StandbyUpWithUp | StandbyUpWithUp:
		*t = StandbyUpWithUp
	case StandbyUpWithDown | StandbyUpWithDown:
		*t = StandbyUpWithDown
	}
}

func (t L) Has(s T) bool {
	for _, e := range t {
		if e == s {
			return true
		}
	}
	return false
}

func List(l ...T) L {
	return L(l)
}

func (t L) Add(s ...T) L {
	return L(append(t, s...))
}

func (t L) String() string {
	l := make([]string, 0)
	for _, s := range t {
		l = append(l, s.String())
	}
	return strings.Join(l, ",")
}
