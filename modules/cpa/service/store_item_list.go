package service

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func get_store_item_list() {

	doc, err := goquery.NewDocument("http://stores.ebay.com.au/aishang2012")
	if err != nil {
		log.Println(err)
	}

	doc.Find("#lvc > table > tbody > tr > td.details > div.ttl").Each(func(i int, s *goquery.Selection) {
		title_link := s.Find("a")
		title := title_link.Text()
		link, _ := title_link.Attr("href")
		fmt.Printf("%02d: %s %s\n", i, title, link)
	})
}
