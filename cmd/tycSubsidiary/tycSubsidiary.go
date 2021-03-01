package main

import (
	"github.com/canc3s/cSubsidiary/internal/runner"
)

func main() {
	options := runner.ParseOptions()

	runner.RunEnumeration(options)
}

