package parser

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

type KutekSiteParser struct {
	client *SiteClient
	url    string
	items  Items
}

func (k *KutekSiteParser) SetItems(items Items) {
	k.items = items
}

func (k *KutekSiteParser) GetItems() Items {
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

			reg, _ := regexp.Compile("\\([A-Z]{1,2}\\)")
			regSlash, _ := regexp.Compile("\\([A-Z]{1,2}\\/")
			regShades, _ := regexp.Compile("\\)[A-Z]{1,2}")
			regColor, _ := regexp.Compile("\\(.*\\)")

			d.Find(".collections-container > a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
				url, _ := a.Attr("href")
				//url = "https://kutek.com.pl/kolekcje/arezzo/are-pl-8n300/" // ARE-PL-8(P)300
				//url = "https://kutek.com.pl/kolekcje/genova/gen-k-1-zma/" // GEN-K-1(ZM/A)
				//url = "https://kutek.com.pl/kolekcje/kamma/kam-zw-6-p/" // KAM-ZW-6(P)
				//url = "https://kutek.com.pl/kolekcje/elve/elv-pods-1z-380-dl/" // ELV-PODS-1-(P)-380-DL
				//url = "https://kutek.com.pl/kolekcje/bellagio/bel-pl-3n470-cr/" // BEL-PL-3(P)470-CR
				//url = "https://kutek.com.pl/kolekcje/ellini/ell-plm-6bn350ii/" // ELL-PLM-6(BN)350/II
				//url = "https://kutek.com.pl/kolekcje/bibione-abazur/bib-zw-6-pasr/" // BIB-ZW-6 (P/A)SR

				b := k.client.request(url)
				if b != nil {
					query := k.client.queryDoc(string(b))

					desc := query.Find(".product-descritpion")

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

					img, _ := query.Find(".gallery-image img").First().Attr("src")

					var colorShades, colors []string
					query.Find(".product-galvanisation > h3").Each(func(i int, h3 *goquery.Selection) {
						text := strings.TrimSpace(h3.Text())
						text = regColor.FindString(text)

						if text == "" {
							return
						}

						text = strings.Trim(text, "()")
						if h3.Find("img").Length() > 0 {
							colors = append(colors, text)
						} else {
							colorShades = append(colorShades, text)
						}
					})

					item := Item{
						Url:         url,
						Name:        name,
						Article:     strings.ReplaceAll(article, " ", ""),
						Height:      height,
						Width:       width,
						Distance:    distance,
						Lamp:        lamp,
						Img:         img,
						Colors:      colors,
						ColorShades: colorShades,
					}

					fmt.Println(item)

					if len(colors) == 0 {
						k.items[item.Article] = &item
						return true
					}

					// цвет
					for _, color := range colors {
						if reg.MatchString(article) {
							article = reg.ReplaceAllString(article, "("+color+")")
						} else if regSlash.MatchString(article) {
							article = regSlash.ReplaceAllString(article, "("+color+"/")
						} else {
							continue
						}

						if len(colorShades) == 0 {
							clone := item
							clone.Article = article
							k.items[clone.Article] = &clone
							continue
						}

						// оттенки
						for _, shade := range colorShades {
							clone := item
							clone.Article = regShades.ReplaceAllString(clone.Article, ")"+shade)
							k.items[clone.Article] = &clone
						}
					}

				}

				return true
			})
		}

		return true
	})
}
