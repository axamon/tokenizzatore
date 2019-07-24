package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/axamon/tokenizzatore/vault"
)

var numkey = flag.Int("nk",3,"Numero di chiavi superAdmin da creare")
var threashhold = flag.Int("nm",2,"Numero minimo di chiavi per sbloccare")

func main() {

	flag.Parse()

	if vault.IsOpen() == true {
		fmt.Println("Il vault è aperto")
		os.Exit(0)
	}

	vault.CreaChiaviSuperAdmin(*numkey,*threashhold)

	vault.Apri(*threashhold)

	
}