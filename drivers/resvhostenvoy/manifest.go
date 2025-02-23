package resvhostenvoy

import (
	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/core/manifest"
	"opensvc.com/opensvc/util/converters"
)

var (
	drvID = driver.NewID(driver.GroupVhost, "envoy")
)

func init() {
	driver.Register(drvID, New)
}

// Manifest exposes to the core the input expected by the driver.
func (t T) Manifest() *manifest.T {
	m := manifest.New(drvID, t)
	m.AddKeyword([]keywords.Keyword{
		{
			Option:    "domains",
			Attr:      "Domains",
			Scopable:  true,
			Converter: converters.List,
			Default:   "{name}",
			Example:   "{name}",
			Text:      "The list of http domains in this expose.",
		},
		{
			Option:    "routes",
			Attr:      "Routes",
			Scopable:  true,
			Converter: converters.List,
			Example:   "route#1 route#2",
			Text:      "The list of route resource identifiers for this vhost.",
		},
	}...)
	return m
}
