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
        "crypto/tls"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "regexp"

)

const (
        endpointTgu   = "https://10.38.34.138:8443/DiagnosticTool/api.php?method=DiagnosticTool&sincrono=N&format=json&tgu="
        endpointEsito = "https://10.38.34.138:8443/DiagnosticTool/api.php?method=DiagnosticTool&sincrono=Y&format=json&cod_esito="
)

type response struct {
        Esito         string `json:"esito"`
        TDResponseCod string `json:"responsecode"`
        TDResponse    string `json:"response"`
        CodEsito      string `json:"cod_esito"`
}

var isAllDigits = regexp.MustCompile(`(?m)^\d+$`)


func diagnostictoolClient(ctx context.Context, tgu string) (result string, err error) {

        // fmt.Println(tgu)

        if !isAllDigits.MatchString(tgu) {
                log.Printf("ERROR formato tgu errato: %s\n", tgu)
                err := fmt.Errorf("ERROR formato tgu errato: %s", tgu)
                return "" , err
        }

        // Costringe il client ad accettare anche certificati https non validi
        // o scaduti.
        transCfg := &http.Transport{
                // Ignora certificati SSL scaduti.
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        client := &http.Client{Transport: transCfg}
        url := endpointTgu + tgu

        // fmt.Println(url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Printf("ERROR Impossibile creare richiesta: %s\n", err.Error())
        }

        // username := os.Getenv("DiagnosticToolUsername")
        // password := os.Getenv("DiagnosticToolPassoword")

        req.SetBasicAuth(username, password)
        req.WithContext(ctx)

        resp, err := client.Do(req)
        if err != nil {
                log.Printf("ERROR Impossibile inviare richiesta http: %s\n", err.Error())
        }
        defer resp.Body.Close()

        responsBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Printf("ERROR Impossibile leggere body reqest: %s\n", err.Error())
        }

        // fmt.Println(string(responsBody))

        risposta := new(response)

        err = json.Unmarshal(responsBody, &risposta)
        if err != nil {
                log.Println(err)
        }

        // fmt.Println(risposta)

        if risposta.Esito != "OK" {
                log.Printf("ERROR risposta interrogazione tgu %s non anadata a buon file: %s\n", tgu, risposta.Esito)
                return
        }

        dati, err := dt(ctx, risposta.CodEsito)
        if err != nil {
                log.Println(err)
        }

        return dati, err

}

func dt(ctx context.Context, code string) (str string, err error) {

        // Costringe il client ad accettare anche certificati https non validi
        // o scaduti.
        transCfg := &http.Transport{
                // Ignora certificati SSL scaduti.
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        client := &http.Client{Transport: transCfg}
        url := endpointEsito + code

        // fmt.Println(url)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Printf("ERROR Impossibile creare richiesta: %s\n", err.Error())
        }

        // username := os.Getenv("DiagnosticToolUsername")
        // password := os.Getenv("DiagnosticToolPassoword")

        req.SetBasicAuth(username, password)
        req.WithContext(ctx)


        resp, err := client.Do(req)
        if err != nil {
                log.Printf("ERROR Impossibile inviare richiesta http: %s\n", err.Error())
        }
        defer resp.Body.Close()

        responsBody, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                log.Printf("ERROR Impossibile leggere body reqest: %s\n", err.Error())
        }

        return string(responsBody), err
}