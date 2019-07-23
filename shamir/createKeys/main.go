package main

import (
	"encoding/json"
	"encoding/base64"
	"os"
	"flag"
	"fmt"
	"log"
	"github.com/corvus-ch/shamir"
)

var secret = flag.String("s","","secret")
var numKeys = flag.Int("nk",3,"Number of keys to create")
var numMin = flag.Int("nm",2,"Minimum number of kesys to unseal")

type Key struct {
	K byte `json:"k"`
	V []byte `jon:"v"`
}

func main() {

	flag.Parse()
	
	if *numMin > *numKeys {
		log.Printf("ERORR the number of minimal keys to use cannot be greater than the number of keys to create: %d > %d", *numMin, *numKeys)
		os.Exit(1)
	}

	m, err := shamir.Split([]byte(*secret),*numKeys,*numMin)
	if err != nil {
		log.Println(err.Error())
	}

	// fmt.Println(m)

	n := 0
	fmt.Printf("Creo %d chiavi, minimo %d serviranno per sbloccare la password\n", *numKeys, *numMin)
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
	}
		




