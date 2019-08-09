package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/corvus-ch/shamir"
	"golang.org/x/crypto/bcrypt"
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

	// Genera hash della masterkey con bcrypt
	hashmastersecret, err := bcrypt.GenerateFromPassword([]byte(mastersecret), bcrypt.DefaultCost)

	// Salvo nel file Vaulthash l'hash della mastersecret.
	err = ioutil.WriteFile(Vaulthash, hashmastersecret, 0600)
	if err != nil {
		log.Println(err.Error())
	}

	// n conta il numero di chiavi Superadmin inserito
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
		keyb64 := base64.StdEncoding.EncodeToString(keyn)

		fmt.Println("Chiave ", n)
		fmt.Println(keyb64)
		fmt.Println()

	}

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
	// token := make([]byte, length)
	// rand.Read(token)

	// fmt.Println(token)

	str := b.String()

	return str
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func generatepass() {
	MIN := 0
	MAX := 94
	SEED := time.Now().Unix()
	var LENGTH int64 = 8

	arguments := os.Args
	switch len(arguments) {
	case 2:
		LENGTH, _ = strconv.ParseInt(os.Args[1], 10, 64)
		if LENGTH <= 0 {
			LENGTH = 8
		}
	default:
		fmt.Println("Using default values!")
	}

	rand.Seed(SEED)
	startChar := "!"
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		fmt.Print(newChar)
		if i == LENGTH {
			break
		}
		i++
	}
	fmt.Println()
}
