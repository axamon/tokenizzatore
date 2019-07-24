package main

import (
	"flag"
	"github.com/axamon/tokenizzatore/vault"
)

var numkey = flag.Int("nk",3,"Numero di chiavi superAdmin da creare")
var threashhold = flag.Int("nm",2,"Numero minimo di chiavi per sbloccare")

func main() {

	flag.Parse()
	vault.CreaChiaviSuperAdmin(*numkey,*threashhold)
}