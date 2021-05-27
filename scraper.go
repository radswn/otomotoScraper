package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func getScrapedData(w http.ResponseWriter, r *http.Request) {

	var input InputUrl
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := input.Url + "&page=" + strconv.Itoa(input.Page)
	log.Println("visiting", url)

	c := colly.NewCollector()
	var response []Offer
	replacer := strings.NewReplacer(" ", "", "km", "", "cm", "")

	c.OnHTML("article", func(e *colly.HTMLElement) {

		nameBlock := e.DOM.Find("h2").Find("a")

		offer := &Offer{
			Name:        strings.TrimSpace(nameBlock.Text()),
			Description: e.ChildText("h3"),
			Url:         nameBlock.AttrOr("href", ""),
		}

		jpgUrl := e.DOM.Find("img").AttrOr("data-srcset", "")
		jpgUrl = jpgUrl[:len(jpgUrl)-5]
		offer.JpgUrl = jpgUrl

		e.DOM.Find("span").Each(func(i int, selection *goquery.Selection) {
			if selection.AttrOr("class", "") == "offer-price__number ds-price-number" {
				offer.Price, _ = strconv.Atoi(strings.Replace(selection.Children().First().Text(), " ", "", -1))
			}
		})

		e.DOM.Find("ul").Children().Each(func(i int, selection *goquery.Selection) {

			switch i {
			case 0:
				offer.Year, _ = strconv.Atoi(strings.TrimSpace(selection.Text()))
			case 1:
				mileage := strings.TrimSpace(selection.Text())
				offer.Mileage, _ = strconv.Atoi(replacer.Replace(mileage))
			case 2:
				engine := strings.TrimSpace(selection.Text())
				engine = engine[:len(engine)-1]
				offer.Engine, _ = strconv.Atoi(replacer.Replace(engine))
			case 3:
				offer.EngineType = strings.TrimSpace(selection.Text())
			}
		})

		location := ""
		e.DOM.Find("h4").Children().Each(func(i int, selection *goquery.Selection) {
			location += selection.Text() + " "
		})
		offer.Location = strings.TrimSpace(location)

		response = append(response, *offer)
	})

	//Command to visit the website
	c.Visit(url)

	// parse our response slice into JSON format
	b, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	// Add some header and write the body for our endpoint
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
