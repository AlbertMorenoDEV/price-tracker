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

func main() {
	d := getData("https://www.vivobarefoot.com/es/mens/everyday/gobi-ii-mens?colour=Black")

	fmt.Println(d)
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
