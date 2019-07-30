package vault

import (
	"github.com/howeyc/gopass"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/axamon/tokenizzatore/vault/creatoken"

	"github.com/corvus-ch/shamir"
)

// Apri apre il Vault verificando che le chiavi passate sblocchino la masterkey.
func Apri() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// m mappa per archiviare le chiavi superAdmin
	m := make(map[byte][]byte)

	// numero minimo di chiavi SuperAdmin per aprire il vault.
	threshold, err := recuperaThreshold()
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Printf("Per Aprire il Vault dovrai inserire %d chiavi SuperAdmin\n", threshold)
	// Richiede un numero di chiavi SuperAdmin pari a threshold.
	for i := 0; i < threshold; i++ {
		// Richiede inserimento di una chiave superAdmin
		fmt.Printf("Inserisci la chiave SuperAdmin numero %d: > ", i+1)
		chiave, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println(err)
		}

		// Elabora la chiave passata da base64 a slice di bytes.
		var decoded []byte
		decoded, err = base64.StdEncoding.DecodeString(string(chiave))
		if err != nil {
			log.Println("decode error:", err)
		}

		// Inserisce in key la le info della chiave.
		var key Key

		err = json.Unmarshal(decoded, &key)
		if err != nil {
			log.Println(err.Error())
		}

		c := byte(key.K)
		cs := []byte(key.V)
		m[c] = cs

	}

	blob, err := shamir.Combine(m)
	if err != nil {
		log.Println(err.Error())
	}

	// fmt.Println(string(blob))

	mastersecret := string(blob)

	aprivault(ctx, mastersecret)

	if IsOpen() == true {

		
		fmt.Println("Il vault del Tokenizzatore è aperto")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			tokengenerated, err := creatoken.OneWeekValidity()
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Fprintf(w, "Benvenuto nel tokenizzatore, il tuo nuovo token:%s\n", tokengenerated)
		})
		log.Println(http.ListenAndServe(":9999", nil))
	}

	return err
}

func recuperaThreshold() (threshold int, err error) {

	// Apre il file di configurazione del Vault VaultConf
	jsonFile, err := os.Open(VaultConf)
	if err != nil {
		log.Println(err.Error())
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err.Error())
	}

	var conf VaultConfStr

	// Parsa quanto è nel file di configurazione sulla variabile stuct.
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		log.Println(err.Error())
	}

	// Recupera il valore di threshold.
	threshold = conf.Threshold

	return threshold, err
}

func aprivault(ctx context.Context, mastersecret string) error {

	hashmasterkeyfromfile, err := ioutil.ReadFile(Vaulthash)
	if err != nil {
		log.Println(err.Error())
	}

	// TODO aprire veramente il vault
	h := sha256.New()
	h.Write([]byte(mastersecret))
	hashmasterkey := h.Sum(nil)

	// fmt.Printf("%x\n", hashmasterkey)

	if string(hashmasterkeyfromfile) == string(hashmasterkey) {

		fmt.Println("ok")

		err := os.Setenv("VAULTISOPEN", "open")
		if err != nil {
			log.Println(err.Error())
		}

	}

	return err
}
