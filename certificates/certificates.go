package certificates

import (
	"bufio"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func GenerateKeys(ctx context.Context, bits int, privatefilename, publicfileame string) {

	// Generate RSA private and public Keys
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println(err.Error())
	}

	publicKey := &privateKey.PublicKey

	// fmt.Println("Private Key : ", privateKey)
	// fmt.Println("Public key ", publicKey)

	exportPrivateKey(ctx, privateKey, privatefilename)

	exportPublicKey(ctx, publicKey, publicfileame)

	return
}

func exportPrivateKey(ctx context.Context, privatekey *rsa.PrivateKey, filename string) (err error) {
	// Export Private Key
	pemPrivateFile, err := os.Create(filename + ".pem")
	defer pemPrivateFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	var pemPrivateBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
	}

	err = pem.Encode(pemPrivateFile, pemPrivateBlock)

	if err != nil {
		fmt.Println(err)
	}

	return err
}

func exportPublicKey(ctx context.Context, publickey *rsa.PublicKey, filename string) (err error) {
	// Export Private Key
	pemPublicFile, err := os.Create(filename + ".pem")
	defer pemPublicFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	var pemPublicBlock = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publickey),
	}

	err = pem.Encode(pemPublicFile, pemPublicBlock)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func ImportPrivateKey(ctx context.Context, privatekeyfile string) (privateKeyImported *rsa.PrivateKey, err error) {
	//Import Private Key
	privateKeyFile, err := os.Open(privatekeyfile)
	defer privateKeyFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	pemfileinfo, _ := privateKeyFile.Stat()
	var size = int64(pemfileinfo.Size())
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	privateKeyImported, err = x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Private Key : ", privateKeyImported)

	return privateKeyImported, err
}

func ImportPublicKey(ctx context.Context, publickeyfile string) (publicKeyImported *rsa.PublicKey, err error) {
	//Import Private Key
	publicKeyFile, err := os.Open(publickeyfile)
	defer publicKeyFile.Close()

	if err != nil {
		fmt.Println(err)
	}

	pemfileinfo, _ := publicKeyFile.Stat()
	var size = int64(pemfileinfo.Size())
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(pembytes)

	data, _ := pem.Decode([]byte(pembytes))

	publicKeyImported, err = x509.ParsePKCS1PublicKey(data.Bytes)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println("Private Key : ", privateKeyImported)

	return publicKeyImported, err
}

func Encrypt(ctx context.Context, publickey *rsa.PublicKey, msg string) (ciphertext []byte, err error) {
	message := []byte(msg)
	label := []byte("")
	hash := sha512.New()

	ciphertext, err = rsa.EncryptOAEP(hash, rand.Reader, publickey, message, label)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	fmt.Println()

	return ciphertext, err
}

func EncryptAndSign(ctx context.Context, privatekey *rsa.PrivateKey, publickey *rsa.PublicKey, msg string) (ciphertext []byte, signaure []byte, err error) {
	message := []byte(msg)
	label := []byte("")
	hash := sha512.New()

	ciphertext, err = rsa.EncryptOAEP(hash, rand.Reader, publickey, message, label)

	if err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("OAEP encrypted [%s] to \n[%x]\n", string(message), ciphertext)
	// fmt.Println()

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := message
	newhash := crypto.SHA512
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privatekey, newhash, hashed, &opts)

	if err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("PSS Signature : %x\n", signature)

	return ciphertext, signature, err
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Printf(err.Error())
	}
	return plaintext
}

func VerifySignature(ctx context.Context, data []byte, signature []byte, publicKey *rsa.PublicKey) bool {

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto // for simple example
	PSSmessage := data
	newhash := crypto.SHA512
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)
	err := rsa.VerifyPSS(publicKey, newhash, hashed, signature, nil)

	if err == nil {
		return true
	}

	return false
}
