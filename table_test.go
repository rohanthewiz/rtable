package rtable

// Code by Rohan Allison
import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"github.com/rohanthewiz/rerr"
)

func TestRTable(t *testing.T) {
	ap := app.New()
	wn := ap.NewWindow("Table Widget")
	if wn == nil {
		t.Log("Unable to create window")
		t.Fail()
	}

	tblOpts := &TableOptions{
		RefWidth: "reference width",
		ColAttrs: AnimalCols,
		Bindings: AnimalBindings,
	}

	tbl := CreateTable(tblOpts,
		func(cell widget.TableCellID) {
			if cell.Row == 0 && cell.Col >= 0 && cell.Col < len(AnimalCols) { // header cell
				fmt.Println("-->", tblOpts.ColAttrs[cell.Col].Header)
				return
			}
			// Other rows
			str, err := GetStrCellValue(cell, tblOpts)
			if err != nil {
				fmt.Println(rerr.StringFromErr(err))
				return
			}
			fmt.Println("-->", str)
		})

	if tbl == nil {
		t.Log("Table was not properly created")
		t.Fail()
	}
}
