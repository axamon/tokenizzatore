package main

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func generaPin(ctx context.Context, length int) string {

	seed := time.Now().UnixNano()

	rand.Seed(seed)

	var pin []string

	for i := 0; i < length; i++ {
		c := strconv.Itoa(rand.Intn(9))
		pin = append(pin, c)
	}

	return strings.Join(pin, "")
}
