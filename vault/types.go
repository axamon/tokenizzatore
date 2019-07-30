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

// Version è la versione del software.
var Version = "ver.2b"

// Key contiene i bytes delle chiavi per superAdmin
type Key struct {
	K byte   `json:"k"`
	V []byte `jon:"v"`
}

// Configuration contiene gli elementi per configurare il tool.
type Configuration struct {
	Token string `json:"token"`
}

var configuration Configuration
