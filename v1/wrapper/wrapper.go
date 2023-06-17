package wrapper

import (
	"strings"
	"strconv"
	"regexp"
	"bufio"
	"time"
	"fmt"
	"sync"
	"os"
	"context"
	"syscall"
	"os/exec"
	"os/signal"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	// "unsafe"
	// "image"
	_ "image/jpeg"
	_ "image/png"
	// "bytes"
	utils "github.com/0187773933/ADBWrapper/v1/utils"
	image_similarity "github.com/0187773933/ADBWrapper/v1/image-similarity"
	try "github.com/manucorporat/try"
	// https://github.com/denismakogon/gocv-alpine
)

const IMAGE_SIMILARITY_THRESHOLD float64 = 1.5
const PAUSE_THRESHOLD = ( 500 * time.Millisecond )

type Wrapper struct {
	ADBPath string `json:"adb_path"`
	HostIP string `json:"host_ip"`
	HostPort string `json:"host_port"`
	Connected bool `json:"connected"`
}

func ConnectIP( adb_path string , host_ip string , host_port string ) ( wrapper Wrapper ) {
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

func ConnectUSB( adb_path string , serial string ) ( wrapper Wrapper ) {
	wrapper.ADBPath = adb_path
	wrapper.Connected = true
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

// adb shell pm list packages
func ( w *Wrapper ) OpenAppName( app_name string ) ( result string ) {
	result = w.Exec( "shell", "monkey", "-p", app_name , "-c", "android.intent.category.LAUNCHER", "1" )
	// fmt.Println( result )
	return
}

func ( w *Wrapper ) PressButtonSequence( buttons ...int ) ( result string ) {
	sequence_string := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(buttons)), " " ), "[]")
	result = w.Exec( "shell" , "input" , "keyevent" , sequence_string )
	return
}

func ( w *Wrapper ) Tap( x int , y int ) ( result string ) {
	result = w.Exec( "shell" , "input" , "tap" , strconv.Itoa( x ) , strconv.Itoa( y ) )
	return
}

func ( w *Wrapper ) PressKey( key_number int ) ( result string ) {
	result = w.Exec( "shell" , "input" , "keyevent" , strconv.Itoa( key_number ) )
	return
}

// Doesn't Work
// https://ktnr74.blogspot.com/2013/06/emulating-touchscreen-interaction-with.html
func ( w *Wrapper ) Swipe( start_x int , start_y int , stop_x int , stop_y int ) ( result string ) {
	result = w.Exec( "shell" , "input" , "swipe" , strconv.Itoa( start_x ) , strconv.Itoa( start_y ) , strconv.Itoa( stop_x ) , strconv.Itoa( stop_y ) , strconv.Itoa( 100 ) )
	return
}

func (w *Wrapper) Type( text string ) ( result string ) {
	text = strings.ReplaceAll(text, " ", "%s")
	text = strings.ReplaceAll(text, "'", "\\'")
	text = strings.ReplaceAll(text, "\"", "\\\"")
	result = w.Exec("shell", "input", "text", text)
	return
}

func ( w *Wrapper ) GetTopWindowInfo() ( results []string ) {
	results = utils.GetTopWindowInfo()
	// for i , line := range results {

	// }
	return
}

func ( w *Wrapper ) GetEventDevices() ( results []string ) {
	results = utils.GetEventDevices()
	// for i , line := range results {

	// }
	return
}

// https://github.com/imba28/image-similarity
// https://github.com/baitisj/android_screen_mirror
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/search.go#L36
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/index.go#L10
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/descriptor.go#L16
// https://github.com/hybridgroup/gocv/blob/e11806566cdf2482485cc90d92ed320fa920e91a/cmd/img-similarity/main.go#L123
func ( w *Wrapper ) Screenshot( save_path string ) ( result string ) {
	if save_path == "" {
		temp_dir := os.TempDir()
		save_path = filepath.Join( temp_dir , "adb_screenshot_4524124.png" )
	}
	result = utils.ExecProcessWithTimeout( ( 1500 * time.Millisecond ) , "bash" , "-c" ,
		fmt.Sprintf( "adb exec-out screencap -p > %s" , save_path ) ,
	)
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

func ( w *Wrapper ) get_current_screen_features() ( features []float64 ) {
	try.This(func() {
		temp_dir := os.TempDir()
		temp_save_path := filepath.Join( temp_dir , "adb_screenshot_4524124.png" )
		utils.ExecProcessWithTimeout( ( 1500 * time.Millisecond ) , "bash" , "-c" ,
			fmt.Sprintf( "adb exec-out screencap -p > %s" , temp_save_path ) ,
		)

		// Option 1 - Just Wait the full 1.5 seconds
		// time.Sleep( 1501 * time.Millisecond )

		// Option 2 - Wait until the file exists
		// Wait for the screenshot to be created
		for {
			_ , err := os.Stat( temp_save_path )
			if err == nil { break }
			time.Sleep( 10 * time.Millisecond ) // sleep for a bit; don't busy wait
		}
		// Wait for the screenshot file size to stabilize
		previousSize := int64( -1 )
		for {
			fileInfo , err := os.Stat( temp_save_path )
			if err != nil { return }
			size := fileInfo.Size()
			if size == previousSize { break } // File size has not changed; assuming it's done writing
			previousSize = size
			time.Sleep( 10 * time.Millisecond )
		}

		features = image_similarity.GetFeatureVector( temp_save_path )
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	return
}

func ( w *Wrapper ) current_screen_similarity_to_reference_image( reference_image_path string ) ( distance float64 ) {
	try.This(func() {
		current_screen_features := w.get_current_screen_features()
		reference_image_features := image_similarity.GetFeatureVector( reference_image_path )
		distance = image_similarity.CalculateDistance( current_screen_features , reference_image_features )
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	return
}

func ( w *Wrapper ) similarity_to_feature_list( features []float64 , reference_image_path string ) ( distance float64 ) {
	try.This(func() {
		current_screen_features := w.get_current_screen_features()
		reference_image_features := image_similarity.GetFeatureVector( reference_image_path )
		distance = image_similarity.CalculateDistance( current_screen_features , reference_image_features )
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	return
}

func ( w *Wrapper ) IsSameScreen( reference_image_path string ) ( result bool ) {
	distance := w.current_screen_similarity_to_reference_image( reference_image_path )
	// fmt.Println( distance )
	if distance > IMAGE_SIMILARITY_THRESHOLD {
		result = false
	} else {
		result = true
	}
	return
}

func ( w *Wrapper ) WaitOnScreen( reference_image_path string , timeout time.Duration ) ( result bool ) {
	done := make(chan bool, 1)

	// Create a timer that will send a message on its channel after the timeout
	timer := time.NewTimer(timeout)

	// Create a ticker that will send a message on its channel every 500ms
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			select {
			// When the ticker ticks, call adb.IsSameScreen
			case <-ticker.C:
				if w.IsSameScreen(reference_image_path) {
					done <- true
					return
				}
			// When the timer finishes, stop checking
			case <-timer.C:
				done <- false
				return
			}
		}
	}()

	// Wait for either the function to finish or the timer to expire
	result = <-done

	// Always stop the ticker and timer when you're done to free resources
	ticker.Stop()
	timer.Stop()

	return
}


type ScreenHit struct {
	Path     string
	Distance float64
}
func ( w *Wrapper ) ClosestScreen( reference_image_path_directory string ) ( result string ) {
	files , _ := os.ReadDir( reference_image_path_directory )

	current_screen_features := w.get_current_screen_features()

	// Prepare the WaitGroup , semaphore , and context
	total_concurrent := 5
	var wg sync.WaitGroup
	semaphore := make( chan struct{} , total_concurrent )
    ctx, cancel := context.WithCancel( context.Background() )
    defer cancel() // make sure all paths cancel the context to release resources

	results := make( chan ScreenHit , len( files ) )
    for _ , f := range files {
        wg.Add( 1 )
        go func(f os.DirEntry) {
            defer wg.Done()
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            select {
	            case <-ctx.Done():
	                return // returning early if context was cancelled
	            default:
	                imagePath := filepath.Join( reference_image_path_directory , f.Name() )
	                distance := w.similarity_to_feature_list( current_screen_features , imagePath )
	                results <- ScreenHit{ imagePath , distance}
	                if distance < IMAGE_SIMILARITY_THRESHOLD {
	                    cancel() // this will cancel all other goroutines once threshold is met
	                }
	            }
        }( f )
    }
	go func() {
		// Wait for all goroutines to finish and then close the results channel
		wg.Wait()
		close( results )
	}()

	// Find the image with the smallest distance
	minDistance := float64( 1<<63 - 1 ) // set to maximum possible float64
	for x := range results {
		if x.Distance < minDistance {
			minDistance = x.Distance
			result = x.Path
		}
	}
	return
}

type Event struct {
	EventNum int
	TypeDec  int
	CodeDec  int
	ValueDec int
}
var (
	streams       [][]Event
	currentStream []Event
	lastEventTime = time.Now()
)
func handle_event_output(cmd *exec.Cmd, done chan bool) {
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	defer stdout.Close()
	r, _ := regexp.Compile(`^(/dev/input/event\d+): ([0-9a-f]+) ([0-9a-f]+) ([0-9a-f]+)$`)
	for scanner.Scan() {
		select {
		case <-done:
			return
		default:
			line := scanner.Text()
			match := r.FindStringSubmatch(line)
			if match != nil {
				eventFile, typeHex, codeHex, valueHex := match[1], match[2], match[3], match[4]
				typeDec, _ := strconv.ParseInt(typeHex, 16, 32)
				codeDec, _ := strconv.ParseInt(codeHex, 16, 32)
				valueDec, _ := strconv.ParseInt(valueHex, 16, 32)
				eventNum, _ := strconv.Atoi(strings.TrimPrefix(eventFile, "/dev/input/event"))

				currentTime := time.Now()
				if currentTime.Sub(lastEventTime) > PAUSE_THRESHOLD {
					if len(currentStream) > 0 {
						streams = append(streams, currentStream)
						fmt.Printf("Stream %d:\n", len(streams))
						for _, event := range currentStream {
							fmt.Println(event)
						}
						currentStream = nil
					}
					lastEventTime = currentTime
				}

				currentStream = append(currentStream, Event{
					EventNum: eventNum,
					TypeDec:  int(typeDec),
					CodeDec:  int(codeDec),
					ValueDec: int(valueDec),
				})
			}
		}
	}
}
func ( w *Wrapper ) SaveEvents( save_path string ) {

	cmd := exec.Command( "adb" , "shell" , "getevent" )
	done := make( chan bool )
	go handle_event_output( cmd, done )
	sigint := make( chan os.Signal , 1 )
	signal.Notify( sigint , os.Interrupt , syscall.SIGTERM )

	go func() {
		<-sigint
		done <- true
	}()

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	<-done

	if len( currentStream ) > 0 {
		streams = append( streams , currentStream )
		fmt.Printf( "\nStream %d:\n" , len( streams ) )
		for _ , event := range currentStream {
			fmt.Println( event )
		}
	}

	utils.WriteJSON( save_path , streams )
}


func adb_shell( commands []string ) {
	cmd := exec.Command("adb", "shell")
	cmd.Args = append(cmd.Args, commands...)
	err := cmd.Run()
	if err != nil {
		fmt.Println("adbShell error:", err)
	}
}

func playback_stream(stream []Event) {
	var commands []string
	for _, event := range stream {
		eventFile := "/dev/input/event" + strconv.Itoa(event.EventNum)
		command := fmt.Sprintf(`S="sendevent %s";$S %d %d %d;`, eventFile, event.TypeDec, event.CodeDec, event.ValueDec)
		commands = append(commands, command)
	}
	adb_shell(commands)
}

func ( w *Wrapper ) PlaybackEvents( save_path string ) {

	// Read the streams from the JSON file
	data, err := ioutil.ReadFile( save_path )
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var streams [][]Event
	err = json.Unmarshal(data, &streams)
	if err != nil {
		fmt.Println("Error parsing JSON file:", err)
		return
	}

	for _, stream := range streams {
		playback_stream(stream)
	}
}


