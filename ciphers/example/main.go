package main

import (
	"context"
	"fmt"
	"os"

	"github.com/axamon/tokenizzatore/ciphers"

	host "github.com/shirou/gopsutil/host"
)

func main() {

	str := os.Args[1]

	h, err := host.InfoWithContext(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}

	id := h.HostID

	fmt.Println(id)

	// prvT, pubT := ciphers.GenerateKeyPair(4048)
	// fmt.Println("CHIAVE T CREATA")
	prvC, pubC := ciphers.GenerateKeyPair(4048)
	fmt.Println("Chiavi client create")

	ciphers.SaveGobKey("private.key", prvC)
	ciphers.SaveGobKey("public.key", pubC)

	ciphers.SavePEMKey("private.pem", prvC)
	ciphers.SavePublicPEMKey("public.pem", *pubC)

	// // // Tokenizzatore cifra con chiave pubblica Tokenizatore
	// // cifratopubT := ciphers.EncryptWithPublicKey([]byte(str), pubT)

	// // Tokenizzatore cifra con chiave pubblica client
	cifratopubC := ciphers.EncryptWithPublicKey(cifratopubT, pubC)

	// // Client
	// incomprensibile := ciphers.DecryptWithPrivateKey(cifratopubTpubC, prvC)

	// cifratopubTpubC := ciphers.EncryptWithPublicKey(cifratopubT, pubC)

	// cifratopubC := ciphers.DecryptWithPrivateKey(cifratopubTpubC, prvT)

	// inchiaro := ciphers.DecryptWithPrivateKey(cifratopubC, prvC)

	// fmt.Println(string(inchiaro))
}
