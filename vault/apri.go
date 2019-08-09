package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"

	"github.com/howeyc/gopass"

	"github.com/axamon/tokenizzatore/vault/creatoken"

	"github.com/corvus-ch/shamir"
	"golang.org/x/crypto/bcrypt"
)

var isOpen = false

// Apri apre il Vault verificando che le chiavi passate sblocchino la masterkey.
func Apri(vaulthash string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// m mappa per archiviare le chiavi superAdmin
	m := make(map[byte][]byte)

	// Ripulisca la mappa ogni tot tempo.
	cleanup := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-cleanup.C:
				if len(m) > 0 {
					for k := range m {
						delete(m, k)
					}
					log.Println("mappa ripulita")
				}
			}
		}
	}()

	// numero minimo di chiavi SuperAdmin per aprire il vault.
	threshold, err := recuperaThreshold()
	if err != nil {
		log.Println(err.Error())
	}

	// Per sbloccare il vault da console
	go func() {

		fmt.Printf("Per Aprire il Vault dovrai inserire %d chiavi SuperAdmin\n", threshold)
		// Richiede un numero di chiavi SuperAdmin pari a threshold.
		//for i := 0; i < threshold; i++ {

		var i int
		for {

			// Richiede inserimento di una chiave superAdmin
			fmt.Printf("Inserisci la chiave SuperAdmin numero %d: > ", i+1)
			chiave, err := gopass.GetPasswdMasked()
			if err != nil {
				if err.Error() == "interrupted" {
					os.Exit(1)
				}
				fmt.Println(err)
			}
			i++

			// Elabora la chiave passata da base64 a slice di bytes.
			var decoded []byte
			decoded, err = base64.StdEncoding.DecodeString(string(chiave))
			if err != nil {
				log.Println("decode error:", err)
			}

			// Inserisce in key le info della chiave.
			var key Key

			err = json.Unmarshal(decoded, &key)
			if err != nil {
				log.Println(err.Error())
			}

			c := byte(key.K)
			cs := []byte(key.V)
			m[c] = cs

			// Almeno 2 chiavi devono essere passate.
			if i < 2 {
				continue
			}
			mastersecret, err := shamir.Combine(m)
			if err != nil {
				log.Println(err.Error())
			}

			aprivault(ctx, vaulthash, mastersecret)

			//	fmt.Println(isOpen) // debug

			// Se il vault si apre esce dal ciclo for.
			if isOpen == true {
				return
			}

		}
	}()

	// Per sbloccare via http
	r := chi.NewRouter()

	r.Post("/superadmin", func(w http.ResponseWriter, r *http.Request) {

		type SuperAdminKeys struct {
			SuperAdminKey string `json:"superadminkey"`
		}

		fmt.Println("e' arrivato un post")
		decoder := json.NewDecoder(r.Body)

		var decodedd SuperAdminKeys
		err := decoder.Decode(&decodedd)
		if err != nil {
			log.Println("errore decodifica", err.Error())
		}

		fmt.Println(decodedd.SuperAdminKey)

		// Elabora la chiave passata da base64 a slice di bytes.
		decoded, err := base64.StdEncoding.DecodeString(decodedd.SuperAdminKey)
		if err != nil {
			log.Println("decode error:", err)
		}

		// Inserisce in key le info della chiave.
		var key Key

		err = json.Unmarshal(decoded, &key)
		if err != nil {
			log.Println(err.Error())
		}

		c := byte(key.K)
		cs := []byte(key.V)
		m[c] = cs

		mastersecret, err := shamir.Combine(m)
		if err != nil {
			log.Println(err.Error())
		}

		aprivault(ctx, vaulthash, mastersecret)

		switch isOpen {
		case true:
			fmt.Fprintf(w, "Il tokenizzatore è attivo")
		case false:
			fmt.Fprintf(w, "Il tokenizzatore è disattivato")
		}

	})

	// Per vedere lo stato del tokenizzatore da web
	r.Get("/stato", func(w http.ResponseWriter, r *http.Request) {
		switch isOpen {
		case true:
			fmt.Fprintf(w, "aperto")
		case false:
			fmt.Fprintf(w, "chiuso")
		}
	})

	// Per richiedere un nuovo token da web
	r.Get("/token", func(w http.ResponseWriter, r *http.Request) {
		if isOpen == true {
			tokengenerated, err := creatoken.FiveMinutes()
			if err != nil {
				log.Println(err.Error())
			}
			fmt.Fprintf(w, tokengenerated)
		}
		if isOpen == false {
			fmt.Fprintf(w, "Il tokenizzatore è chiuso. Contatta i SuperAdmin per riaprirlo.")
		}
	})

	http.ListenAndServe(":9999", r)

	fmt.Scanln()
	fmt.Println("done")

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

func aprivault(ctx context.Context, vaulthash string, mastersecret []byte) error {

	// Recupera hash della mastersecret dal file dove è stata salvata.
	hashmasterkeyfromfile, err := ioutil.ReadFile(vaulthash)
	if err != nil {
		log.Println(err.Error())
	}

	err = bcrypt.CompareHashAndPassword(hashmasterkeyfromfile, mastersecret)

	// Verifica se i due hash corrispondono.
	if err == nil {
		isOpen = true

		// fmt.Println("ok")

		err := os.Setenv("VAULTISOPEN", "open")
		if err != nil {
			log.Println(err.Error())
		}

	}

	return err
}
