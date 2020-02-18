package gate

import "github.com/ory/ladon"

// IsAllow is ...
func IsAllow(r ladon.Request) error {
	return Warden().IsAllowed(&r)
}
