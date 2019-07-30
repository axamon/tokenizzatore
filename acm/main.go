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
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

// Configuration contiene gli elemnti per configurare il tool.
type Configuration struct {
	Token string `json:"token"`
}

var configuration Configuration

var scadenza int
var username, password string

var configfile = flag.String("c", "conf.json", "File di configurazione")

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	err := Verificatoken()
	if err != nil {
		log.Printf("ERROR Token non verificato: %s\n", err.Error())
		os.Exit(1)
	}

	// var h []string
	//h, err := host.InfoWithContext(ctx)

	//fmt.Println(h.HostID)

	dati, err := diagnostictoolClient(ctx, os.Args[1])
	if err != nil {
		log.Printf("ERROR Impossiibile recuparare dati: %s", err.Error())
	}

	fmt.Println(dati)
}
