package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var server = flag.String("server", "127.0.0.1:9999/superadmin", "Connessione al server host:port")

var superAdminKey = flag.String("key", "", "SuperAdmin Key")

type SuperAdminKeys struct {
	SuperAdminKey string `json:"superadminkey"`
}

func main() {

	flag.Parse()

	var sak SuperAdminKeys

	sak.SuperAdminKey = *superAdminKey

	bytesRepresentation, err := json.Marshal(sak)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(bytesRepresentation))

	resp, err := http.Post("http://127.0.0.1:9999/superadmin", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp)

	fmt.Println(resp.Status, resp.StatusCode)

}
