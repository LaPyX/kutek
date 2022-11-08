package parser

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
	Url     string
	Name    string
	Article string
	Width   string
	Height  string
	Lamp    string
	Colors  []string
	Img     string
}

type SiteParser interface {
	Run()
	GetItems() map[int]*Item
	SetItems(map[int]*Item)
}

type SiteClient struct {
	client *http.Client
}

func (s *SiteClient) request(uri string) []byte {
	method := "GET"

	payload := &bytes.Buffer{}

	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	res, err := s.client.Do(req)
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

func (s *SiteClient) queryDoc(data string) *goquery.Document {
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

func newClient() *SiteClient {
	return &SiteClient{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// See comment above.
					// UNSAFE!
					// DON'T USE IN PRODUCTION!
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func GetParser(url string) SiteParser {
	switch url {
	case SiteKutekUrl:
		return &KutekSiteParser{
			url:    url,
			client: newClient(),
			items:  make(map[int]*Item),
		}
	case SiteKutekMoodUrl:
		return &KutekMoodSiteParser{
			url:    url,
			client: newClient(),
			items:  make(map[int]*Item),
		}
	}

	return nil
}
