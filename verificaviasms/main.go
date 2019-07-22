package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	v := make(chan []byte)
	var pin string

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Troppo tempo")
			os.Exit(1)
		case pinr := <-v:
			if pin == string(pinr) {
				fmt.Println("ok")
			}
		case <-time.After(30 * time.Second):
			fmt.Println("Tempo scaduto")
			os.Exit(1)
		}
	}()

	// Genera un pin di lunghezza 6
	pin = generaPin(ctx, 6)

	// fmt.Println(pin)

	// Richiede numero cell
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("Inserisci il tuo numero di cellulare: > ")
	cellulare, err := buf.ReadBytes('\n')
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(cellulare))
	}

	cell := string(cellulare)

	message := "il tuo pin è: " + pin

	err = inviasms(ctx, cell, message)
	if err != nil {
		log.Printf("ERROR Impossibile inviare sms: %s\n", err.Error())
	}

	// Richiede inserimento pin ricevuto su cell
	bufp := bufio.NewReader(os.Stdin)
	fmt.Print("Inserisci il pin ricevuto sul cellulare: > ")
	pinricevuto, err := bufp.ReadBytes('\n')
	v <- pinricevuto
	if err != nil {
		fmt.Println(err)
	}

}
