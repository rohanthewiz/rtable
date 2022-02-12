package rtable

import "fyne.io/fyne/v2/data/binding"

type Animal struct {
	Name, Type, Color, Weight string
}

var animals = []Animal{
	{Name: "Frisky", Type: "cat", Color: "gray", Weight: "10"},
	{Name: "Ella", Type: "dog", Color: "brown", Weight: "50"},
	{Name: "Mickey", Type: "mouse", Color: "black", Weight: "1"},
}

var AnimalCols = []ColAttr{
	{ColName: "Name", Header: "Name", WidthPercent: 100},
	{ColName: "Type", Header: "Type", WidthPercent: 66},
	{ColName: "Color", Header: "Color", WidthPercent: 100},
	{ColName: "Weight", Header: "Weight", WidthPercent: 64},
}

var AnimalBindings []binding.DataMap

// Create a binding for each animal data
func init() {
	for i := 0; i < len(animals); i++ {
		AnimalBindings = append(AnimalBindings, binding.BindStruct(&animals[i]))
	}
}
