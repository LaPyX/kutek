package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"kutek/parser"
	"strconv"
	"strings"
)

var _Cells = [12]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}

type Excel struct {
	name       string
	file       *excelize.File
	currentRow int
}

func (e *Excel) generate(items map[int]*parser.Item) {
	sheet := "Sheet1"

	e.file = excelize.NewFile()
	// Create a new sheet.
	e.file.NewSheet(sheet)

	// set headers
	e.setHeader(sheet, []string{
		"Name",
		"Type",
		"Article",
		"Img",
		"Width",
		"Height",
		"Distance",
		"Lamp",
		"Colors",
		"ColorShades",
		"Url",
	})

	// Set value of a cell.
	for _, item := range items {
		e.setRow(sheet, []string{
			item.Name,
			item.Type,
			item.Article,
			item.Img,
			item.Width,
			item.Height,
			item.Distance,
			item.Lamp,
			strings.Join(item.Colors, ", "),
			strings.Join(item.ColorShades, ", "),
			item.Url,
		})
	}
	// Save spreadsheet by the given path.
	if err := e.file.SaveAs(e.name); err != nil {
		fmt.Println(err)
	}
}

func (e *Excel) setHeader(sheet string, cells []string) {
	e.currentRow = 1
	style, _ := e.file.NewStyle(`{"font":{"bold":true,"size":14}}`)
	e.file.SetRowStyle(sheet, 1, 1, style)
	for i, name := range cells {
		e.file.SetCellValue(sheet, e.getChar(i), name)
	}
	e.nextRow()
}

func (e *Excel) setRow(sheet string, cells []string) {
	for i, value := range cells {
		e.file.SetCellValue(sheet, e.getChar(i), strings.TrimSpace(value))
	}
	e.nextRow()
}

func (e *Excel) getChar(i int) string {
	chr := _Cells[i]
	return chr + strconv.Itoa(e.currentRow)
}

func (e *Excel) nextRow() {
	e.currentRow++
}
