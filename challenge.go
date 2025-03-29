package bunnyshieldm2m

import (
	"encoding/hex"
	"errors"
	"strings"
)

type Challenge struct {
	userKey   []byte
	challenge []byte
	raw       string

	difficulty int
}

func (c Challenge) string() string {
	return c.raw
}

var ErrInvalidChallenge = errors.New("bunnyshield: invalid challenge")

func ParseChallenge(text string) (Challenge, error) {
	parts := strings.Split(text, "#")
	if len(parts) < 2 {
		return Challenge{}, ErrInvalidChallenge
	}

	userKey, err := hex.DecodeString(parts[0])
	if err != nil {
		return Challenge{}, errors.Join(ErrInvalidChallenge, err)
	}
	challenge, err := hex.DecodeString(parts[1])
	if err != nil {
		return Challenge{}, errors.Join(ErrInvalidChallenge, err)
	}

	// currently the difficulty seems to always be 13 bits
	// it's not part of the challenge string, but rather hardcoded in the javascript
	// maybe the JS is dynamic? I don't know yet
	return Challenge{userKey, challenge, text, 13}, nil
}

func ParseChallengeFromHTML(html string) (Challenge, error) {
	// the challenge is in the "data-pow" attribute of the "body" tag
	offset := strings.Index(html, `data-pow="`)
	if offset == -1 {
		return Challenge{}, ErrInvalidChallenge
	}
	offset += len(`data-pow="`)
	end := strings.Index(html[offset:], `"`)
	if end == -1 {
		return Challenge{}, ErrInvalidChallenge
	}

	return ParseChallenge(html[offset : offset+end])
}
