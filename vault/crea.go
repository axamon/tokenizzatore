package vault

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/corvus-ch/shamir"
)

// Crea crea la mastersecret, il vault.db e le chiavi SuperAdmin per avviare o aprire il vault del tokenizzatore.
func Crea(numKeys, threshold int) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if threshold > numKeys {
		log.Printf("ERORR the number of minimal keys to use cannot be greater than the number of keys to create: %d > %d", threshold, numKeys)
		os.Exit(1)
	}

	// Crea la chiave mastersecret con cui verrà criptato il Vault.
	mastersecret := createMastersecret(ctx, 100)

	// fmt.Println(mastersecret) //debug

	// Crea un numero di chiavi di tipo SuperAdmin pari a numKeys, saranno necessarie un
	// numero di chiavi di tipo SuperAdmin pari a threshold per aprire il vault.
	m, err := shamir.Split([]byte(mastersecret), numKeys, threshold)
	if err != nil {
		log.Println(err.Error())
	}

	// Creo hash della mastersecret
	h := sha256.New()
	h.Write([]byte(mastersecret))
	hashmastersecret := h.Sum(nil)

	// Salvo nel file Vaulthash l'hash della mastersecret.
	err = ioutil.WriteFile(Vaulthash, hashmastersecret, 0600)
	if err != nil {
		log.Println(err.Error())
	}

	n := 0
	fmt.Printf("Creo %d chiavi, minimo %d serviranno per sbloccare la password\n", numKeys, threshold)
	for c, cs := range m {
		n++

		var keyn []byte
		key := Key{K: c, V: cs}
		keyn, err := json.Marshal(key)
		// fmt.Println(string(keyn))
		if err != nil {
			log.Println(err.Error())
		}
		keyb := base64.StdEncoding.EncodeToString(keyn)

		fmt.Println("Chiave ", n)
		fmt.Println(keyb)
		fmt.Println()

	}
	// fmt.Println("hash masterkey:")
	// fmt.Printf("%x\n", string(hashmastersecret))

	// Cifra dei dati con la mastersecret
	// data := simmetric.Encrypt([]byte("DB VAult"), mastersecret)

	data := VaultConfStr{
		Email:     "alberto.bregliano@telecomitalia.it",
		Version:   Version,
		Threshold: threshold}

	dataj, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
	}

	err = ioutil.WriteFile(VaultConf, dataj, 0600)
	if err != nil {
		log.Println(err.Error())
	}

	// Crea il file VauldDB
	err = ioutil.WriteFile(VaultDB, nil, 0600)
	if err != nil {
		log.Println(err.Error())
	}

	return
}

func createMastersecret(ctx context.Context, length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" + ".-_@;%£")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}
