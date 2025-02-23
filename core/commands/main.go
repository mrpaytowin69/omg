package commands

func mergeSelector(selector string, subsysSelector string, kind string, defaultSelector string) string {
	var s string
	switch {
	case selector != "":
		s = selector
	case subsysSelector != "":
		s = subsysSelector
	default:
		s = defaultSelector
	}
	if kind != "" {
		kindSelector := "*/" + kind + "/*"
		if s == "" {
			s = kindSelector
		} else {
			s += "+" + kindSelector
		}
	}
	return s
}
