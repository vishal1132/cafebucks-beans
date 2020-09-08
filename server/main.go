package main

import (
	"log"

	"github.com/vishal1132/cafebucks/config"
)

var beanSlice []beans

var coffeebeans map[string]beans

func main() {
	seedBeans()
	seedCoffeBeans()
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Println(err)
	}
	l := config.DefaultLogger(cfg)
	if err := runserver(cfg, l); err != nil {
		l.Fatal().Err(err).Msg("Failed to run order service server")
	}
}

func seedCoffeBeans() {
	coffeebeans = make(map[string]beans)
	coffeebeans["cappuccino"] = beans{"cappu", 300, 30}
	coffeebeans["frappuccino"] = beans{"frappu", 200, 25}
	coffeebeans["americano"] = beans{"ameri", 150, 20}
	coffeebeans["indiano"] = beans{"ameri", 150, 20}
	coffeebeans["espresso"] = beans{"frappu", 200, 25}
}

func seedBeans() {
	beanSlice = make([]beans, 0, 10)
	beanSlice = []beans{
		{"cappu", 300, 30},
		{"frappu", 200, 25},
		{"ameri", 150, 20},
	}
}
