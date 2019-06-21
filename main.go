package main

import (
	"log"
	"net/http"
	"time"

	"github.com/d2r2/go-max7219"
)

func main() {
	// Create new LED matrix with number of cascaded devices is equal to 1
	mtx := max7219.NewMatrix(4)
	// Open SPI device with spibus and spidev parameters equal to 0 and 0.
	// Set LED matrix brightness is equal to 7
	err := mtx.Open(0, 0, 7)
	if err != nil {
		log.Fatal(err)
	}
	queue := make(chan string)
	defer mtx.Close()
	go func(mtx *max7219.Matrix, queue chan string) {
		i := "SMS Texter 2.0"
		for {
			select {
			case j := <-queue:
				i = j
				mtx.SlideMessage(i,
					max7219.FontCP437, true, 60*time.Millisecond)
				for index := 0; index < 20; index++ {
					mtx.Device.ScrollLeft(true)
					time.Sleep(60 * time.Millisecond)
				}
			default:
				mtx.SlideMessage(i,
					max7219.FontCP437, true, 60*time.Millisecond)
				for index := 0; index < 20; index++ {
					mtx.Device.ScrollLeft(true)
					time.Sleep(60 * time.Millisecond)
				}
			}
		}
	}(mtx, queue)
	http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
		queue <- r.FormValue("Body")
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
