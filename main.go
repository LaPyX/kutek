package main

import (
	"kutek/parser"
)

func main() {
	kutek := parser.GetParser(parser.SiteKutekUrl)
	kutek.Run()

	kutekMood := parser.GetParser(parser.SiteKutekMoodUrl)
	kutekMood.SetItems(kutek.GetItems())
	kutekMood.Run()

	excel := &Excel{name: "Book1.xlsx"}
	excel.generate(kutekMood.GetItems())
}
