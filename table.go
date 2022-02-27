package rtable

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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

func CreateTable(tblOpts *TableOptions, clkHdlr ...func(cell widget.TableCellID)) (tbl *widget.Table) {
	tbl = widget.NewTable(
		// Dimensions (rows, cols)
		func() (int, int) {
			return len(tblOpts.Bindings) + 1, len(tblOpts.ColAttrs) // + 1 row for a hdr
		},
		// Default
		func() fyne.CanvasObject {
			chk := widget.NewCheck("", func(c bool) {
				// log.Println("Chk Clicked")
			})
			ctr := container.NewMax(chk, widget.NewLabel(""))
			chk.Hide()
			return ctr
		},
		// Set Values
		func(position widget.TableCellID, cvObj fyne.CanvasObject) {
			// Type Switch
			// cvObj.(*widget.Label).Bind(datum.(binding.String))
			con := cvObj.(*fyne.Container)
			for _, conObj := range con.Objects {
				switch obj := conObj.(type) {
				case *widget.Label:
					if position.Row == 0 { // header row
						obj.Alignment = fyne.TextAlignCenter
						obj.TextStyle = fyne.TextStyle{Bold: true}
						obj.SetText(tblOpts.ColAttrs[position.Col].Header)
						return
					}
					if position.Col == 0 { // checkboxes only in 1st col
						obj.Hide()
						return
					}
					datum, err := getTableDatum(position, tblOpts)
					if err != nil {
						fmt.Println(rerr.StringFromErr(err))
						return
					}
					obj.Bind(datum.(binding.String))
				case *widget.Check:
					if position.Row > 0 && position.Col == 0 {
						obj.SetChecked(true) // hard-wired true for now
						obj.OnChanged = func(b bool) {
							fmt.Println("Clicked =-> rowIdx:", position.Row, "colIdx", position.Col)
						}
						obj.Show()
					} else {
						obj.Hide() // may not be necessary, but making sure
					}
				}
			}

			if position.Row == 0 { // header row
				label, ok := cvObj.(*widget.Label) // Maybe the best approach is to do the type switch above 1st
				if ok {
					label.Alignment = fyne.TextAlignCenter
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.SetText(tblOpts.ColAttrs[position.Col].Header)
				}
				return
			}
			// Get the datum for the non-hdr positions
			// for no binding just use SetText -> cvObj.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	if len(clkHdlr) > 0 {
		tbl.OnSelected = clkHdlr[0]
	}

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
