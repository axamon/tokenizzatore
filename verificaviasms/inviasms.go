package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/axamon/easyapiclient"
	"github.com/tkanos/gonfig"
)

// Configuration tiene gli elementi di configurazione
type Configuration struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var conf Configuration
var file = flag.String("file", "conf.json", "File di configurazione")

func inviasms(ctx context.Context, cell string, message string) error {

	// Parsa i parametri non di default passati all'avvio.
	flag.Parse()

	// Recupera valori dal file di configurazione passato come argomento.
	err := gonfig.GetConf(*file, &conf)

	if err != nil {
		log.Printf("Errore Impossibile recuperare informazioni dal file di configurazione: %s", *file)
	}

	// Recupera un token sms valido.
	token, _, err := easyapiclient.RecuperaToken(ctx, conf.Username, conf.Password)

	if err != nil {
		log.Printf("Errore nel recupero del token sms: %s\n", err.Error())
	}

	// fmt.Printf("token %s in scadenza tra %d secondi\n", token, scadenza)

	// Recupera lo shortnumber da usare per inviare sms.
	shortnumber, err := numerobreve(ctx, token)

	if err != nil {
		log.Printf("Errore, impossibile recuperare shortnumber %s\n", err.Error())
	}

	// Invia sms.
	err = easyapiclient.InviaSms(ctx, token, shortnumber, cell, message)

	if err != nil {
		log.Printf("Errore, sms non inviato: %s\n", err)
	}

	return err
}

func numerobreve(ctx context.Context, token string) (shortnumber string, err error) {

	type ShortNum struct {
		Number string `xml:"shortNumber"`
	}

	sNum := new(ShortNum)

	urlinfo := "https://easyapi.telecomitalia.it:8248/sms/v1/info"
	bearertoken := "Bearer " + token

	// Accetta anche certificati https non validi.
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Crea il cliet http.
	client := &http.Client{Transport: tr}

	// Crea la request da inviare.
	req, err := http.NewRequest("GET", urlinfo, nil)
	if err != nil {
		log.Printf("Errore creazione request: %v\n",
			req)
	}

	// Aggiunge alla request il contesto.
	req.WithContext(ctx)

	// Aggiunge alla request l'autenticazione.
	req.Header.Set("Authorization", bearertoken)

	// Aggiunge alla request gli header per passare le informazioni.
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Invia la request HTTP.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Errore %s\n", err.Error())
	}

	// Se la http response ha un codice di errore esce.
	if resp.StatusCode > 299 {
		fmt.Printf("Errore %d\n", resp.StatusCode)
		return
	}

	// Legge il body della risposta.
	bodyresp, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf(
			"Error Impossibile leggere risposta client http: %s\n",
			err.Error())
	}

	// Come da specifiche va chiuso il body.
	defer resp.Body.Close()

	err = xml.Unmarshal(bodyresp, &sNum)

	if err != nil {
		log.Printf(
			"Error Impossibile effettuare caricamento shortnumber: %s\n",
			err.Error())
	}

	// fmt.Println(sNum.Number)

	return sNum.Number, err

}
