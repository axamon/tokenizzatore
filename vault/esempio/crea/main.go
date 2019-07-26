package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/axamon/tokenizzatore/vault"
)

var numkey = flag.Int("n", 3, "Numero di chiavi superAdmin da creare")
var threashhold = flag.Int("t", 2, "Numero minimo di chiavi per sbloccare")

func main() {

	flag.Parse()

	if vault.IsOpen() == true {
		fmt.Println("Il vault è aperto")
		os.Exit(0)
	}

	vault.Crea(*numkey, *threashhold)

}
