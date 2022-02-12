package rtable

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
)

type TableOptions struct {
	RefWidth string
	ColAttrs []ColAttr
	Bindings []binding.DataMap
}

type ColAttr struct {
	ColName      string
	Header       string
	WidthRef     string
	WidthPercent int
}

func CreateTable(tblOpts *TableOptions, clkHdlr func(cell widget.TableCellID)) (tbl *widget.Table) {
	tbl = widget.NewTable(
		// Dimensions (rows, cols)
		func() (int, int) {
			return len(tblOpts.Bindings) + 1, len(tblOpts.ColAttrs) // + 1 row for a hdr
		},
		// Default
		func() fyne.CanvasObject {
			return widget.NewLabel(" - ")
		},
		// Set Values
		func(cell widget.TableCellID, cnvObj fyne.CanvasObject) {
			// for no binding just use SetText -> cnvObj.(*widget.Label).SetText(data[i.Row][i.Col])
			if cell.Row == 0 { // header row
				label := cnvObj.(*widget.Label)
				label.Alignment = fyne.TextAlignCenter
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.SetText(tblOpts.ColAttrs[cell.Col].Header)
				return
			}

			datum, err := getTableDatum(cell, tblOpts)
			if err != nil {
				fmt.Println(rerr.StringFromErr(err))
				return
			}
			cnvObj.(*widget.Label).Bind(datum.(binding.String))
		})

	tbl.OnSelected = clkHdlr

	refWidth := widget.NewLabel(tblOpts.RefWidth).MinSize().Width

	// Set Column widths
	for i, colAttr := range tblOpts.ColAttrs {
		tbl.SetColumnWidth(i, float32(colAttr.WidthPercent)/100.0*refWidth)
	}
	return
}

func GetStrCellValue(cell widget.TableCellID, tblOpts *TableOptions) (str string, err error) {
	datum, err := getTableDatum(cell, tblOpts)
	if err != nil {
		return str, rerr.Wrap(err)
	}

	str, err = datum.(binding.String).Get()
	if err != nil {
		return str, rerr.Wrap(err)
	}
	return
}

func getTableDatum(cell widget.TableCellID, tblOpts *TableOptions,
) (datum binding.DataItem, err error) {
	// Bounds check
	if cell.Row < 0 || cell.Row > len(tblOpts.Bindings) { // hdr is first row
		msg := "No data binding for row"
		log.Println(msg, cell.Row)
		return datum, rerr.NewRErr(msg)
	}
	if cell.Col < 0 || cell.Col > len(tblOpts.ColAttrs)-1 {
		return datum, rerr.NewRErr(fmt.Sprintf("Column %d is out of bounds", cell.Col))
	}

	// Get the data binding for the row
	bndg := tblOpts.Bindings[cell.Row-1]

	datum, err = bndg.GetItem(tblOpts.ColAttrs[cell.Col].ColName)
	if err != nil {
		log.Println(rerr.StringFromErr(rerr.Wrap(err)))
		return datum, rerr.Wrap(err)
	}
	return
}
