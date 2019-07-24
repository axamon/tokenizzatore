package vault

import (
	"fmt"
	"context"
	"github.com/corvus-ch/shamir"
	"encoding/json"
	"log"
	"encoding/base64"
)

// Apri apre il Vault verificando che le chiavi passate insieme corrispondano alla masterkey.
func Apri(chiavi []string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// m mappa per archiviare le chiavi superAdmin
	m := make(map[byte][]byte)


	for _, chiave := range chiavi {

		var decoded []byte
		decoded, err := base64.StdEncoding.DecodeString(chiave)
		if err != nil {
			log.Println("decode error:", err)
		}

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

	masterkey := string(blob)

	aprivault(ctx, masterkey)
	 
	
return
}

func aprivault(ctx context.Context, masterkey string) {
	fmt.Println(masterkey)
}