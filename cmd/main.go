package main

import (
	"log"

	"github.com/renatospaka/library/configs"
)

func main() {
	log.Println("iniciando a aplicação...")
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}

	log.Println("configuração carregada:", config.DBDriver)
}
