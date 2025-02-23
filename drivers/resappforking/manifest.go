package resappforking

import (
	"opensvc.com/opensvc/core/driver"
	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/core/manifest"
	"opensvc.com/opensvc/drivers/resapp"
)

var (
	drvID = driver.NewID(driver.GroupApp, "forking")
)

func init() {
	driver.Register(drvID, New)
}

// Manifest ...
func (t T) Manifest() *manifest.T {
	var keywordL []keywords.Keyword
	keywordL = append(keywordL, resapp.BaseKeywords...)
	keywordL = append(keywordL, resapp.UnixKeywords...)
	keywordL = append(keywordL, Keywords...)
	m := manifest.New(drvID, &t)
	m.AddContext([]manifest.Context{
		{
			Key:  "path",
			Attr: "Path",
			Ref:  "object.path",
		},
		{
			Key:  "nodes",
			Attr: "Nodes",
			Ref:  "object.nodes",
		},
		{
			Key:  "objectID",
			Attr: "ObjectID",
			Ref:  "object.id",
		},
	}...)
	m.AddKeyword(keywordL...)
	return m
}
