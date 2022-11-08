package parser

import (
	"fmt"
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
	body := k.client.request(k.url)

	doc := k.client.queryDoc(string(body))

	fmt.Println(doc)
	//log.Fatal("stop")
}
