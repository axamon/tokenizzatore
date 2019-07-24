package vault

// File dove è salvato l'hash della masterkey.
var Vaulthash = "vaulthash"

// Key contiene i bytes delle chiavi per superAdmin
type Key struct {
	K byte `json:"k"`
	V []byte `jon:"v"`
}
