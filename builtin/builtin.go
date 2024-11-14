package builtin

import (
	"io"
	"math/rand"

	"github.com/dece2183/hexowl/builtin/types"
)

var descriptor types.Descriptor

func init() {
	descriptor.Constants = constants
	descriptor.Functions = functions
	descriptor.System.Stdout = io.Discard
}

// Provide your own system description.
//
// Use this function to implement native integration into your application.
func SystemInit(sys types.System) {
	descriptor.System = sys
	if sys.RandomSeed != 0 {
		rand.Seed(sys.RandomSeed)
	}
}
