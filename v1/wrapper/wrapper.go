package wrapper

import (
	"strings"
	"strconv"
	"regexp"
	"bufio"
	"time"
	"math/rand"
	"fmt"
	"image"
	"sync"
	"sort"
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
	"image/png"
	// "bytes"
	utils "github.com/0187773933/ADBWrapper/v1/utils"
	image_similarity "github.com/0187773933/ADBWrapper/v1/image-similarity"
	try "github.com/manucorporat/try"
	// https://github.com/denismakogon/gocv-alpine
)

const IMAGE_SIMILARITY_THRESHOLD float64 = 1.5
const PAUSE_THRESHOLD = ( 500 * time.Millisecond )
const EXEC_TIMEOUT = ( 1500 * time.Millisecond )

type Wrapper struct {
	ADBPath string `json:"adb_path"`
	Serial string `json:"serial"`
	Connected bool `json:"connected"`
	Screen bool `json:"screen_on"`
}

func ConnectIP( adb_path string , host_ip string , host_port string ) ( wrapper Wrapper ) {
	wrapper.ADBPath = adb_path
	wrapper.Serial = host_ip + ":" + host_port
	connection_result := utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , adb_path , "connect" , host_ip )
	if strings.Contains( connection_result , "already connected" ) {
		wrapper.Connected = true
	} else if strings.Contains( connection_result , "failed to connect" ) {
		wrapper.Connected = false
	} else if len( strings.TrimSpace( connection_result ) ) == 0 {
		wrapper.Connected = false
	}
	// if force_screen_on == true { wrapper.ScreenOn() }
	return
}

func ConnectUSB( adb_path string , serial string ) ( wrapper Wrapper ) {
	wrapper.ADBPath = adb_path
	wrapper.Serial = serial
	wrapper.Connected = true
	// if force_screen_on == true { wrapper.ScreenOn() }
	return
}

func ( w *Wrapper ) Exec( arguments ...string ) ( result string ) {
	args := append( []string{ "-s" , w.Serial } , arguments... )
	result = utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , w.ADBPath , args... )
	return
}

func ( w *Wrapper ) Shell( arguments ...string ) ( result string ) {
	args := append( []string{ "-s" , w.Serial , "shell" } , arguments... )
	result = utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , w.ADBPath , args... )
	return
}

func ( w *Wrapper ) GetScreenState() ( result bool ) {
	// Display Power: state=OFF
	x := w.Shell( "dumpsys" , "power" )
	lines := strings.Split( x , "\n" )
	for _ , line := range lines {
		if strings.Contains( line , "Display Power: state=" ) {
			parts := strings.Split( strings.TrimSpace( line ) , "state=" )
			if len( parts ) < 1 { break; }
			state := strings.TrimSpace( parts[ 1 ] )
			switch state {
				case "ON":
					result = true
					break;
				case "OFF":
					result = false
					break;
			}
			break
		}
	}
	return
}

func ( w *Wrapper )  ScreenOn() {
	w.Screen = w.GetScreenState()
	if w.Screen == true { return; }
	fmt.Println( w.PressKey( 26 ) )
}

func ( w *Wrapper ) ScreenOff() {
	w.Screen = w.GetScreenState()
	if w.Screen == false { return; }
	fmt.Println( w.PressKey( 26 ) )
}

func ( w *Wrapper ) ForceScreenOn() ( screen_was_off bool ) {
	w.ScreenOn()
	screen_was_off = true
	if w.Screen == false {
		time.Sleep( 500 * time.Millisecond )
		w.ScreenOn()
		if w.Screen == false {
			time.Sleep( 500 * time.Millisecond )
			w.ScreenOn()
		}
	} else { screen_was_off = false }
	// fmt.Println( "Connected , Screen On ===" , w.Screen , " , Screen Was Off ===" , screen_was_off )
	return
}

// adb shell "sqlite3 /data/data/com.android.providers.settings/databases/settings.db  \"update system set value='-1' where name='screen_off_timeout'\";"
func ( w *Wrapper ) DisableScreenTimeout() {
	// w.Shell( "settings" , "put" , "system" , "background_power_saving_enable" , "0" )
	// w.Shell( "settings" , "put" , "system" , "screen_off_timeout" , "2147483647" )
	w.Shell( "svc" , "power" , "stayon" , "true" )
}

func ( w *Wrapper ) EnableScreenTimeout() {
	w.Shell( "svc" , "power" , "stayon" , "false" )
}

func ( w *Wrapper ) SetVolume( level int ) {
	w.Shell( "media" , "volume" , "--stream" , "3" , "--set" , strconv.Itoa( level ) )
}

func (w *Wrapper) SetVolumePercent( percent int ) {
	output := w.Shell( "media", "volume", "--stream", "3", "--get" )
	re := regexp.MustCompile( `volume is (\d+) in range \[(\d+)\.\.(\d+)\]` )
	matches := re.FindStringSubmatch(output)
	if len(matches) != 4 {
		fmt.Println("Failed to parse volume information")
		return
	}

	// Parse the current volume and range
	currentVolume, _ := strconv.Atoi(matches[1])
	minVolume, _ := strconv.Atoi(matches[2])
	maxVolume, _ := strconv.Atoi(matches[3])

	fmt.Printf("Current volume: %d, Min volume: %d, Max volume: %d\n", currentVolume, minVolume, maxVolume)

	// Calculate the desired volume based on the percentage
	desiredVolume := minVolume + (maxVolume-minVolume)*percent/100

	// Set the volume
	fmt.Printf("Setting volume to %d\n", desiredVolume)
	w.Shell("media", "volume", "--stream", "3", "--set", strconv.Itoa(desiredVolume))
}

func ( w *Wrapper )  GetTopWindowInfo() ( lines []string ) {
	result := w.Shell( "dumpsys" , "window" , "windows" )
	fmt.Println( result )
	// command := exec.Command( "bash" , "-c" , "adb shell dumpsys window windows" )
	// var outb , errb bytes.Buffer
	// command.Stdout = &outb
	// command.Stderr = &errb
	// command.Start()
	// time.AfterFunc( ( EXEC_TIMEOUT * time.Millisecond ) , func() {
	// 	command.Process.Signal( syscall.SIGTERM )
	// })
	// command.Wait()
	// result := outb.String()
	// start := strings.Split( result , "Window #1" )[ 1 ]
	// middle := strings.Split( start , "Window #2" )[ 0 ]
	// non_empty_lines := strings.Replace( middle , "\n\n" , "\n" , -1 )
	// lines = strings.Split( non_empty_lines , "\n" )
	return
}

// type EventDevice struct {
// 	DevicePath string
// 	Bus        string
// 	Vendor     string
// 	Product    string
// 	Version    string
// 	Name       string
// 	Location   string
// 	ID         string
// 	Events     string
// 	Props      string
// }
// just run === adb shell getevent -il
// to find you device name and events and stuff
// http://ktnr74.blogspot.com/2013/06/emulating-touchscreen-interaction-with.html
func ( w *Wrapper ) GetEventDevices() ( lines []string ) {

	result := w.Shell( "getevent" , "-il" )
	fmt.Println( result )

	// TODO , finish moving this from utils to here for ( w *Wrapper ) context
	// command := exec.Command( "bash" , "-c" , "adb shell getevent -il" )
	// var outb , errb bytes.Buffer
	// command.Stdout = &outb
	// command.Stderr = &errb
	// command.Start()
	// time.AfterFunc( ( EXEC_TIMEOUT * time.Millisecond ) , func() {
	// 	command.Process.Signal( syscall.SIGTERM )
	// })
	// command.Wait()
	// result := outb.String()
	// non_empty_lines := strings.Replace( result , "\n\n" , "\n" , -1 )
	// lines = strings.Split( non_empty_lines , "\n" )

	// start := strings.Split( result , "Window #1" )[ 1 ]
	// middle := strings.Split( start , "Window #2" )[ 0 ]
	// non_empty_lines := strings.Replace( middle , "\n\n" , "\n" , -1 )
	// lines = strings.Split( non_empty_lines , "\n" )
	return
}


func ( w *Wrapper ) OpenURI( uri string ) ( result string ) {
	result = w.Shell( "am" , "start" , "-a" , "android.intent.action.VIEW" , "-d" , uri )
	return
}

// adb shell pm list packages
func ( w *Wrapper ) OpenAppName( app_name string ) ( result string ) {
	result = w.Shell( "monkey", "-p", app_name , "-c", "android.intent.category.LAUNCHER", "1" )
	// fmt.Println( result )
	return
}

func ( w *Wrapper ) CloseAppName( app_name string ) ( result string ) {
	result = w.Shell( "am", "force-stop", app_name )
	// fmt.Println( result )
	return
}

func ( w *Wrapper ) GetRunningApps() ( packages []string ) {
	result := w.Shell( "dumpsys", "activity" )
	lines := strings.Split( result , "\n" )
	package_map := make( map[ string ] bool )
	re := regexp.MustCompile( `A=(\S+:[\w\.]+)` )
	for _ , line := range lines {
		if strings.Contains( line , "TaskRecord{" ) == false { continue; }
		matches := re.FindStringSubmatch( line )
		if len( matches ) < 1 { continue; }
		package_name := strings.Split( matches[ 1 ] , ":" )[ 1 ]
		if package_name == "com.amazon.firelauncher" { continue }
		package_map[ package_name ] = true
	}
	for key := range package_map { packages = append( packages , key ); }
	return
}

func ( w *Wrapper ) StopAllApps() {
	open_apps := w.GetRunningApps()
	for _ , app := range open_apps {
		w.Shell( "am", "force-stop", app )
	}
	return
}

func ( w *Wrapper ) PressButtonSequence( buttons ...int ) ( result string ) {
	sequence_string := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(buttons)), " " ), "[]")
	result = w.Shell( "input" , "keyevent" , sequence_string )
	return
}

func ( w *Wrapper ) Tap( x int , y int ) ( result string ) {
	result = w.Shell( "input" , "tap" , strconv.Itoa( x ) , strconv.Itoa( y ) )
	return
}

func ( w *Wrapper ) PressKey( key_number int ) ( result string ) {
	result = w.Shell( "input" , "keyevent" , strconv.Itoa( key_number ) )
	return
}

func ( w *Wrapper ) PressKeyName( key_name string ) ( result string ) {
	result = w.Shell( "input" , "keyevent" , key_name )
	return
}


// https://ktnr74.blogspot.com/2013/06/emulating-touchscreen-interaction-with.html
func ( w *Wrapper ) Swipe( start_x int , start_y int , stop_x int , stop_y int ) ( result string ) {
	result = w.Shell( "input" , "swipe" , strconv.Itoa( start_x ) , strconv.Itoa( start_y ) , strconv.Itoa( stop_x ) , strconv.Itoa( stop_y ) , strconv.Itoa( 100 ) )
	return
}

func ( w *Wrapper ) Type( text string ) ( result string ) {
	text = strings.ReplaceAll( text , " " , "%s" )
	text = strings.ReplaceAll( text , "'" , "\\'" )
	text = strings.ReplaceAll( text , "\"" , "\\\"" )
	result = w.Shell( "input" , "text" , text )
	return
}

// https://github.com/imba28/image-similarity
// https://github.com/baitisj/android_screen_mirror
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/search.go#L36
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/index.go#L10
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/descriptor.go#L16
// https://github.com/hybridgroup/gocv/blob/e11806566cdf2482485cc90d92ed320fa920e91a/cmd/img-similarity/main.go#L123
func ( w *Wrapper ) Screenshot( save_path string , crop ...int ) ( result string ) {
	if save_path == "" {
		temp_dir := os.TempDir()
		save_path = filepath.Join( temp_dir , "adb_screenshot_4524124.png" )
	}
	// args := []string{"-s", w.Serial, "exec-out", "screencap", "-p"}
	// command := exec.Command(w.ADBPath, args...)
	// output, err := command.Output()
	// if err != nil {
	// 	panic(err)
	// }
	// err = ioutil.WriteFile(save_path, output, 0644)
	cmd_str := fmt.Sprintf( "%s -s %s exec-out screencap -p > %s", w.ADBPath, w.Serial, save_path)
	cmd := exec.Command( "bash", "-c", cmd_str )
	cmd.Run()

	// Crop if bounding-box is present
	if len( crop ) == 4 {
		// x1 , y1 , x2 , y2 := crop[ 0 ] , crop[ 1 ] , crop[ 2 ] , crop[ 3 ]
		x1 , y1 , width , height := crop[ 0 ] , crop[ 1 ] , crop[ 2 ] , crop[ 3 ]
		crop_file , crop_file_err := os.Open( save_path )
		if crop_file_err != nil { fmt.Println( crop_file_err ); return }
		defer crop_file.Close()
		crop_img_src , crop_img_src_err := png.Decode( crop_file )
		if crop_file_err != nil { fmt.Println( crop_img_src_err ); return }
		// crop_img := crop_img_src.(*image.NRGBA).SubImage(image.Rect(x1, y1, x2, y2)).(*image.NRGBA)
		crop_img := crop_img_src.(*image.NRGBA).SubImage(image.Rect(x1, y1, x1+width, y1+height)).(*image.NRGBA)
		crop_img_out_file , crop_img_out_file_err := os.Create( save_path )
		if crop_file_err != nil { fmt.Println( crop_img_out_file_err ); return }
		defer crop_img_out_file.Close()
		encode_err := png.Encode( crop_img_out_file , crop_img )
		if encode_err != nil { fmt.Println( encode_err ); return }
	}

	return
}

// func ( w *Wrapper ) Screenshot() ( result string ) {
// 	// Bad
// 	// result = w.Exec( "exec-out" , "screencap" , "-p > test.png" )
// 	// result = w.Shell( "screencap -p > test.png" )
// 	// Good
// 	// screenshot_bytes = w.Exec( "exec-out" , "stty raw; screencap -p" )
// 	screenshot_bytes := utils.ExecProcessWithTimeoutGetBytes( ( EXEC_TIMEOUT * time.Millisecond ) ,
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

func ( w *Wrapper ) get_current_screen_features( crop ...int ) ( features []float64 ) {
	try.This(func() {
		temp_dir := os.TempDir()
		temp_save_path := filepath.Join( temp_dir , "adb_screenshot_4524124.png" )
		utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , "bash" , "-c" ,
			fmt.Sprintf( "%s -s %s exec-out screencap -p > %s" , w.ADBPath , w.Serial , temp_save_path ) ,
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

		// Crop if bounding-box is present
		if len( crop ) == 4 {
			// x1 , y1 , x2 , y2 := crop[ 0 ] , crop[ 1 ] , crop[ 2 ] , crop[ 3 ]
			x1 , y1 , width , height := crop[ 0 ] , crop[ 1 ] , crop[ 2 ] , crop[ 3 ]
			crop_file , crop_file_err := os.Open( temp_save_path )
			if crop_file_err != nil { fmt.Println( crop_file_err ); return }
			defer crop_file.Close()
			crop_img_src , crop_img_src_err := png.Decode( crop_file )
			if crop_file_err != nil { fmt.Println( crop_img_src_err ); return }
			// crop_img := crop_img_src.(*image.NRGBA).SubImage(image.Rect(x1, y1, x2, y2)).(*image.NRGBA)
			crop_img := crop_img_src.(*image.NRGBA).SubImage(image.Rect(x1, y1, x1+width, y1+height)).(*image.NRGBA)
			crop_img_out_file , crop_img_out_file_err := os.Create( temp_save_path )
			if crop_file_err != nil { fmt.Println( crop_img_out_file_err ); return }
			defer crop_img_out_file.Close()
			encode_err := png.Encode( crop_img_out_file , crop_img )
			if encode_err != nil { fmt.Println( encode_err ); return }
		}

		features = image_similarity.GetFeatureVector( temp_save_path )
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	return
}
func ( w *Wrapper ) current_screen_similarity_to_reference_image( reference_image_path string , crop ...int ) ( distance float64 ) {
	try.This(func() {
		current_screen_features := w.get_current_screen_features( crop... )
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

func ( w *Wrapper ) IsSameScreen( reference_image_path string , crop ...int ) ( result bool ) {
	distance := w.current_screen_similarity_to_reference_image( reference_image_path , crop... )
	// fmt.Println( distance )
	if distance > IMAGE_SIMILARITY_THRESHOLD {
		result = false
	} else {
		result = true
	}
	return
}

func ( w *Wrapper ) IsSameScreenV2( reference_image_path string , crop ...int ) ( result bool , distance float64 ) {
	distance = w.current_screen_similarity_to_reference_image( reference_image_path , crop... )
	// fmt.Println( distance )
	if distance > IMAGE_SIMILARITY_THRESHOLD {
		result = false
	} else {
		result = true
	}
	return
}

func ( w *Wrapper ) WaitOnScreen( reference_image_path string , timeout time.Duration , crop ...int ) ( result bool ) {
	done := make(chan bool, 1)

	// Create a timer that will send a message on its channel after the timeout
	timer := time.NewTimer(timeout)

	// Create a ticker that will send a message on its channel every 500ms
	ticker := time.NewTicker( 500 * time.Millisecond )

	go func() {
		for {
			select {
			// When the ticker ticks, call adb.IsSameScreen
			case <-ticker.C:
				if w.IsSameScreen( reference_image_path , crop... ) {
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
func ( w *Wrapper ) ClosestScreen( reference_image_path_directory string , crop ...int ) ( result string ) {
	files , _ := os.ReadDir( reference_image_path_directory )

	current_screen_features := w.get_current_screen_features( crop... )

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
	EventNum int `json:"EventNum"`
	TypeDec  int `json:"TypeDec"`
	CodeDec  int `json:"CodeDec"`
	ValueDec int `json:"ValueDec"`
	Time time.Time `json:"Time"`
}
var (
	streams       [][]Event
	currentStream []Event
	lastEventTime = time.Now()
)
func handle_event_output( cmd *exec.Cmd , done chan bool ) {
	stdout, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner( stdout )
	defer stdout.Close()
	r , _ := regexp.Compile( `^(/dev/input/event\d+): ([0-9a-f]+) ([0-9a-f]+) ([0-9a-f]+)$` )
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
						fmt.Printf( "Stream %d:\n", len(streams))
						for _, event := range currentStream {
							fmt.Println(event)
						}
						currentStream = nil
					}
					lastEventTime = currentTime
				}

				currentStream = append(currentStream, Event{
					EventNum: eventNum,
					TypeDec:  int(typeDec) ,
					CodeDec:  int(codeDec) ,
					ValueDec: int(valueDec) ,
					Time: time.Now() ,
				})
			}
		}
	}
}

// go slow
func ( w *Wrapper ) SaveEvents( save_path string ) {
	args := []string{ "-s" , w.Serial , "shell" , "getevent" }
	fmt.Println( w.ADBPath , args )
	cmd := exec.Command( w.ADBPath , args... )
	done := make( chan bool )
	go handle_event_output( cmd , done )
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

func ( w *Wrapper ) PlaybackEvents( save_path string ) {

	data , err := ioutil.ReadFile( save_path )
	if err != nil {
		fmt.Println( "Error reading JSON file:" , err )
		return
	}
	var streams [][]Event
	err = json.Unmarshal( data , &streams )
	if err != nil {
		fmt.Println( "Error parsing JSON file:" , err )
		return
	}
	dropout_percentage := 0.4
	fmt.Println( "total of streams" , len( streams ) )
	for _ , stream := range streams {

		// Remove % of the events
		minus_percent := int( float64( len( stream ) - 2 ) * dropout_percentage ) // Exclude the first and last elements
		minus_percent_indexes := make( map[ int ]bool )
		rand.Seed( time.Now().UnixNano() )
		for i := 0; i < minus_percent; i++ {
			for {
				index := rand.Intn(len(stream)-2) + 1 // +1 to avoid the first element
				if !minus_percent_indexes[index] {
					minus_percent_indexes[index] = true
					break
				}
			}
		}

		var new_events []Event
		for i, event := range stream {
			if i == 0 || i == len(stream)-1 || !minus_percent_indexes[i] {
				new_events = append( new_events , event )
			}
		}

		// Sort By Time
		sort.Slice( new_events , func( i , j int ) bool {
			return new_events[i].Time.Before( new_events[j].Time )
		})
		fmt.Println( "total of events" , len( new_events ) )

		// Build Commands
		var commands []string
		for _ , event := range new_events {
			eventFile := "/dev/input/event" + strconv.Itoa( event.EventNum )
			command := fmt.Sprintf( `S="sendevent %s";$S %d %d %d;` , eventFile , event.TypeDec , event.CodeDec , event.ValueDec )
			commands = append( commands , command )
		}

		args := []string{ "-s" , w.Serial , "shell" }
		args = append( args , commands... )
		fmt.Println( len( commands ) )
		cmd := exec.Command( w.ADBPath , args... )
		cmd.Run()
	}
}


