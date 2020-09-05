package main

import (
	"log"

	"github.com/vishal1132/cafebucks/config"
)

var beanSlice []beans

func main() {
	seedBeans()
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Println(err)
	}
	l := config.DefaultLogger(cfg)
	if err := runserver(cfg, l); err != nil {
		l.Fatal().Err(err).Msg("Failed to run order service server")
	}
}

func seedBeans() {
	beanSlice = make([]beans, 0, 10)
	beanSlice = []beans{
		{"cappu", 300, 30},
		{"frappu", 200, 25},
		{"ameri", 150, 20},
	}
}
