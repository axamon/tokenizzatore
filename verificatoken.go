
package main

import (
	"os"
	"strings"
	"github.com/tkanos/gonfig"
	"time"
	"strconv"
	"encoding/hex"
	"log"
	"crypto/cipher"
	"crypto/aes"
	"crypto/md5"
)
		
		
const (
		passphrase = "vvkidtbcjujhtglivdjtlkgtetbtdejlivgukincfhdt"
	)
		
		
		
// Verificiatoken verifica che il token sia valido.		
func Verificatoken() (err error) {
		
	err = gonfig.GetConf("conf.json", &configuration)
	if err != nil {
		log.Printf("ERROR Problema con il file di configurazione conf.json: %s\n", err.Error())
		return
	}
		
	credBlob, _ := hex.DecodeString(configuration.Token)
	userEpass := string(decrypt(credBlob))
	credenziali := strings.Split(userEpass, " ")
		
	scadenza, err := strconv.Atoi(credenziali[0])
	if err != nil {
		log.Printf("ERROR Impossibile parsare scadenza del token: %s\n", err.Error())
	}
				// username = credenziali[1]
				// password = credenziali[2]
		
	oggi := time.Now().Unix()
		
	if oggi > int64(scadenza) {
		log.Println("Token scaduto. Impossibile proseguire.")
			os.Exit(1)
		}
		
		return err
	}
		
		
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
		
		
func decrypt(data []byte) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println(err.Error())
	}
	return plaintext
}