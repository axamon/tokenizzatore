package vault

import (
	"os"
)

// IsOpen verifica se il vault è open o no.
func IsOpen() bool {

	// Verificare se vault è open
	isopen := os.Getenv("VAULTISOPEN")

	switch {
	case isopen == "open":
		return true
	default:
		return false

	}

}