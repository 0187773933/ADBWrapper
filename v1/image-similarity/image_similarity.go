package image_similarity

import (
	// "fmt"
	"math"
	"gocv.io/x/gocv"
	"image"
	"image/color"
)

// https://github.com/imba28/image-similarity

const (
	hueBins        = 12
	saturationBins = 12
	brightnessBins = 4
)

func featuresInSegment(img gocv.Mat, mask gocv.Mat, hist gocv.Mat) ([]float64, error) {
	// h, s, v channel
	channels := []int{0, 1, 2}
	bins := []int{hueBins, saturationBins, brightnessBins}
	// hue range = 0-180, saturation range = 0-256, value/brightness range = 0-256
	ranges := []float64{0, 180, 0, 256, 0, 256}

	gocv.CalcHist([]gocv.Mat{img}, channels, mask, &hist, bins, ranges, false)
	gocv.Normalize(hist, &hist, 1, 0, gocv.NormL2)

	float64Hist := gocv.NewMat()
	defer float64Hist.Close()
	hist.ConvertTo(&float64Hist, gocv.MatTypeCV64F)
	f, err := float64Hist.DataPtrFloat64()

	// copy slice to golang memory, because after returning from this function the allocated memory is released.
	// accessing the underlying array would result in undefined behaviour.
	fCopy := make([]float64, len(f))
	copy(fCopy, f)
	return fCopy, err
}

func GetFeatureVectorFromFilePath( image_path string ) ( features []float64 ) {
	img := gocv.IMRead( image_path , gocv.IMReadColor)
	if img.Empty() { panic( "image is empty" ) }
	defer img.Close()

	// convert img to hsv color space
	// img.ConvertTo(&img, gocv.ColorBGRToHSV)
	// img.ConvertTo(&img, gocv.ColorBGRToHSV)
	img.ConvertTo( &img , 40 )

	black := color.RGBA{0, 0, 0, 0}
	white := color.RGBA{255, 255, 255, 0}

	width, height := img.Size()[1], img.Size()[0]
	cx, cy := width/2, height/2

	segments := [][]int{
		{0, 0, cx, cy},          // top left
		{cx, 0, width, cy},      // top right
		{0, cy, cx, height},     // bottom left
		{cx, cy, width, height}, // bottom right
	}

	axesX, axesY := int((float32(width)*0.75)/2), int((float32(height)*0.75)/2)
	ellipMask := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC1)
	defer ellipMask.Close()
	gocv.Ellipse(&ellipMask, image.Point{cx, cy}, image.Point{axesX, axesY}, 0, 0, 360, white, -1)

	segmentMask := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC1)
	defer segmentMask.Close()

	segmentHistogram := gocv.NewMat()
	defer segmentHistogram.Close()
	for _, segment := range segments {
		// reset mask
		gocv.Rectangle(&segmentMask, image.Rect(0, 0, width, height), black, -1)

		// calculate intersection of current segment and elliptic mask
		gocv.Rectangle(&segmentMask, image.Rect(segment[0], segment[1], segment[2], segment[3]), white, -1)
		gocv.Subtract(segmentMask, ellipMask, &segmentMask)

		// ShowMask(segmentMask)

		segmentFeatures, _ := featuresInSegment(img, segmentMask, segmentHistogram)
		features = append(features, segmentFeatures...)
	}

	// ShowMask(ellipMask)
	// why

	ellipFeatures, _ := featuresInSegment(img, ellipMask, segmentHistogram)
	features = append(features, ellipFeatures...)

	return
}

func GetFeatureVector( image_bytes []byte ) ( features []float64 ) {
	img , _ := gocv.IMDecode( image_bytes , gocv.IMReadColor )
	if img.Empty() { panic( "image is empty" ) }
	defer img.Close()

	// convert img to hsv color space
	// img.ConvertTo(&img, gocv.ColorBGRToHSV)
	// img.ConvertTo(&img, gocv.ColorBGRToHSV)
	img.ConvertTo( &img , 40 )

	black := color.RGBA{0, 0, 0, 0}
	white := color.RGBA{255, 255, 255, 0}

	width, height := img.Size()[1], img.Size()[0]
	cx, cy := width/2, height/2

	segments := [][]int{
		{0, 0, cx, cy},          // top left
		{cx, 0, width, cy},      // top right
		{0, cy, cx, height},     // bottom left
		{cx, cy, width, height}, // bottom right
	}

	axesX, axesY := int((float32(width)*0.75)/2), int((float32(height)*0.75)/2)
	ellipMask := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC1)
	defer ellipMask.Close()
	gocv.Ellipse(&ellipMask, image.Point{cx, cy}, image.Point{axesX, axesY}, 0, 0, 360, white, -1)

	segmentMask := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC1)
	defer segmentMask.Close()

	segmentHistogram := gocv.NewMat()
	defer segmentHistogram.Close()
	for _, segment := range segments {
		// reset mask
		gocv.Rectangle(&segmentMask, image.Rect(0, 0, width, height), black, -1)

		// calculate intersection of current segment and elliptic mask
		gocv.Rectangle(&segmentMask, image.Rect(segment[0], segment[1], segment[2], segment[3]), white, -1)
		gocv.Subtract(segmentMask, ellipMask, &segmentMask)

		// ShowMask(segmentMask)

		segmentFeatures, _ := featuresInSegment(img, segmentMask, segmentHistogram)
		features = append(features, segmentFeatures...)
	}

	// ShowMask(ellipMask)
	// why

	ellipFeatures, _ := featuresInSegment(img, ellipMask, segmentHistogram)
	features = append(features, ellipFeatures...)

	return
}



// https://stats.stackexchange.com/questions/470720/why-is-it-called-chi2-distance-kernel#471093
func chi2Distance(v1, v2 []float64) float64 {
	d := 0.
	min := int(math.Min(float64(len(v1)), float64(len(v2))))

	for i := 0; i < min; i++ {
		distance := math.Pow(v1[i]-v2[i], 2) / (v1[i] + v2[i] + 1e-10)
		d += distance
	}

	return d * 0.5
}

func CalculateDistance( reference_image_features []float64 , test_image_features []float64 ) ( result float64 ) {
	result = chi2Distance( reference_image_features , test_image_features )
	return
}