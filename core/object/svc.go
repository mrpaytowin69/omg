package object

import (
	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/core/path"
	"opensvc.com/opensvc/util/funcopt"
	"opensvc.com/opensvc/util/key"
)

type (
	svc struct {
		actor
	}

	//
	// Svc is the svc-kind object.
	//
	// These objects contain front facing resources like app and containers.
	//
	Svc interface {
		Actor
	}
)

// NewSvc allocates a svc kind object.
func NewSvc(p path.T, opts ...funcopt.O) (*svc, error) {
	s := &svc{}
	err := s.init(s, p, opts...)
	return s, err
}

func (t svc) KeywordLookup(k key.T, sectionType string) keywords.Keyword {
	return keywordLookup(keywordStore, k, t.path.Kind, sectionType)
}
