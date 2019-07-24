package vault

// Key contiene i bytes delle chiavi per superAdmin
type Key struct {
	K byte `json:"k"`
	V []byte `jon:"v"`
}
