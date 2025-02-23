package placement

import (
	"bytes"
	"encoding/json"

	"opensvc.com/opensvc/util/xmap"
)

type (
	// Policy is an integer representing the opensvc object placement policy.
	Policy int

	// State is optimal when the running <n> instances are placed on the <n> first
	// nodes of the candidate list sorted by the selected placement policy.
	State int
)

const (
	// Invalid is for invalid kinds
	Invalid Policy = iota
	// None is the policy used by special objects like sec, cfg, usr.
	None
	// NodesOrder is the policy where node priorities are assigned to nodes descending from left to right in the nodes list.
	NodesOrder
	// LoadAvg is the policy where node priorities are assigned to nodes based on load average. The higher the load, the lower the priority.
	LoadAvg
	// Shift is the policy where node priorities are assigned to nodes based on the scaler slice number.
	Shift
	// Spread is the policy where node priorities are assigned to nodes based on nodename hashing. The node priority is stable as long as the cluster members are stable.
	Spread
	// Score is the policy where node priorities are assigned to nodes based on score. The higher the score, the higher the priority.
	Score

	NotApplicable State = iota
	Optimal
	NonOptimal
)

var (
	policyToString = map[Policy]string{
		None:       "none",
		NodesOrder: "nodes order",
		LoadAvg:    "load avg",
		Shift:      "shift",
		Spread:     "spread",
		Score:      "score",
	}

	policyToID = map[string]Policy{
		"none":        None,
		"nodes order": NodesOrder,
		"load avg":    LoadAvg,
		"shift":       Shift,
		"spread":      Spread,
		"score":       Score,
	}

	stateToString = map[State]string{
		NotApplicable: "n/a",
		Optimal:       "optimal",
		NonOptimal:    "non-optimal",
	}

	stateToID = map[string]State{
		"n/a":         NotApplicable,
		"optimal":     Optimal,
		"non-optimal": NonOptimal,
	}
)

func (t State) String() string {
	return stateToString[t]
}

// MarshalJSON marshals the enum as a quoted json string
func (t State) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(stateToString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *State) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*t = stateToID[j]
	return nil
}

func (t Policy) String() string {
	return policyToString[t]
}

// NewPolicy returns a id from its string representation.
func NewPolicy(s string) Policy {
	t, ok := policyToID[s]
	if ok {
		return t
	}
	return Invalid
}

// MarshalJSON marshals the enum as a quoted json string
func (t Policy) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(policyToString[t])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (t *Policy) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	*t = policyToID[j]
	return nil
}

func PolicyNames() []string {
	return xmap.Keys(policyToID)
}
