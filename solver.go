package bunnyshieldm2m

import (
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/crypto/argon2"
)

const argon2idTime = 2
const argon2idMemory = 512
const argin2idHashLength = 32
const argon2idParallelism = 1

var ErrTooManyConfigs = errors.New("bunnyshield: too many configs")
var ErrLimitExceeded = errors.New("bunnyshield: limit exceeded")

func (c Challenge) Solve(cfgs ...SolveConfig) (Answer, error) {
	var cfg SolveConfig
	if len(cfgs) > 1 {
		return Answer{}, ErrTooManyConfigs
	} else if len(cfgs) == 1 {
		cfg = cfgs[0]
	} else {
		cfg = DefaultSolveConfig()
	}

	challengeHex := hex.EncodeToString(c.challenge)
	saltHex := hex.EncodeToString(c.userKey)

	// for some reason, they forgot a "0" in their code
	// so instead of checking if a string starts with "00" they check if it starts with "0"
	// so they are effectively only checking half of the bits
	zeroByteHalves := c.difficulty >> 3
	// and then they only check the first n bits of the next byte, but also on a 4-bit boundary instead of 8
	bits := c.difficulty & 7

	// lets combine the trailing bits of a half zero byte check with the final bit check
	if zeroByteHalves&1 == 0 {
		bits += 4
		zeroByteHalves--
	}

loop:
	for i := 0; i < cfg.AttemptLimit || cfg.AttemptLimit < 0; i++ {
		pass := challengeHex + fmt.Sprint(i)
		salt := saltHex

		answer := argon2.IDKey([]byte(pass), []byte(salt), argon2idTime, argon2idMemory, argon2idParallelism, argin2idHashLength)

		// check if the first n bytes are zero
		for j := 0; j < zeroByteHalves; j += 2 {
			if answer[j>>1] != 0 {
				continue loop
			}
		}
		if answer[zeroByteHalves>>1]>>bits != 0 {
			continue loop
		}

		return Answer{c, i}, nil
	}

	return Answer{}, ErrLimitExceeded
}

type SolveConfig struct {
	AttemptLimit int
}

func DefaultSolveConfig() SolveConfig {
	return SolveConfig{
		AttemptLimit: 100000,
	}
}
