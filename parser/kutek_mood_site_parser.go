package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

type KutekMoodSiteParser struct {
	client *SiteClient
	url    string
	items  map[int]*Item
}

func (k *KutekMoodSiteParser) SetItems(items map[int]*Item) {
	k.items = items
}

func (k *KutekMoodSiteParser) GetItems() map[int]*Item {
	return k.items
}

func (k *KutekMoodSiteParser) Run() {
	body := k.client.request(k.getUrl("/en/2-offer"))

	doc := k.client.queryDoc(string(body))

	doc.Find(".records.ibs .item > a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
		url, _ := a.Attr("href")

		b := k.client.request(k.getUrl(url))
		if b != nil {
			d := k.client.queryDoc(string(b))

			d.Find(".records.ibs .item > a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
				url, _ := a.Attr("href")
				url = k.getUrl(url)
				//url = "http://www.kutekmood.com/en/p77-aba-zw-8-n-ii"

				b := k.client.request(url)
				if b != nil {
					item := k.client.queryDoc(string(b))
					name := item.Find(".grp > .name").First().Text()
					_type := item.Find(".grp > .type").First().Text()
					img, _ := item.Find(".imgs > a").First().Attr("href")
					d, _ := item.Find(".grp .con > p").First().Html()

					desc := strings.Split(d, "<br/>")
					reg, _ := regexp.Compile("<strong>(.*)</strong>")
					height := reg.ReplaceAllString(desc[0], "")
					width := reg.ReplaceAllString(desc[1], "")

					var bulbs, distance string
					regBulbs, _ := regexp.Compile("Bulbs")
					if len(desc) > 3 && regBulbs.MatchString(desc[3]) {
						distance = reg.ReplaceAllString(desc[2], "")
						bulbs = reg.ReplaceAllString(desc[3], "")
					} else {
						bulbs = reg.ReplaceAllString(desc[2], "")
					}

					var colors []string
					item.Find(".grp .con > h5 > span").Each(func(i int, h5 *goquery.Selection) {
						colors = append(colors, strings.TrimSpace(h5.Text()))
					})

					var colorShades []string
					item.Find(".grp > .colors i").Each(func(i int, t *goquery.Selection) {
						colorShades = append(colorShades, strings.TrimSpace(t.Text()))
					})

					k.items[len(k.items)] = &Item{
						Url:         url,
						Name:        name,
						Article:     name,
						Type:        _type,
						Height:      height,
						Width:       width,
						Lamp:        bulbs,
						Distance:    distance,
						Img:         k.getUrl(img),
						Colors:      colors,
						ColorShades: colorShades,
					}

					fmt.Println(k.items[len(k.items)-1])
				}

				return true
			})
		}

		return false
	})
}

func (k *KutekMoodSiteParser) getUrl(url string) string {
	return k.url + url
}
