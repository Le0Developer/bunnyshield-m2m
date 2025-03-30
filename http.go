package bunnyshieldm2m

import (
	"errors"
	"io"
	"net/http"
	"time"
)

var ErrChallengeFailed = errors.New("bunnyshield: challenge failed")
var ErrServerUnexpectedFailure = errors.New("bunnyshield: server returned unexpected status code")

func SolveResponse(response *http.Response, cfgs ...HTTPSolverConfig) ([]*http.Cookie, error) {
	var cfg HTTPSolverConfig
	if len(cfgs) > 1 {
		return nil, ErrTooManyConfigs
	} else if len(cfgs) == 1 {
		cfg = cfgs[0]
	} else {
		cfg = DefaultHTTPSolverConfig()
	}

	var html []byte

	if cfg.Response == nil {
		var err error
		html, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	} else {
		html = cfg.Response
	}

	challenge, err := ParseChallengeFromHTML(string(html))
	if err != nil {
		return nil, err
	}

	start := time.Now()

	answer, err := challenge.Solve()
	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	if elapsed < cfg.Delay {
		time.Sleep(cfg.Delay - elapsed)
	}

	solveURL := response.Request.URL
	solveURL.Path = "/.bunny-shield/verify-pow"

	req, err := http.NewRequest("POST", solveURL.String(), nil)
	if err != nil {
		return nil, err
	}

	// this is a lie since there is no body, but that's what the javascript does
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("BunnyShield-Challenge-Response", answer.String())

	res, err := cfg.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 500 {
		return nil, ErrServerUnexpectedFailure
	} else if res.StatusCode >= 400 {
		return nil, ErrChallengeFailed
	}

	return res.Cookies(), nil
}

type HTTPSolverConfig struct {
	Response    []byte
	Client      http.Client
	SolveConfig SolveConfig
	Delay       time.Duration
}

func DefaultHTTPSolverConfig() HTTPSolverConfig {
	return HTTPSolverConfig{
		Client:      http.Client{},
		SolveConfig: DefaultSolveConfig(),
		// this minimum delay is present in the Javascript but not enforced by the server
		// it's probably a silent flag which will get you monitored if you ignore it
		Delay: 3500 * time.Millisecond,
	}
}
