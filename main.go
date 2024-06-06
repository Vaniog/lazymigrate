package main

import (
	"github.com/Vaniog/lazymigrate/app"
	"log"
)

func main() {
	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
