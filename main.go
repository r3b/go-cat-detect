package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

var faceFound bool = false

func main() {
	deviceID := 0
	xmlFile := "face.xml"

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("I C U")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	red := color.RGBA{255, 0, 0, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Cannot open classifier file: %v\n", xmlFile)
		return
	}

	fmt.Printf("Open default camera: %v\n", deviceID)
	fmt.Println("Press any key to exit")
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("error reading device: %d\n", deviceID)
			// return
		}
		if img.Empty() {
			continue
		}

		// Have a look
		rects := classifier.DetectMultiScale(img)
		if len(rects) > 0 && !faceFound {
			faceFound = true
			fmt.Printf("found %d faces\n", len(rects))
		}

		for i, r := range rects {
			err := gocv.Rectangle(&img, r.Inset(-20), red, 3)
			if err != nil {
				fmt.Println("Rectangle Error: ", err)
			}
			label := fmt.Sprintf("Somebody #%d", i)
			size := gocv.GetTextSize(label, gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			err = gocv.PutText(&img, label, pt, gocv.FontHersheyPlain, 1.2, red, 2)
			if err != nil {
				fmt.Println("PutText Error: ", err)
			}
		}

		// show the image in the window, and wait 1 millisecond
		err = window.IMShow(img)
		if err != nil {
			fmt.Println("IMShow Error: ", err)
		}
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
