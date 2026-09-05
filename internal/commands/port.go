package commands

import "fmt"

const (
	portMin = 1
	portMax = 65535
)

func validatePort(p int) error {
	if p < portMin || p > portMax {
		return fmt.Errorf("invalid port (%d-%d)", portMin, portMax)
	}
	return nil
}

// ValidatePort is the single source of truth for port range validation,
// shared by CLI pre-checks and filter construction.
func ValidatePort(p int) error {
	return validatePort(p)
}
