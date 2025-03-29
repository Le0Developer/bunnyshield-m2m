package main

import (
	"os"

	bunnyshieldm2m "github.com/le0developer/bunnyshield-m2m"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: cli <challenge>")
		os.Exit(1)
	}
	challengeText := os.Args[1]

	challenge, err := bunnyshieldm2m.ParseChallenge(challengeText)
	if err != nil {
		println("Failed to parse challenge: ", err)
		os.Exit(1)
	}

	answer, err := challenge.Solve()
	if err != nil {
		println("Failed to solve challenge: ", err)
		os.Exit(1)
	}

	println(answer.String())
}
