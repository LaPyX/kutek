package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	SiteKutekUrl     = "https://kutek.com.pl/kolekcje/"
	SiteKutekMoodUrl = "http://kutekmood.com/"
)

type Item struct {
	url     string
	name    string
	article string
	width   string
	height  string
	lamp    string
	colors  []string
	img     string
}

type kutekSiteParse struct {
	url    string
	client *http.Client
	items  map[int]*Item
}

func main() {
	kutek := &kutekSiteParse{
		url:    SiteKutekUrl,
		client: newClient(),
		items:  make(map[int]*Item),
	}

	kutek.Run()
}

func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				// See comment above.
				// UNSAFE!
				// DON'T USE IN PRODUCTION!
				InsecureSkipVerify: true,
			},
		},
	}
}

func (k kutekSiteParse) Run() {
	body := k.request(k.url)

	doc := queryDoc(string(body))

	doc.Find(".collections-container > a").Each(func(_ int, a *goquery.Selection) {
		url, _ := a.Attr("href")

		b := k.request(url)
		if b != nil {
			d := queryDoc(string(b))

			d.Find(".collections-container > a").Each(func(_ int, a *goquery.Selection) {
				url, _ := a.Attr("href")

				b := k.request(url)
				if b != nil {
					item := queryDoc(string(b))

					desc := item.Find(".product-descritpion")

					name := desc.Find("h1").First().Text()
					article := desc.Find("h2").First().Text()
					height := desc.Find("ul > li").Eq(0).Text()
					width := desc.Find("ul > li").Eq(1).Text()
					lamp := desc.Find("ul > li").Eq(2).Text()
					img, _ := item.Find(".gallery-image img").First().Attr("src")

					var colors []string
					item.Find(".product-galvanisation > h3").Each(func(i int, h3 *goquery.Selection) {
						if i == 0 {
							return
						}
						colors = append(colors, strings.TrimSpace(h3.Text()))
					})

					k.items[len(k.items)] = &Item{
						url:     url,
						name:    name,
						article: article,
						height:  height,
						width:   width,
						lamp:    lamp,
						img:     img,
						colors:  colors,
					}

					fmt.Println(k.items[len(k.items)-1])

					excel := &Excel{name: "Book1.xlsx"}
					excel.generate(k.items)
				}
			})
		}

		log.Fatal("stop")
	})
}

func (k *kutekSiteParse) request(uri string) []byte {
	method := "GET"

	payload := &bytes.Buffer{}

	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	res, err := k.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//fmt.Println(res.StatusCode)

	return body
}

func queryDoc(data string) *goquery.Document {
	node, err := html.Parse(strings.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	doc := goquery.NewDocumentFromNode(node)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}
