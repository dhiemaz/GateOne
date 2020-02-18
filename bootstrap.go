package gate

import (
	"github.com/ory/ladon"
)

var (
	warden *ladon.Ladon
)

func init() {}

// Warden is ...
func Warden() *ladon.Ladon {
	return warden
}
