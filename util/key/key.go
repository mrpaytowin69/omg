package key

import "strings"

type (
	T struct {
		Section string `json:"section"`
		Option  string `json:"option"`
	}
	L []T
)

func New(section, option string) T {
	if section == "" && option != "" {
		section = "DEFAULT"
	}
	return T{
		Section: section,
		Option:  option,
	}
}

func ParseL(l []string) L {
	kws := make(L, len(l))
	for _, s := range l {
		kws = append(kws, Parse(s))
	}
	return kws
}

func Parse(s string) T {
	l := strings.SplitN(s, ".", 2)
	switch len(l) {
	case 1:
		if strings.Index(s, "#") >= 0 {
			return T{s, ""}
		}
		return T{"DEFAULT", s}
	case 2:
		return T{l[0], l[1]}
	default:
		return T{}
	}
}

func (t T) BaseOption() string {
	l := strings.SplitN(t.Option, "@", 2)
	return l[0]
}

func (t T) Scope() string {
	l := strings.Split(t.Option, "@")
	switch len(l) {
	case 2:
		return l[1]
	default:
		return ""
	}
}

func (t T) String() string {
	if t.Section == "DEFAULT" {
		return t.Option
	}
	if t.Option == "" {
		return t.Section
	}
	return t.Section + "." + t.Option
}
