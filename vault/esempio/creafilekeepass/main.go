package main

import (
	"context"
	"os"

	"github.com/axamon/tokenizzatore/vault/keepass"
)

func main() {

	ctx := context.TODO()

	keepass.CreaFile(ctx, os.Args[1], os.Args[2])
}
