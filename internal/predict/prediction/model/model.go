package model

import "fmt"

// the Input for the prediction
type Input struct {
	ID string `json:"id"`
}

// the Result of the prediction
type Result struct {
	Class       string  `json:"class"`
	Probability float32 `json:"probability"`
}

func (r *Result) String() string {
	isCat := r.Class == "cats"
	var suffix string
	if isCat {
		suffix = "this is a cat"
	} else {
		suffix = "this is not a cat"
	}

	return fmt.Sprintf("I'm pretty sure %s", suffix)
}

// the ErrorResult of the prediction
type ErrorResult struct {
	Message string `json:"message"`
}
