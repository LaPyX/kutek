package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

var _Cells = [8]string{"A", "B", "C", "D", "E", "F", "G", "H"}

type Excel struct {
	name       string
	file       *excelize.File
	currentRow int
}

func (e *Excel) generate(items map[int]*Item) {
	sheet := "Sheet1"

	e.file = excelize.NewFile()
	// Create a new sheet.
	e.file.NewSheet(sheet)

	// set headers
	e.setHeader(sheet, []string{
		"Name",
		"Article",
		"Img",
		"Width",
		"Height",
		"Lamp",
		"Colors",
		"Url",
	})

	// Set value of a cell.
	for _, item := range items {
		e.setRow(sheet, []string{
			item.name,
			item.article,
			item.img,
			item.width,
			item.height,
			item.lamp,
			strings.Join(item.colors, ", "),
			item.url,
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
		e.file.SetCellValue(sheet, e.getChar(i), value)
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
