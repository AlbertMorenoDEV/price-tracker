package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
)

type (
	pageUrl string
	data    struct {
		Title string
		Price float64
	}
	response struct {
		dat data
		err error
	}
)

var pagesUrl = []pageUrl{
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
	workc := make(chan pageUrl, 8)
	resultc := make(chan response, 8)

	go func() {
		for {
			select {
			case res := <-resultc:
				if res.err != nil {
					log.Println(res.err)
					continue
				}

				fmt.Println(res.dat)
			default:
			}
		}
	}()

	for i := 0; i < 4; i++ {
		go func() {
			for {
				select {
				case pag := <-workc:
					resultc <- getResponse(pag)
				default:
				}
			}
		}()
	}

	for _, u := range pagesUrl {
		fmt.Println("URL: " + u)
		workc <- u

		time.Sleep(5 * time.Second)
	}

	fmt.Println("Done!")
}

func getResponse(u pageUrl) response {
	resp, err := soup.Get(string(u))
	if err != nil {
		return response{
			err: err,
		}
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
		return response{
			err: err,
		}
	}

	return response{
		dat: data{
			Title: title,
			Price: res.Defaults.Eur,
		},
		err: nil,
	}
}
