package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
)

const sheetWidth = 297
const sheetHeight = 210

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("??????")
	fmt.Fprintf(w, "home")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%#v", r)
	log.Println("????? not found ?")
	fmt.Fprintf(w, " 404 not found lol :P")
}

func tstHnd(w http.ResponseWriter, r *http.Request) {
	img := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
	encoder := png.Encoder{}
	gc := draw2dimg.NewGraphicContext(img)

	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})

	gc.BeginPath()
	gc.MoveTo(10, 10)
	gc.LineTo(100, 50)
	gc.QuadCurveTo(100, 10, 10, 10)
	gc.Close()
	gc.FillStroke()
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	gc.MoveTo(100, 100)
	gc.SetFillColor(image.Black)
	gc.SetFontSize(32)
	gc.FillStringAt("xsdfjklgh ksdjfghk sdjfghsdfjd", 40, 40)
	encoder.Encode(w, img)
}

var usernames []string
var numSheets = 1

func init() {
	usernames = []string{"nuuls", "fourtf", "gempir", "Lordborne"}
}

type BingoColumn struct {
	// -1 indicates FREE SPACE
	// a bingo row has 5 numbers
	Rows []int
}

const (
	numRows    = 5
	numColumns = 5
)

type BingoSheet struct {
	// 24 numbers for the sheet. middle spot is always free
	// a bingo sheet has 5 rows
	Columns []BingoColumn
}

func numberExistsInSet(needle int, set []int) bool {
	for _, number := range set {
		if needle == number {
			return true
		}
	}

	return false
}

func NewBingoColumn(column int) BingoColumn {
	bingoColumn := BingoColumn{
		Rows: make([]int, numRows),
	}

	generatedNumbers := make([]int, 0)

	for row := 0; row < numRows; row++ {
		number := rng.Intn(14) + 1 + (column * 15)

		for numberExistsInSet(number, generatedNumbers) {
			number = rng.Intn(14) + 1 + (column * 15)
		}
		generatedNumbers = append(generatedNumbers, number)
		bingoColumn.Rows[row] = number
	}

	return bingoColumn
}

func NewBingoSheet() *BingoSheet {
	bingoSheet := &BingoSheet{
		Columns: make([]BingoColumn, numColumns),
	}

	for column := 0; column < numColumns; column++ {
		bingoSheet.Columns[column] = NewBingoColumn(column)
	}

	return bingoSheet
}

type coordinate struct {
	Row    int
	Column int
}

func (c *coordinate) matches(row, column int) bool {
	return c.Row == row && c.Column == column
}

var FreeSlot = coordinate{
	Row:    2,
	Column: 2,
}

func createBingoSheetImage(sheet *BingoSheet) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, 297, 210.0))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold | draw2d.FontStyleItalic})
	gc.SetFillColor(image.Black)
	gc.SetFontSize(32)

	const offsetX = 60
	const offsetY = 40

	// Draw numbers
	for row := 0; row < numRows; row++ {
		for column := 0; column < numColumns; column++ {
			realX := float64(column * offsetX)
			realY := float64(40 + row*offsetY)
			if FreeSlot.matches(row, column) {
				gc.SetFillColor(color.RGBA{0xFF, 0x00, 0xFF, 0xFF})
				gc.FillStringAt("XD", realX, realY)
			} else {
				gc.SetFillColor(image.Black)
				gc.FillStringAt(strconv.Itoa(sheet.Columns[column].Rows[row]), realX, realY)
			}
		}
	}

	// Draw horizontal row
	for row := 1; row < numRows; row++ {
		y := float64(row * sheetHeight / numRows)
		gc.MoveTo(0, y)
		gc.LineTo(sheetWidth, y)
		gc.Stroke()
	}

	// Draw vertical columns
	for column := 1; column < numColumns; column++ {
		x := float64(column * sheetWidth / numColumns)
		gc.MoveTo(x, 0)
		gc.LineTo(x, sheetHeight)
		gc.Stroke()
	}

	return img
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Generate bingo images xd")

	for _, username := range usernames {
		for sheetIndex := 1; sheetIndex <= numSheets; sheetIndex++ {

			sheet := NewBingoSheet()

			img := createBingoSheetImage(sheet)
			draw2dimg.SaveToPngFile(fmt.Sprintf("images/%s-%d.png", username, sheetIndex), img)
		}
	}

	fmt.Fprintf(w, "xd generated images")
}

func sheetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	isValid := false
	for _, validUsername := range usernames {
		if validUsername == username {
			isValid = true
			break
		}
	}

	if !isValid {
		fmt.Fprintf(w, "<h1>Invalid username</h1>")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html>")
	fmt.Fprintf(w, "<h1>%s</h1>", username)
	for i := 1; i <= numSheets; i++ {
		fmt.Fprintf(w, "Sheet %d:<br>", i)
		fmt.Fprintf(w, "<img src=\"/bingo/images/"+username+"-"+strconv.Itoa(i)+".png\"><br>")
	}
	fmt.Fprintf(w, "</html>")
}

var rng *rand.Rand

func main() {
	rng = rand.New(rand.NewSource(123))
	fmt.Println("server up")
	r := mux.NewRouter()
	s := r.PathPrefix("/bingo").Subrouter()
	s.HandleFunc("/", HomeHandler)
	s.HandleFunc("", HomeHandler)
	s.HandleFunc("/test", tstHnd)
	s.HandleFunc("/sheets/{username}", sheetsHandler)
	s.HandleFunc("/generate", generateHandler)
	s.PathPrefix("/images/").Handler(
		http.StripPrefix("/bingo/images/", http.FileServer(http.Dir("./images/"))),
	)

	r.NotFoundHandler = http.HandlerFunc(notFound)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
