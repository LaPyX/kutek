package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type KutekSiteParser struct {
	client *SiteClient
	url    string
	items  map[int]*Item
}

func (k *KutekSiteParser) SetItems(items map[int]*Item) {
	k.items = items
}

func (k *KutekSiteParser) GetItems() map[int]*Item {
	return k.items
}

func (k *KutekSiteParser) Run() {
	body := k.client.request(k.url)

	doc := k.client.queryDoc(string(body))

	doc.Find(".collections-container > a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		url, _ := a.Attr("href")

		b := k.client.request(url)
		if b != nil {
			d := k.client.queryDoc(string(b))

			d.Find(".collections-container > a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
				url, _ := a.Attr("href")

				b := k.client.request(url)
				if b != nil {
					item := k.client.queryDoc(string(b))

					desc := item.Find(".product-descritpion")

					name := desc.Find("h1").First().Text()
					article := desc.Find("h2").First().Text()
					li := desc.Find("ul > li")
					height := li.Eq(0).Text()
					width := li.Eq(1).Text()

					var lamp, distance string
					if li.Length() > 3 {
						distance = li.Eq(2).Text()
						lamp = li.Eq(3).Text()
					} else {
						lamp = li.Eq(2).Text()
					}

					img, _ := item.Find(".gallery-image img").First().Attr("src")

					var colors []string
					item.Find(".product-galvanisation > h3").Each(func(i int, h3 *goquery.Selection) {
						if i == 0 {
							return
						}
						colors = append(colors, strings.TrimSpace(h3.Text()))
					})

					k.items[len(k.items)] = &Item{
						Url:      url,
						Name:     name,
						Article:  article,
						Height:   height,
						Width:    width,
						Distance: distance,
						Lamp:     lamp,
						Img:      img,
						Colors:   colors,
					}

					fmt.Println(k.items[len(k.items)-1])
				}

				return true
			})
		}

		return false
	})
}
