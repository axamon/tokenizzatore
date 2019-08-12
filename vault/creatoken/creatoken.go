// Copyright (c) 2019 Alberto Bregliano
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package creatoken

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	m "math/rand"
	"strconv"
	"strings"
	"time"
)

const simmetricpass = "vvkidtbcjujhtglivdjtlkgtetbtdejlivgukincfhdt"

var Dinamicsimmetricpass = generatepass(context.TODO(), 14)

var giorni = flag.Int("d", 7, "Giorni validità token in giorni")
var user = flag.String("user", "", "Username")
var pass = flag.String("pass", "", "Password")

// TokenWithCredentials crea un token con user e pass all'interno.
func TokenWithCredentials(user, pass string) (token string, err error) {
	oggi := time.Now()
	scadenzaUmana := oggi.Add(time.Duration(7) * time.Hour * 24)

	scadenza := scadenzaUmana.Unix()

	scadenzaStr := strconv.Itoa(int(scadenza))

	fmt.Println("Scadenza token: ", scadenzaUmana.Format("2006-01-02T15:04 UTC"))
	var elementi []string
	elementi = append(elementi, scadenzaStr, user, pass)
	str := strings.Join(elementi, " ")

	c := encrypt([]byte(str), simmetricpass)

	token = hex.EncodeToString(c)

	return token, err
}

// OneWeekValidity crea un token con la durata di sette giorni.
func OneWeekValidity() (token string, err error) {

	oggi := time.Now()
	scadenzaUmana := oggi.Add(time.Duration(7) * time.Hour * 24)

	scadenza := scadenzaUmana.Unix()

	scadenzaStr := strconv.Itoa(int(scadenza))

	fmt.Println("Scadenza token: ", scadenzaUmana.Format("2006-01-02T15:04 UTC"))
	var elementi []string
	elementi = append(elementi, scadenzaStr, "test", "test")
	str := strings.Join(elementi, " ")

	c := encrypt([]byte(str), simmetricpass)

	token = hex.EncodeToString(c)

	return token, err
}

// FiveMinutes crea un token con durata di 5 minuti.
func FiveMinutes() (token string, err error) {

	oggi := time.Now()
	scadenzaUmana := oggi.Add(time.Duration(7) * time.Minute * 5)

	scadenza := scadenzaUmana.Unix()

	scadenzaStr := strconv.Itoa(int(scadenza))

	fmt.Println("Scadenza token: ", scadenzaUmana.Format("2006-01-02T15:04 UTC"))
	var elementi []string
	elementi = append(elementi, scadenzaStr, "test", "test")
	str := strings.Join(elementi, " ")

	c := encrypt([]byte(str), simmetricpass)

	token = hex.EncodeToString(c)

	return token, err
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
//	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	block, _ := aes.NewCipher([]byte(createHash(Dinamicsimmetricpass)))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func generatepass(ctx context.Context, length int) string {
	m.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" + ".-_@;%£")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[m.Intn(len(chars))])
	}

	str := b.String()

	return str
}
