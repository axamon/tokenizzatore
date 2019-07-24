package main

import (
	"crypto/sha256"
	"github.com/corvus-ch/shamir"
	"log"
	"os"
	"fmt"
	"encoding/base64"
	"encoding/json"
)

type Key struct {
	K byte `json:"k"`
	V []byte `jon:"v"`
}


func main() {

	m := make(map[byte][]byte)


	for _, chiave := range os.Args[1:] {

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

	h := sha256.New()
	h.Write(blob)
	hashmasterkey := h.Sum(nil)

	 
	fmt.Printf("%x\n", hashmasterkey)


	
}