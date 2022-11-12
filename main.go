package main

import "kutek/parser"

func main() {
	excelReader := &ExcelPriceRead{
		name:   "prices.xlsx",
		prices: make(map[string]string),
	}
	excelReader.read()

	kutek := parser.GetParser(parser.SiteKutekUrl)
	kutek.Run()

	//kutekMood := parser.GetParser(parser.SiteKutekMoodUrl)
	//kutekMood.SetItems(kutek.GetItems())
	//kutekMood.Run()
	//

	excel := &Excel{
		name:   "Book1.xlsx",
		prices: excelReader.prices,
	}
	excel.generate(kutek.GetItems())
}
