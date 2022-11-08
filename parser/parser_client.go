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
	"time"
)

const (
	SiteKutekUrl     = "https://kutek.com.pl/kolekcje/"
	SiteKutekMoodUrl = "http://www.kutekmood.com"
)

type Item struct {
	Url         string
	Name        string
	Article     string
	Width       string
	Height      string
	Distance    string
	Lamp        string
	Colors      []string
	Img         string
	Type        string
	ColorShades []string
}

type SiteParser interface {
	Run()
	GetItems() map[int]*Item
	SetItems(map[int]*Item)
}

type SiteClient struct {
	client   *http.Client
	interval time.Duration
}

func (s *SiteClient) request(uri string) []byte {
	if s.interval > 0 {
		time.Sleep(s.interval)
	}

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

func (s *SiteClient) SetInterval(interval time.Duration) *SiteClient {
	s.interval = interval
	return s
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
			client: newClient().SetInterval(100 * time.Millisecond),
			items:  make(map[int]*Item),
		}
	}

	return nil
}
