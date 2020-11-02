package main

import (
	"github.com/fogleman/gg"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// http://www.mini.pw.edu.pl/MiNIwyklady/fraktale/Dywan/dywan.html
func main() {

	r := mux.NewRouter()
	handler := http.HandlerFunc(sierpinskisCarpetHandler)
	r.HandleFunc("/sierpinski/carpet/{step}", handler).Methods(http.MethodGet)

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}

func sierpinskisCarpetHandler(rw http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["step"]
	step, err := strconv.Atoi(s)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("Step variable must be an integer"))
		return
	}

	dc := gg.NewContext(1200, 1200)
	sierpinskisCarpet(1200, step, dc)

	rw.Header()["Content-Type"] = []string{"image/png"}
	dc.EncodePNG(rw)
}

func sierpinskisCarpet(size float64, stepCount int, dc *gg.Context) {

	// Draw initial square
	dc.DrawRectangle(0, 0, size, size)
	dc.SetRGBA(0, 0, 0, 1)
	dc.Fill()

	// Start recursive steps
	sierpinskisCarpetStep(0, 0, size, stepCount, dc)
}

func sierpinskisCarpetStep(x, y, size float64, step int, dc *gg.Context) {

	// Recursion stop condition
	if step == 0 {
		return
	}

	// Plot white triangle in the middle
	newSize := size / 3
	newStep := step - 1

	dc.DrawRectangle(x+newSize, y+newSize, newSize, newSize)
	dc.SetRGBA(1, 1, 1, 1)
	dc.Fill()

	// Recursive calls for 8 squares
	sierpinskisCarpetStep(x, y, newSize, newStep, dc)
	sierpinskisCarpetStep(x+newSize, y, newSize, newStep, dc)
	sierpinskisCarpetStep(x+2*newSize, y, newSize, newStep, dc)

	sierpinskisCarpetStep(x, y+newSize, newSize, newStep, dc)
	sierpinskisCarpetStep(x+2*newSize, y+newSize, newSize, newStep, dc)

	sierpinskisCarpetStep(x, y+2*newSize, newSize, newStep, dc)
	sierpinskisCarpetStep(x+newSize, y+2*newSize, newSize, newStep, dc)
	sierpinskisCarpetStep(x+2*newSize, y+2*newSize, newSize, newStep, dc)
}
