package vault

import (
	"io/ioutil"
	"crypto/sha256"
	"context"
	"time"
	"strings"
	"math/rand"
	"encoding/json"
	"encoding/base64"
	"os"
	"fmt"
	"log"
	"github.com/corvus-ch/shamir"
)


// CreaChiaviSuperAdmin crea le chiavi per avviare o aprire il vault del tokenizzatore.
func CreaChiaviSuperAdmin(numKeys, numMin int) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	if numMin > numKeys {
		log.Printf("ERORR the number of minimal keys to use cannot be greater than the number of keys to create: %d > %d", numMin, numKeys)
		os.Exit(1)
	}

	mastersecret := createPassword(ctx, 100)

	fmt.Println(mastersecret)

	m, err := shamir.Split([]byte(mastersecret),numKeys,numMin)
	if err != nil {
		log.Println(err.Error())
	}

	
	h := sha256.New()
	h.Write([]byte(mastersecret))
	hashmasterkey := h.Sum(nil)

	// fmt.Println(m)

	n := 0
	fmt.Printf("Creo %d chiavi, minimo %d serviranno per sbloccare la password\n", numKeys, numMin)
	for c,cs := range m {
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
	fmt.Println("hash masterkey:")
	fmt.Printf("%x\n",string(hashmasterkey))

	// Salvo nel file l'hash della masterkey.
	err = ioutil.WriteFile(Vaulthash, hashmasterkey, 0600)
    if err != nil {
		log.Println(err.Error())
	}
		

	return
}




	func createPassword(ctx context.Context, length int) string {
		rand.Seed(time.Now().UnixNano())
		chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
		
		var b strings.Builder
		for i := 0; i < length; i++ {
			   b.WriteRune(chars[rand.Intn(len(chars))])
		}
		str := b.String() // E.g. "ExcbsVQs"

		

	return str
	}


	func salvamastersecret(ctx context.Context, mastersecret string) {

		fmt.Println(mastersecret)
		return
	}