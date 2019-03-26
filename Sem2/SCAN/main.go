package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
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
	sort.Ints(values.Values)
	newArr = shortestSeekTimeFirst(values.Values, values.Start)
	drawPlot(newArr)
}

func shortestSeekTimeFirst(request []int, head int) []int {
	if len(request) == 0 {
		return nil
	}

	l := len(request)
	diff := make([][]int, l)
	// initialize array
	for i, _ := range diff {
		diff[i] = make([]int, 2)

	}

	// stores sequence in which disk
	// access is done
	seek_sequence := make([]int, l+1)

	for i := 0; i < l; i++ {
		seek_sequence[i] = head
		calculateDifference(request, head, diff)
		index := findMin(diff)

		diff[index][1] = 1
		head = request[index]
	}
	//  for last accessed track
	seek_sequence[len(seek_sequence)-1] = head

	if contains(request, head) {
		// Remove the element at index 0 from seek_sequence.
		copy(seek_sequence[0:], seek_sequence[1:])           // Shift a[i+1:] left one index.
		seek_sequence[len(seek_sequence)-1] = 0              // Erase last element (write zero value).
		seek_sequence = seek_sequence[:len(seek_sequence)-1] // Truncate slice.
	}

	if contains(request, head) {
		for i := 0; i < l; i++ {
			fmt.Println(seek_sequence[i])
		}
	} else {
		for i := 0; i <= l; i++ {
			fmt.Println(seek_sequence[i])
		}
	}
	return seek_sequence
}

func calculateDifference(queue []int, head int, diff [][]int) {
	for i, _ := range diff {
		diff[i][0] = int(math.Abs(float64(queue[i] - head)))
	}
}
func findMin(diff [][]int) int {
	index := -1
	min := 999999999
	for i, _ := range diff {
		if diff[i][1] != 1 && min > diff[i][0] {
			min = diff[i][0]
			index = i
		}
	}
	return index
}
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
