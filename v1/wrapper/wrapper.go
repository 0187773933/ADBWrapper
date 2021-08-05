package wrapper

import (
	"strings"
	"time"
	"fmt"
	"os"
	// "unsafe"
	// "image"
	"image/png"
	// "bytes"
	utils "github.com/0187773933/ADBWrapper/v1/utils"
	image_similarity "github.com/0187773933/ADBWrapper/v1/image-similarity"
	// https://github.com/denismakogon/gocv-alpine
)

type Wrapper struct {
	ADBPath string `json:"adb_path"`
	HostIP string `json:"host_ip"`
	HostPort string `json:"host_port"`
	Connected bool `json:"connected"`
}

func Connect( adb_path string , host_ip string , host_port string ) ( wrapper Wrapper ) {
	wrapper.ADBPath = adb_path
	wrapper.HostIP = host_ip
	wrapper.HostPort = host_port
	connection_result := utils.ExecProcessWithTimeout( ( 1500 * time.Millisecond ) , adb_path , "connect" , host_ip )
	if strings.Contains( connection_result , "already connected" ) {
		wrapper.Connected = true
	} else if strings.Contains( connection_result , "failed to connect" ) {
		wrapper.Connected = false
	} else if len( strings.TrimSpace( connection_result ) ) == 0 {
		wrapper.Connected = false
	}
	return
}

func ( w *Wrapper ) Exec( arguments ...string ) ( result string ) {
	result = utils.ExecProcessWithTimeout( ( 1500 * time.Millisecond ) , w.ADBPath , arguments... )
	return
}

func ( w *Wrapper ) OpenURI( uri string ) ( result string ) {
	result = w.Exec( "shell" , "am" , "start" , "-a" , "android.intent.action.VIEW" , "-d" , uri )
	return
}

func ( w *Wrapper ) PressButtonSequence( buttons ...int ) ( result string ) {
	sequence_string := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(buttons)), " " ), "[]")
	result = w.Exec( "shell" , "input" , "keyevent" , sequence_string )
	return
}

func ( w *Wrapper ) Tap( x int , y int ) ( result string ) {
	result = w.Exec( "shell" , "input" , "tap" , string( x ) , string( y ) )
	return
}

// https://github.com/imba28/image-similarity
// https://github.com/baitisj/android_screen_mirror
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/search.go#L36
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/index.go#L10
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/descriptor.go#L16
// https://github.com/hybridgroup/gocv/blob/e11806566cdf2482485cc90d92ed320fa920e91a/cmd/img-similarity/main.go#L123
func ( w *Wrapper ) Screenshot() ( result string ) {
	result = utils.ExecProcessWithTimeout( ( 1500 * time.Millisecond ) , "bash" , "-c" ,
		"adb exec-out screencap -p > /tmp/adb_screenshot_4524124.png" ,
	)
	time.Sleep( 1500 * time.Millisecond )
	screenshot_image , _ := os.Open( "/tmp/adb_screenshot_4524124.png" )
	defer screenshot_image.Close()
	// screenshot_image_data , _ , _ := image.Decode( screenshot_image )
	// fmt.Println( screenshot_image_data )
	// screenshot_image.Seek( 0 , 0 )
	screenshot_png_data , _ := png.Decode( screenshot_image )
	fmt.Println( screenshot_png_data )
	return
}

// func ( w *Wrapper ) Screenshot() ( result string ) {
// 	// Bad
// 	// result = w.Exec( "exec-out" , "screencap" , "-p > test.png" )
// 	// result = w.Exec( "shell" , "screencap -p > test.png" )
// 	// Good
// 	// screenshot_bytes = w.Exec( "exec-out" , "stty raw; screencap -p" )
// 	screenshot_bytes := utils.ExecProcessWithTimeoutGetBytes( ( 1500 * time.Millisecond ) ,
// 		// w.ADBPath , "exec-out" , "stty raw; screencap -p" ,
// 		w.ADBPath , "exec-out" , "screencap -p" ,
// 	)
// 	// fmt.Println( screenshot_bytes )
// 	img , _ , err := image.Decode( bytes.NewReader( screenshot_bytes ) )
// 	fmt.Println( err )
// 	fmt.Println( img )
// 	// img , _ , _ := image.Decode( bytes.NewReader( screenshot_bytes ) )
// 	// out , _ := os.Create( "screenshot.png" )
// 	// defer out.Close()
// 	// png.Encode( out , img )
// 	// var opts jpeg.Options
// 	// opts.Quality = 1

// 	// err = jpeg.Encode(out, img, &opts)
// 	// //jpeg.Encode(out, img, nil)
// 	// if err != nil {
// 	//     log.Println(err)
// 	// }
// 	return
// }


func ( w *Wrapper ) CurrentScreenSimilarityToReferenceImage( reference_image_path string ) ( distance float64 ) {
	utils.ExecProcessWithTimeout( ( 1500 * time.Millisecond ) , "bash" , "-c" ,
		"adb exec-out screencap -p > /tmp/adb_screenshot_4524124.png" ,
	)
	time.Sleep( 1500 * time.Millisecond )
	current_screen_features := image_similarity.GetFeatureVector( "/tmp/adb_screenshot_4524124.png" )
	reference_image_features := image_similarity.GetFeatureVector( reference_image_path )
	distance = image_similarity.CalculateDistances( current_screen_features , reference_image_features )
	return
}