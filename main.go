package main

import (
	//"fmt"
	"log"
	//"os"
	"todo/internal/cli"
	"todo/internal/store"
)

const FILE_PATH = "todos.json"

func main() {
	s := store.NewStore(FILE_PATH)
	if err := s.Load(); err != nil {
		log.Fatalln(err)
	}

	if err := cli.Handle(s); err != nil {
		log.Fatalln(err)
	}
}
