package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/axamon/tokenizzatore/vault"
)

var file = flag.String("f", "../crea/vaulthash", "File vaulthash da aprire")

func main() {

	flag.Parse()

	if vault.IsOpen() == true {
		fmt.Println("Il vault è aperto")
		os.Exit(0)
	}

	err := vault.Apri(*file)
	if err != nil {
		log.Println(err.Error())
	}

}
