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

package main

import (
	"strconv"
	"strings"
	"time"
	"flag"
	"fmt"
	"crypto/cipher"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"io"
	"encoding/hex"
)

var giorni = flag.Int("d",7,"Giorni validità token: default 7")
var user = flag.String("u","","Username")
var pass = flag.String("p","","Password")


func main() {

	flag.Parse()

	oggi := time.Now()
	scadenzaUmana := oggi.Add(time.Duration(*giorni)*time.Hour*24)
	
	scadenza := scadenzaUmana.Unix()

	scadenzaStr := strconv.Itoa(int(scadenza))

	fmt.Println("Scadenza token: ", scadenzaUmana.Format("2006-01-02T15:04 UTC"))
	var elementi []string
    elementi = append(elementi, scadenzaStr, *user, *pass)
	str := strings.Join(elementi," ")

	c := encrypt([]byte(str), "vvkidtbcjujhtglivdjtlkgtetbtdejlivgukincfhdt")

	h := hex.EncodeToString(c)
	fmt.Println(h)

}


func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
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