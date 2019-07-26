package main

import (
	"fmt"
	"log"
	"os"

	"github.com/axamon/tokenizzatore/vault"
)

func main() {

	if vault.IsOpen() == true {
		fmt.Println("Il vault è aperto")
		os.Exit(0)
	}

	err := vault.Apri()
	if err != nil {
		log.Println(err.Error())
	}

}
