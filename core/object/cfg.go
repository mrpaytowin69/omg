package object

import (
	"encoding/base64"
	"strings"
	"unicode"

	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/util/funcopt"
	"opensvc.com/opensvc/util/key"
)

type (
	cfg struct {
		keystore
	}

	//
	// Cfg is the cfg-kind object.
	//
	// These objects are unencrypted key-value store.
	// Values can be binary or text.
	//
	// A Key can be installed as a file in a Vol, then exposed to apps
	// and containers.
	// A Key can be exposed as a environment variable for apps and
	// containers.
	// A Signal can be sent to consumer processes upon exposed key value
	// changes.
	//
	Cfg interface {
		Keystore
	}
)

// NewCfg allocates a cfg kind object.
func NewCfg(p any, opts ...funcopt.O) (*cfg, error) {
	s := &cfg{}
	s.customEncode = cfgEncode
	s.customDecode = cfgDecode
	if err := s.init(s, p, opts...); err != nil {
		return s, err
	}
	s.Config().RegisterPostCommit(s.postCommit)
	return s, nil
}

func (t cfg) KeywordLookup(k key.T, sectionType string) keywords.Keyword {
	return keywordLookup(keywordStore, k, t.path.Kind, sectionType)
}

func cfgEncode(b []byte) (string, error) {
	switch {
	case isAsciiPrintable(b):
		return "literal:" + string(b), nil
	default:
		return "base64:" + base64.URLEncoding.Strict().EncodeToString(b), nil
	}
}

func cfgDecode(s string) ([]byte, error) {
	switch {
	case strings.HasPrefix(s, "base64:"):
		return base64.URLEncoding.DecodeString(s[7:])
	case strings.HasPrefix(s, "literal:"):
		return []byte(s[8:]), nil
	default:
		return []byte(s), nil
	}
}

func isAsciiPrintable(bytes []byte) bool {
	for _, b := range bytes {
		r := rune(b)
		if r > unicode.MaxASCII || !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
