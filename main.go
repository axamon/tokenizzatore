package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/axamon/tokenizzatore/certificates"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	certificates.GenerateKeys(ctx, 2048, "private_key", "public_key")
	certificates.GenerateKeys(ctx, 2048, "GinoPrv", "GinoPub")
	certificates.GenerateKeys(ctx, 2048, "AndreaPrv", "AndreaPub")

	privateKeyImported, err := certificates.ImportPrivateKey(ctx, "private_key.pem")
	if err != nil {
		log.Println(err.Error())
	}

	publicKeyImported, err := certificates.ImportPublicKey(ctx, "public_key.pem")
	if err != nil {
		log.Println(err.Error())
	}

	ginoPubKey, err := certificates.ImportPublicKey(ctx, "GinoPub.pem")
	if err != nil {
		log.Println(err.Error())
	}

	// ciphertext, err := certificates.Encrypt(ctx, ginoPubKey, os.Args[1])
	// if err != nil {
	// 	log.Println(err.Error())
	// }

	//	fmt.Printf("%x\n", ciphertext)

	ciphertext, signature, err := certificates.EncryptAndSign(ctx, privateKeyImported, ginoPubKey, os.Args[1])
	if err != nil {
		log.Println(err.Error())
	}

	ginoPrivKey, err := certificates.ImportPrivateKey(ctx, "GinoPrv.pem")
	if err != nil {
		log.Println(err.Error())
	}

	messaggio := certificates.DecryptWithPrivateKey(ciphertext, ginoPrivKey)

	isSigned := certificates.VerifySignature(ctx, []byte(messaggio), signature, publicKeyImported)
	fmt.Println(string(messaggio), isSigned)

}
