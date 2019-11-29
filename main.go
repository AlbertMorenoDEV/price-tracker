package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
)

type data struct {
	Title string
	Price float64
}

var p = []string{
"https://www.vivobarefoot.com/es/mens/everyday/gobi-ii-mens?colour=Black",
"https://www.vivobarefoot.com/eu/mens/outdoor/primus-trail-fg-mens-ow5?colour=Charcoal",
"https://www.vivobarefoot.com/eu/mens/outdoor/tracker-fg-mens?colour=Dark+Brown",
"https://www.vivobarefoot.com/eu/mens/active/primus-lite-mens-ow5?colour=Obsidian",
"https://www.vivobarefoot.com/eu/mens/outdoor/primus-trail-sg-mens?colour=Black",
"https://www.vivobarefoot.com/eu/mens/everyday/ra-ii-mens?colour=Midnight+Navy",
"https://www.vivobarefoot.com/eu/mens/active/primus-knit-lux-mens?colour=Black",
"https://www.vivobarefoot.com/eu/mens/active/primus-knit-mens?colour=Mood+Indigo",
}

func main() {
	n := len(p)
	d := make(chan data, n)

	go getAllUrls(d)

	for r := range d {
		fmt.Println(r)
	}
}

func getAllUrls(d chan data) {
	for _, u := range p {
		fmt.Println("URL: " + u)
		d <- getData(u)
	}
	close(d)
}

func getData(u string) data {
	resp, err := soup.Get(u)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)

	title := doc.Find("h1").Text()
	dataFpAttr := doc.Find("div", "class", "price").Find("span").Attrs()["data-fp"]
	dataFpAttr = strings.Replace(dataFpAttr, "'", "\"", -1)

	res := struct {
		Defaults struct {
			Eur float64 `json:"EUR"`
		} `json:"defaults"`
	}{}
	if err := json.Unmarshal([]byte(dataFpAttr), &res); err != nil {
		panic(err)
	}

	return data{
		Title: title,
		Price: res.Defaults.Eur,
	}
}
