package main

import (
	"log"

	"github.com/vishal1132/cafebucks/config"
)

var beanSlice []*beans

var coffeebeans map[string]*beans

var (
	cappu  beans
	frappu beans
	ameri  beans
)

func main() {
	seedBeans()
	seedBeansSlice()
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
	coffeebeans = make(map[string]*beans)
	coffeebeans["cappuccino"] = &cappu
	coffeebeans["frappuccino"] = &frappu
	coffeebeans["americano"] = &ameri
	coffeebeans["indiano"] = &cappu
	coffeebeans["espresso"] = &frappu
}

func seedBeansSlice() {
	beanSlice = make([]*beans, 0, 10)
	beanSlice = []*beans{&cappu, &frappu, &ameri}
}

func seedBeans() {
	cappu = beans{"cappu", 300, 30}
	frappu = beans{"frappu", 200, 25}
	ameri = beans{"ameri", 150, 20}
}
