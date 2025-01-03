package main

import (
	"context"
	"github.com/tlb_katia/auth/internal/app"
	"log"
)

func main() {
	a, err := app.NewApp(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err = a.Run(); err != nil {
		log.Fatal(err)
	}
}
