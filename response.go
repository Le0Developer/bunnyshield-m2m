package bunnyshieldm2m

import "fmt"

type Answer struct {
	challenge Challenge
	answer    int
}

func (r Answer) String() string {
	return r.challenge.string() + "#" + fmt.Sprint(r.answer)
}
