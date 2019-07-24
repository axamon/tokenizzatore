package vault

import (
	"net/http"
	"io/ioutil"
	"crypto/sha256"
	"os"
	"bufio"
	"fmt"
	"context"
	"github.com/corvus-ch/shamir"
	"encoding/json"
	"log"
	"encoding/base64"
)

// Apri apre il Vault verificando che le chiavi passate insieme corrispondano alla masterkey.
func Apri(threshold int) {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// m mappa per archiviare le chiavi superAdmin
	m := make(map[byte][]byte)


	// Richiede un numero di chiavi superAdmin pari a threshold. 
	for i:=0 ; i<threshold ; i++ {
	// Richiede inserimento di una chiave superAdmin
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("Inserisci la tua chiave superAdmin: > ")
	chiave, err := buf.ReadBytes('\n')   // TODO: Evitare che ci sia echo
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
		http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Benvenuto nel tokenizzatore, Lavori in corso\n")
		})
		log.Println(http.ListenAndServe(":9999", nil))
	}
	
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

		err := os.Setenv("VAULTISOPEN","open")
		if err != nil {
			log.Println(err.Error())
		}

	}
	
	
	return err
}