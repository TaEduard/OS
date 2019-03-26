package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
)

func main() {
	// Users Values which contains
	// an array of int values
	type Values struct {
		Start  int   `json:"startingPoint"`
		Values []int `json:"values"`
	}
	// Open our jsonFile
	jsonFile, err := os.Open("data.json")
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println("Successfully Opened data.json")
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Values array
	var values Values

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'values' which we defined above
	json.Unmarshal(byteValue, &values)
	var newArr []int
	// we iterate through every user within our values array and
	startingpointFoundInValues := false
	for i := 0; i < len(values.Values); i++ {
		if values.Values[i] == values.Start {
			newArr = values.Values[i:]
			newArr = append(newArr, values.Values[:i]...)
			fmt.Println(newArr)
			startingpointFoundInValues = true
		}
	}
	if !startingpointFoundInValues {
		if getMin(values.Values) > values.Start || getMax(values.Values) < values.Start {
			newArr = append(newArr, values.Start)
			newArr = append(newArr, values.Values...)
			fmt.Println(newArr)
		} else {
			for i := 0; i < len(values.Values); i++ {
				if values.Start < values.Values[i] {
					newArr = append(newArr, values.Start)
					newArr = append(newArr, values.Values[i:]...)
					newArr = append(newArr, values.Values[:i]...)
					fmt.Println(newArr)
				}
			}
		}
	}
	drawPlot(newArr)
}

func getMin(v []int) int {
	m := 0
	for i, e := range v {
		if i == 0 || e < m {
			m = e
		}
	}
	return m
}

func getMax(v []int) int {
	m := 0
	for i, e := range v {
		if i == 0 || e > m {
			m = e
		}
	}
	return m
}

func drawPlot(newArr []int) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "FCFS"
	p.X.Label.Text = "turn"
	p.Y.Label.Text = "Value"
	p.Y.Tick.Marker = commaTicks{}
	p.X.Tick.Marker = commaTicks{}

	err = plotutil.AddLinePoints(p, "FCFS", MakePoints(newArr))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(400, 400, "points.png"); err != nil {
		panic(err)
	}
}

// MakePoints returns points
func MakePoints(n []int) plotter.XYs {
	pts := make(plotter.XYs, len(n))
	for i := range pts {
		if i == 0 {
			pts[i].X = 1
		} else {
			pts[i].X = pts[i-1].X + 1
		}
		pts[i].Y = float64(n[i])
	}
	return pts
}

type commaTicks struct{}

// Ticks computes the default tick marks, but inserts commas
// into the labels for the major tick marks.
func (commaTicks) Ticks(min, max float64) []plot.Tick {

	tks := plot.DefaultTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" { // Skip minor ticks, they are fine.
			if i != len(tks) {
				tksint, _ := strconv.Atoi(tks[i-1].Label)
				tk := tksint + 1
				tks[i].Label = strconv.Itoa(tk)
			}
		}
		tks[i].Label = t.Label
	}
	return tks
}
