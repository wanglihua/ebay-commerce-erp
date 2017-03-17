package service

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func Get_store_item_list() {

	doc, err := goquery.NewDocument("http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=aishang2012&_pgn=1&_skc=0&rt=nc")
	if err != nil {
		log.Println(err)
	}

	doc.Find("ul#ListViewInner li h3.lvtitle").Each(func(i int, s *goquery.Selection) {
		title_link := s.Find("a")
		title := title_link.Text()
		link, _ := title_link.Attr("href")
		fmt.Printf("%02d: %s %s\n", i, title, link)
	})
}

/*

http://www.ebay.com.au/usr/aishang2012

http://www.ebay.com.au/sch/aishang2012/m.html?_nkw=&_armrs=1&_ipg=&_from=&rt=nc&LH_BIN=1
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=aishang2012&_pgn=1&_skc=0&rt=nc
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=aishang2012&_pgn=2&_skc=50&rt=nc
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=aishang2012&_pgn=3&_skc=100&rt=nc
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=aishang2012&_pgn=4&_skc=150&rt=nc


http://www.ebay.com.au/usr/sparesbox_auto

http://www.ebay.com.au/sch/sparesbox_auto/m.html?_nkw=&_armrs=1&_ipg=&_from=&rt=nc&LH_BIN=1
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=sparesbox_auto&_pgn=1&_skc=0&rt=nc
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=sparesbox_auto&_pgn=2&_skc=50&rt=nc
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=sparesbox_auto&_pgn=3&_skc=100&rt=nc
http://www.ebay.com.au/sch/m.html?_nkw=&_armrs=1&_from=&LH_BIN=1&_ssn=sparesbox_auto&_pgn=4&_skc=150&rt=nc


 */
