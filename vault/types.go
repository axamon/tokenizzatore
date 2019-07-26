package vault

// Vaulthash dove è salvato l'hash della masterkey.
var Vaulthash = "vaulthash"

// VaultDB dove è salvato l'hash della masterkey.
var VaultDB = "vault.db"

// VaultConf dove è salvato l'hash della masterkey.
var VaultConf = "vault.json"

// VaultConfStr è la struttura del file di configurazione.
type VaultConfStr struct {
	Email     string
	Version   string
	Threshold int
}

var Version = "ver.1b"

// Key contiene i bytes delle chiavi per superAdmin
type Key struct {
	K byte   `json:"k"`
	V []byte `jon:"v"`
}
