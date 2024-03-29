package wrapper

import (
	"strings"
	"strconv"
	"regexp"
	"bufio"
	"bytes"
	"time"
	"math"
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
	// fsnotify "github.com/fsnotify/fsnotify"
	// "unsafe"
	// "image"
	_ "image/jpeg"
	"image/png"
	color "image/color"
	// "bytes"
	utils "github.com/0187773933/ADBWrapper/v1/utils"
	image_similarity "github.com/0187773933/ADBWrapper/v1/image-similarity"
	try "github.com/manucorporat/try"
	// https://github.com/denismakogon/gocv-alpine
)



const IMAGE_SIMILARITY_THRESHOLD float64 = 1.5
const PAUSE_THRESHOLD = ( 500 * time.Millisecond )
const EXEC_TIMEOUT = ( 1500 * time.Millisecond )

func open_with_preview(imagePath string) {
	cmd := exec.Command("open", "-a", "Preview", imagePath)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to open image:", err)
	}
}

type Wrapper struct {
	ADBPath string `json:"adb_path"`
	Serial string `json:"serial"`
	Connected bool `json:"connected"`
	Screen bool `json:"screen_on"`
}

func ConnectIP( adb_path string , host_ip string , host_port string ) ( wrapper Wrapper ) {
	wrapper.ADBPath = adb_path
	wrapper.Serial = host_ip + ":" + host_port
	connection_result := utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , adb_path , "connect" , wrapper.Serial )
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

func ( w *Wrapper ) RestartServer() {
	utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , w.ADBPath , "kill-server" )
	time.Sleep( 100 * time.Millisecond )
	utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , w.ADBPath , "start-server" )
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

func ( w *Wrapper ) GetCPUArchitecture() ( result string ) {
	x := w.Shell( "getprop" , "ro.product.cpu.abi" )
	lines := strings.Split( x , "\n" )
	result = lines[ 0 ]
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

type Status struct {
	DisplayOn bool `json:"display_on"`
	Volume int `json:"volume"`
	Activity string `json:"activity"`
	MediaSession MediaSession `json:"media_session"`
}

func ( w *Wrapper ) GetStatus() ( result Status ) {
	result.DisplayOn = w.GetScreenState()
	result.Volume = w.GetVolume()
	result.Activity = w.GetActivity()
	result.MediaSession = w.GetMediaSessionInfo()
	return
}

type MediaSession struct {
	Type string `json:"type"`
	Activity string `json:"activity"`
	Package string `json:"package"`
	State string `json:"state"`
	Position string `json:"position"`
	BufferedPosition string `json:"buffered_position"`
	UpdatedTime string `json:"updated_time"`
	Speed string `json:"speed"`
	Description string `json:"description"`
}
func ( w *Wrapper ) GetMediaSessionInfo() ( result MediaSession ) {
	media_session_dump := w.Shell( "dumpsys" , "media_session" )
	media_session_dump_lines := strings.Split( media_session_dump , "\n" )
	for line_index , line := range media_session_dump_lines {
		if strings.Contains( line , "active=true" ) {
			session_type_line := media_session_dump_lines[ ( line_index - 5 ) ]
			if strings.Contains( session_type_line , "bluetooth" ) { continue; }
			session := strings.Split( session_type_line , " " )
			session = utils.RemoveEmpties( session )
			result.Type = session[ 0 ]
			result.Activity = session[ 1 ]
			state_line := media_session_dump_lines[ ( line_index + 4 ) ]
			state_key_values := strings.Split( state_line , "," )
			state_key_values = utils.RemoveEmpties( state_key_values )
			state_num := strings.Split( state_key_values[ 0 ] , "state=" )[ 2 ]
			switch state_num {
				case "0":
					result.State = "none"
				case "1":
					result.State = "stopped"
				case "2":
					result.State = "paused"
				case "3":
					result.State = "playing"
			}
			result.Position = strings.Split( state_key_values[ 1 ] , "position=" )[ 1 ]
			result.BufferedPosition = strings.Split( state_key_values[ 2 ] , "buffered position=" )[ 1 ]
			result.Speed = strings.Split( state_key_values[ 3 ] , "speed=" )[ 1 ]
			result.UpdatedTime = strings.Split( state_key_values[ 4 ] , "updated=" )[ 1 ]
			description_line := media_session_dump_lines[ ( line_index + 7 ) ]
			// description_line_items := strings.Split( description_line , "=" )
			result.Description = strings.Split( description_line , "description=" )[ 1 ]
		}
	}
	return
}

func ( w *Wrapper ) ScreenOn() {
	w.Screen = w.GetScreenState()
	if w.Screen == true { return; }
	fmt.Println( w.PressKey( 26 ) )
}

func ( w *Wrapper ) ScreenOff() {
	w.Screen = w.GetScreenState()
	if w.Screen == false { return; }
	fmt.Println( w.PressKey( 26 ) )
	w.Shell( "am" , "broadcast" , "-a" , "android.intent.action.SCREEN_OFF" )
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

func ( w *Wrapper ) GetVolume() ( result int ) {
	output := w.Shell( "media", "volume", "--stream", "3", "--get" )
	re := regexp.MustCompile( `volume is (\d+) in range \[(\d+)\.\.(\d+)\]` )
	matches := re.FindStringSubmatch(output)
	if len( matches ) != 4 {
		// fmt.Println( "Failed to parse volume information" )
		return
	}
	result, _ = strconv.Atoi(matches[1])
	// fmt.Printf("Current volume: %d\n")
	return
}


func ( w *Wrapper ) SetVolumePercent( percent int ) {
	output := w.Shell( "media", "volume", "--stream", "3", "--get" )
	re := regexp.MustCompile( `volume is (\d+) in range \[(\d+)\.\.(\d+)\]` )
	matches := re.FindStringSubmatch(output)
	if len( matches ) != 4 {
		// fmt.Println( "Failed to parse volume information" )
		return
	}

	// Parse the current volume and range
	currentVolume , _ := strconv.Atoi( matches[ 1 ] )
	minVolume , _ := strconv.Atoi( matches[ 2 ] )
	maxVolume , _ := strconv.Atoi( matches[ 3 ] )

	fmt.Printf( "Current volume: %d, Min volume: %d, Max volume: %d\n" , currentVolume , minVolume , maxVolume )

	// Calculate the desired volume based on the percentage
	desiredVolume := minVolume + (maxVolume-minVolume)*percent/100

	// Set the volume
	fmt.Printf("Setting volume to %d\n", desiredVolume)
	w.Shell("media", "volume", "--stream", "3", "--set", strconv.Itoa(desiredVolume))
}

type Window struct {
	Number int `json:"number"`
	Package string `json:"package"`
	Activity string `json:"activity"`
	IsOnScreen bool `json:"is_on_screen"`
	IsVisible bool `json:"is_visible"`
}
func ( w *Wrapper ) GetWindowStack() ( windows []Window ) {
	result := w.Shell( "dumpsys" , "window" , "windows" )
	var current_window *Window
	for _ , line := range strings.Split( result , "\n" ) {
		if current_window == nil {
			if strings.Contains( line , "Window #" ) {
				current_window = &Window{}
				win_num_parts_one := strings.Split( line , "Window #" )
				win_num_parts := strings.Split( win_num_parts_one[ 1 ] , " " )
				current_window.Number , _ = strconv.Atoi( win_num_parts[ 0 ] )
				parts := strings.Fields( line )
				last_part := parts[ ( len( parts ) - 1 ) ]
				// current_window.Activity = strings.Split( last_part , "}:" )[ 0 ]
				pa := strings.Split( last_part , "}:" )[ 0 ]
				pa_parts := strings.Split( pa , "/" )
				switch len( pa_parts ) {
					case 1:
						current_window.Package = pa_parts[ 0 ]
						current_window.Activity = ""
						break;
					case 2:
						current_window.Package = pa_parts[ 0 ]
						current_window.Activity = pa_parts[ 1 ]
						break;
				}
				continue
			}
		} else {
			if strings.Contains( line , "isOnScreen=" ) {
				value := strings.Split( line , "isOnScreen=" )[ 1 ]
				current_window.IsOnScreen = ( value == "true" )
			} else if strings.Contains( line , "isVisible=" ) {
				value := strings.Split( line , "isVisible=" )[ 1 ]
				current_window.IsVisible = ( value == "true" )
				windows = append( windows , *current_window )
				current_window = nil
			}
		}
	}
	sort.Slice( windows , func( i , j int ) bool {
		return windows[i].IsOnScreen && !windows[j].IsOnScreen
	})
	sort.Slice( windows , func( i , j int ) bool {
		return windows[i].IsVisible && !windows[j].IsVisible
	})
	obscuring_window_parts := strings.Split( result , "mObscuringWindow=" )
	if len( obscuring_window_parts ) > 1 {
		if strings.HasPrefix( obscuring_window_parts[ 1 ] , "Window{" )  {
			obscuring_window_line := strings.Split( obscuring_window_parts[ 1 ] , "\n" )[ 0 ]
			// fmt.Println( "we have a hidden window to deal with" )
			obscuring_window_line_parts := strings.Split( obscuring_window_line , " " )
			last_part := obscuring_window_line_parts[ ( len( obscuring_window_line_parts ) - 1 ) ]
			new_top_activity := strings.Split( last_part , "}" )[ 0 ]
			pa_parts := strings.Split( new_top_activity , "/" )
			x_package := ""
			x_activity := ""
			switch len( pa_parts ) {
				case 1:
					x_package = pa_parts[ 0 ]
					break;
				case 2:
					x_package = pa_parts[ 0 ]
					x_activity = pa_parts[ 1 ]
					break;
			}
			new_top_window := Window{
				Number: 0 ,
				Package:  x_package ,
				Activity: x_activity ,
				IsOnScreen: true ,
				IsVisible: true ,
			}
			new_window_stack := []Window{ new_top_window , windows[ 0 ] }
			// new_index := 1
			for _ , window := range windows[ 1 : ] {
				if window.Activity == x_activity { continue }
				// window.Number = new_index
				// new_index += 1
				new_window_stack = append( new_window_stack , window )
			}
			// windows = append( []Window{ new_top_window } , windows... )
			windows = new_window_stack
		}
	}
	return
}

func ( w *Wrapper ) GetTopWindow() ( result Window ) {
	windows := w.GetWindowStack()
	if len( windows ) > 0 {
		result = windows[ 0 ]
	}
	return
}

func ( w *Wrapper ) GetActivity() ( result string ) {
	windows := w.GetWindowStack()
	if len( windows ) > 0 {
		result = windows[ 0 ].Activity
	}
	return
}

func ( w *Wrapper ) GetPackage() ( result string ) {
	windows := w.GetWindowStack()
	if len( windows ) > 0 {
		result = windows[ 0 ].Package
	}
	return
}

// func ( w *Wrapper ) GetCurrentPackage() ( result string ) {
// 	// result = w.Shell( "dumpsys" , "window" , "windows" )
// 	// for _ , line := range strings.Split( result , "\n" ) {
// 	// 	if strings.Contains( line , "mCurrentFocus" ) {
// 	// 		parts := strings.Fields( line )
// 	// 		last_part := parts[ ( len( parts ) - 1 ) ]
// 	// 		result = strings.Split( last_part , "}" )[ 0 ]
// 	// 		result = strings.Split( last_part , "/" )[ 0 ]
// 	// 		break
// 	// 	}
// 	// }
// 	window := w.GetTopWindow()
// 	result = window.Activity
// 	return
// }

func ( w *Wrapper ) GetPlaybackPositionForce() ( package_name string , position int ) {
	// package_name = w.GetCurrentPackage()
	package_name = w.GetTopWindow().Package

	// the only way to get media_session to update is after user interaction ( pause , resume )
	// pausing is fine on hulu movies/tv , but *can* break on twitch livestreams
	// TODO = check if livestream needs restarted
	w.Shell( "input" , "keyevent" , "KEYCODE_MEDIA_PLAY_PAUSE" , "KEYCODE_MEDIA_PLAY_PAUSE" )
	result := w.Shell( "dumpsys" , "media_session" )
	lines := strings.Split( result , "\n" )
	for line_index , line := range lines {
		if strings.Contains( line , "active=true" ) {
			session_type_line := lines[ ( line_index - 5 ) ]
			if strings.Contains( session_type_line , package_name ) {
				playback_line := lines[ ( line_index + 4 ) ]
				position_str := strings.Split( playback_line , "position=" )[ 1 ]
				position_str = strings.Split( position_str , "," )[ 0 ]
				position , _ = strconv.Atoi( position_str )
				return
			}
		}
	}
	return
}

func ( w *Wrapper ) GetPlaybackPosition() ( package_name string , position int ) {
	// package_name = w.GetCurrentPackage()
	package_name = w.GetTopWindow().Package
	result := w.Shell( "dumpsys" , "media_session" )
	lines := strings.Split( result , "\n" )
	for line_index , line := range lines {
		if strings.Contains( line , "active=true" ) {
			session_type_line := lines[ ( line_index - 5 ) ]
			if strings.Contains( session_type_line , package_name ) {
				playback_line := lines[ ( line_index + 4 ) ]
				position_str := strings.Split( playback_line , "position=" )[ 1 ]
				position_str = strings.Split( position_str , "," )[ 0 ]
				position , _ = strconv.Atoi( position_str )
				return
			}
		}
	}
	return
}

// state=PlaybackState {state=1, position=0, buffered position=0, speed=1.0, updated=2195684109, actions=1049468, custom actions=[], active item id=-1, error=null}
type PlaybackResult struct {
	PackageStr string `json:"package_str"`
	Type string `json:"type"`
	State string `json:"state"`
	Position int `json:"position"`
	Updated int `json:"updated"`
}
func parse_playback_line( line string ) ( result PlaybackResult ) {
	position_str_parts := strings.Split( line , "position=" )
	if len( position_str_parts ) > 1 {
		position_str_parts := strings.Split( position_str_parts[ 1 ] , "," )
		if len( position_str_parts ) > 1 {
			result.Position , _ = strconv.Atoi( position_str_parts[ 0 ] )
		}
	}
	updated_str_parts := strings.Split( line , "updated=" )
	if len( updated_str_parts ) > 1 {
		updated_str_parts := strings.Split( updated_str_parts[ 1 ] , "," )
		if len( position_str_parts ) > 1 {
			result.Updated , _ = strconv.Atoi( updated_str_parts[ 0 ] )
		}
	}
	state_str_parts := strings.Split( line , "{state=" )
	if len( state_str_parts ) > 1 {
		state_str_parts = strings.Split( state_str_parts[ 1 ] , "," )
		if len( state_str_parts ) > 1 {
			switch state_str_parts[ 0 ] {
				case "0":
					result.State = "none"
				case "1":
					result.State = "stopped"
				case "2":
					result.State = "paused"
				case "3":
					result.State = "playing"
				default:
					result.State = "unknown"
			}
		}
	}
	return
}

func ( w *Wrapper ) GetPlaybackPositions() ( result map[string]PlaybackResult ) {
	// package := w.GetCurrentPackage()
	result = make( map[string]PlaybackResult )
	ms_result := w.Shell( "dumpsys" , "media_session" )
	lines := strings.Split( ms_result , "\n" )
	for line_index , line := range lines {
		// fmt.Println( line_index , line )
		if strings.Contains( line , "active=true" ) {
			session_type_line := lines[ ( line_index - 5 ) ]
			session_type_line = strings.TrimSpace( session_type_line )
			session_parts := strings.Fields( session_type_line )
			// package_str := strings.ToLower( session_parts[ 1 ] )
			package_str := strings.Join( session_parts[ 1 : ] , " " )
			type_str := strings.ToLower( session_parts[ 0 ] )
			x := parse_playback_line( lines[ ( line_index + 4 ) ] )
			x.PackageStr = package_str
			x.Type = type_str
			result[ type_str ] = x
		}
	}
	return
}

func ( w *Wrapper ) GetUpdatedPlaybackPosition( x_input PlaybackResult ) ( result PlaybackResult ) {
	ms_result := w.Shell( "dumpsys" , "media_session" )
	lines := strings.Split( ms_result , "\n" )
	for line_index , line := range lines {
		if strings.Contains( line , "active=true" ) {
			session_type_line := lines[ ( line_index - 5 ) ]
			session_type_line = strings.TrimSpace( session_type_line )
			session_parts := strings.Fields( session_type_line )
			type_str := strings.ToLower( session_parts[ 0 ] )
			if type_str != x_input.Type { continue; }
			result = parse_playback_line( lines[ ( line_index + 4 ) ] )
			result.PackageStr = x_input.PackageStr
			result.Type = type_str
			return
		}
	}
	return
}

func ( w *Wrapper ) WaitOnUpdatedPlaybackPosition( x_input PlaybackResult ) ( result PlaybackResult ) {
	max_tries := 30
	for i := 0; i < max_tries; i++ {
		// fmt.Println( "attempt" , i , "of" , max_tries )
		result = w.GetUpdatedPlaybackPosition( x_input )
		if result.Position != x_input.Position {
			return
		}
		time.Sleep( 500 * time.Millisecond )
	}
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
	return
}

func ( w *Wrapper ) OpenActivity( activity_name string ) ( result string ) {
	result = w.Shell( "am" , "start" , "-n" , activity_name )
	return
}

func ( w *Wrapper ) CloseAppName( app_name string ) ( result string ) {
	result = w.Shell( "am", "force-stop", app_name )
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

func ( w *Wrapper ) GetPackages() ( packages []string ) {
	result := w.Shell( "pm", "list" , "packages" )
	lines := strings.Split( result , "\n" )
	for _ , line := range lines {
		parts := strings.Split( line , "package:" )
		if len( parts ) < 2 { continue }
		packages = append( packages , parts[ 1 ] )
	}
	return
}

func ( w *Wrapper ) GetInstalledPackages() ( packages []string ) {
	result := w.Shell( "pm", "list" , "packages" , "-3" )
	lines := strings.Split( result , "\n" )
	for _ , line := range lines {
		parts := strings.Split( line , "package:" )
		if len( parts ) < 2 { continue }
		packages = append( packages , parts[ 1 ] )
	}
	return
}

func ( w *Wrapper ) GetPackagesLog( package_name string ) ( log_lines []string ) {
	result := w.Shell( "pm", "dump" , package_name )
	log_lines = strings.Split( result , "\n" )
	return
}

func ( w *Wrapper ) GetPackagesDefaultActivity( package_name string ) ( result string ) {
	// result = w.Shell( "cmd" , "package" , "resolve-activity" , "--brief" , package_name , "tail" , "-n" , "1" )
	x := w.Shell( "cmd" , "package" , "resolve-activity" , package_name )
	lines := strings.Split( x , "\n" )
	for _ , line := range lines {
		if strings.HasPrefix( line , "  name=" ) {
			parts := strings.Split( line , "name=" )
			result = parts[ 1 ]
			return
		}
	}
	return
}

// really you need to pull/dump apk and extract stuff from manifest
// adb pull .apk path
// aapt dump badging <pulledfile.apk>
// https://stackoverflow.com/questions/12698814/get-launchable-activity-name-of-package-from-adb
// https://stackoverflow.com/questions/2789462/find-package-name-for-android-apps-to-use-intent-to-launch-market-app-from-web/7502519#7502519
func ( w *Wrapper ) GetPackagesActivities( package_name string ) ( activities []string ) {
	log_lines := w.GetPackagesLog( package_name )
	seen := make( map[ string ] bool )
	package_split_part := fmt.Sprintf( "%s/" , package_name )
	for _ , line := range log_lines {
		line_lower := strings.ToLower( line )
		if strings.Contains( line_lower , "main" ) == false { continue }
		// fmt.Println( line )
		cmp_parts := strings.Split( line , "cmp=" )
		if len( cmp_parts ) > 1 {
			cmp_activity_parts := strings.Split( cmp_parts[ 1 ] , package_split_part )
			if len( cmp_activity_parts ) < 2 { continue }
			cmp_activity := cmp_activity_parts[ 1 ]
			cmp_activity = strings.Fields( cmp_activity )[ 0 ]
			if strings.HasSuffix( cmp_activity , "}" ) {
				cmp_activity = strings.TrimSuffix( cmp_activity , "}" )
			}
			_ , ok := seen[ cmp_activity ]
			if ok == false {
				seen[ cmp_activity ] = true
				activities = append( activities , cmp_activity )
			}
		}
		class_parts := strings.Split( line , "class=" )
		if len( class_parts ) > 1 {
			class_activity := strings.Fields( class_parts[ 1 ] )[ 0 ]
			_ , ok := seen[ class_activity ]
			if ok == false {
				seen[ class_activity ] = true
				activities = append( activities , class_activity )
			}
		}
	}
	return
}


func ( w *Wrapper ) StopAllApps() {
	open_apps := w.GetRunningApps()
	for _ , app := range open_apps {
		w.Shell( "am", "force-stop", app )
	}
	return
}

func ( w *Wrapper ) Sleep() {
	w.Shell( "input" , "keyevent" , "KEYCODE_SLEEP" ) // 223
	return
}

func ( w *Wrapper ) Wakeup() {
	w.Shell( "input" , "keyevent" , "KEYCODE_WAKEUP" ) // 224
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

func ( w *Wrapper ) Landscape() ( result string ) {
	result = w.Shell( "settings" , "put" , "system" , "user_rotation" , "1" )
	// w.Shell( "service" , "call" , "window" , "18" , "i32" , "1" )
	// w.Shell( "setprop" , "persist.demo.hdmirotationlock" , "true" )
	return
}

func ( w *Wrapper ) Portrait() ( result string ) {
	result = w.Shell( "settings" , "put" , "system" , "user_rotation" , "0" )
	// w.Shell( "service" , "call" , "window" , "18" , "i32" , "0" )
	return
}

func ( w *Wrapper ) Brightness( value int ) ( result string ) {
	result = w.Shell( "settings" , "put" , "system" , "screen_brightness" , string( "0" ) )
	// w.Shell( "service" , "call" , "window" , "18" , "i32" , "0" )
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
	text = strings.ReplaceAll( text , "$" , "\\$" )
	result = w.Shell( "input" , "text" , text )
	return
}

// https://github.com/imba28/image-similarity
// https://github.com/baitisj/android_screen_mirror
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/search.go#L36
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/index.go#L10
// https://github.com/imba28/image-similarity/blob/6f921fdf4f5ab8b37d4d563684de99601cc88d5b/pkg/descriptor.go#L16
// https://github.com/hybridgroup/gocv/blob/e11806566cdf2482485cc90d92ed320fa920e91a/cmd/img-similarity/main.go#L123
func ( w *Wrapper ) ScreenshotToFile( save_path string , crop ...int ) ( result string ) {
	utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , "bash" , "-c" ,
		fmt.Sprintf( "%s -s %s exec-out screencap -p > %s" , w.ADBPath , w.Serial , save_path ) ,
	)
	// // TODO , still even clean this up with event.Op&fsnotify.Write == fsnotify.Write {
	// // import "github.com/fsnotify/fsnotify" etc etc
	// // file_stable := make( chan bool )
	for {
		_ , err := os.Stat( save_path )
		if err == nil { break }
		time.Sleep( 10 * time.Millisecond )
	}
	// Wait for the screenshot file size to stabilize
	previous_size := int64( -1 )
	for {
		file_info , err := os.Stat( save_path )
		if err != nil { return }
		size := file_info.Size()
		if size == previous_size { break } // File size has not changed; assuming it's done writing
		previous_size = size
		time.Sleep( 20 * time.Millisecond )
	}
	// watcher.Add( save_path )
	// for {
	// 	select {
	// 	case event := <-watcher.Events:
	// 		// Check for the write event
	// 		if event.Op&fsnotify.Write == fsnotify.Write {
	// 			log.Println("File write detected:", event.Name)
	// 			// Do your work here
	// 			return
	// 		}
	// 	case err := <-watcher.Errors:
	// 		log.Println("Error:", err)
	// 	}
	// }
	fmt.Println( "Screen Shot Captured" )
	return
}

func ( w *Wrapper ) ScreenshotToBytes( crop ...int ) ( result []byte ) {
	rand.Seed( time.Now().UnixNano() )
	random_number := ( rand.Intn( 9000000 ) + 1000000 )
	temp_file , _ := ioutil.TempFile( "" , fmt.Sprintf( "%d-" , random_number ) )
	temp_save_path := temp_file.Name()
	defer os.Remove( temp_save_path )
	utils.ExecProcessWithTimeout( ( EXEC_TIMEOUT * time.Millisecond ) , "bash" , "-c" ,
		fmt.Sprintf( "%s -s %s exec-out screencap -p > %s" , w.ADBPath , w.Serial , temp_save_path ) ,
	)
	// TODO , still even clean this up with event.Op&fsnotify.Write == fsnotify.Write {
	// import "github.com/fsnotify/fsnotify" etc etc
	// file_stable := make( chan bool )
	for {
		_ , err := os.Stat( temp_save_path )
		if err == nil { break }
		time.Sleep( 10 * time.Millisecond )
	}
	// Wait for the screenshot file size to stabilize
	previous_size := int64( -1 )
	for {
		file_info , err := os.Stat( temp_save_path )
		if err != nil { return }
		size := file_info.Size()
		if size == previous_size { break } // File size has not changed; assuming it's done writing
		previous_size = size
		time.Sleep( 20 * time.Millisecond )
	}
	// fmt.Println( "Screen Shot Captured" )

	image_bytes , _ := ioutil.ReadFile( temp_save_path )

	// If we don't have to crop , return early
	if len( crop ) != 4 {
		result = image_bytes
		return
	}

	// Crop
	temp_image_byte_reader := bytes.NewReader( image_bytes )
	temp_image , _ := png.Decode( temp_image_byte_reader )
	x1 , y1 , width , height := crop[ 0 ] , crop[ 1 ] , crop[ 2 ] , crop[ 3 ]
	crop_area := image.Rect( x1 , y1 , ( x1 + width ) , ( y1 + height ) )
	crop_img := temp_image.(*image.NRGBA).SubImage( crop_area ).(*image.NRGBA)
	var crop_buffer bytes.Buffer
	png.Encode( &crop_buffer , crop_img )
	result = crop_buffer.Bytes()
	return
}

func ( w *Wrapper ) ScreenshotToFeatures( crop ...int ) ( result []float64 ) {
	screenshot := w.ScreenshotToBytes( crop... )
	result = image_similarity.GetFeatureVector( screenshot )
	return
}

func ( w *Wrapper ) ScreenshotToPNG( crop ...int ) ( result image.Image ) {
	screenshot := w.ScreenshotToBytes( crop... )
	temp_image_byte_reader := bytes.NewReader( screenshot )
	result , _ = png.Decode( temp_image_byte_reader )
	return
}

func ( w *Wrapper ) GetPixelColor( x int , y int ) ( result color.RGBA ) {
	screenshot := w.ScreenshotToPNG()
	pixel := screenshot.At( x , y )
	r , g , b , a := pixel.RGBA()
	result.R = uint8( r )
	result.G = uint8( g )
	result.B = uint8( b )
	result.A = uint8( a )
	return
}

func ( w *Wrapper ) IsPixelTheSameColor( x int , y int , x_color color.Color ) ( result bool ) {
	pixel_color := w.GetPixelColor( x , y )
	result = ( pixel_color == x_color )
	return
}

func ( w *Wrapper ) CurrentScreenSimilarityToReferenceImage( reference_image_path string , crop ...int ) ( distance float64 ) {
	try.This(func() {
		current_screen_features := w.ScreenshotToFeatures( crop... )
		reference_image_features := image_similarity.GetFeatureVectorFromFilePath( reference_image_path )
		distance = image_similarity.CalculateDistance( current_screen_features , reference_image_features )
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	return
}

func ( w *Wrapper ) SimilarityToFeatureList( features []float64 , reference_image_path string ) ( distance float64 ) {
	try.This(func() {
		current_screen_features := w.ScreenshotToFeatures()
		reference_image_features := image_similarity.GetFeatureVectorFromFilePath( reference_image_path )
		distance = image_similarity.CalculateDistance( current_screen_features , reference_image_features )
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	return
}

func ( w *Wrapper ) IsSameScreen( reference_image_path string , crop ...int ) ( result bool ) {
	distance := w.CurrentScreenSimilarityToReferenceImage( reference_image_path , crop... )
	fmt.Println( "ADBWrapper --> IsSameScreen() --> Distance ===" ,  distance , IMAGE_SIMILARITY_THRESHOLD )
	if distance > IMAGE_SIMILARITY_THRESHOLD {
		result = false
	} else {
		result = true
	}
	return
}

func ( w *Wrapper ) ScreenDistance( reference_image_path string , crop ...int ) ( distance float64 ) {
	distance = w.CurrentScreenSimilarityToReferenceImage( reference_image_path , crop... )
	return
}

func ( w *Wrapper ) IsSameScreenV2( reference_image_path string , crop ...int ) ( result bool , distance float64 ) {
	distance = w.CurrentScreenSimilarityToReferenceImage( reference_image_path , crop... )
	// fmt.Println( distance )
	if distance > IMAGE_SIMILARITY_THRESHOLD {
		result = false
	} else {
		result = true
	}
	return
}

func ( w *Wrapper ) WaitOnScreen( reference_image_path string , timeout time.Duration , crop ...int ) ( result bool ) {
	done := make( chan bool , 1 )

	// Create a timer that will send a message on its channel after the timeout
	timer := time.NewTimer( timeout )

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

	current_screen_features := w.ScreenshotToFeatures( crop... )

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
					distance := w.SimilarityToFeatureList( current_screen_features , imagePath )
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

func ( w *Wrapper ) ClosestScreenInList( file_paths []string , crop ...int ) ( result string ) {
	current_screen_features := w.ScreenshotToFeatures( crop... )

	distances := make( []float64 , len( file_paths ) )
	for i := 0; i < len( file_paths ); i++ {
		reference_image_features := image_similarity.GetFeatureVectorFromFilePath( file_paths[ i ] )
		distance := image_similarity.CalculateDistance( current_screen_features , reference_image_features )
		fmt.Println( file_paths[ i ] , distance )
		distances[ i ] = distance
	}

	min_index := 0
	min_value := math.MaxFloat64
	for i, value := range distances {
		if value < min_value {
			min_value = value
			min_index = i
		}
	}

	result = file_paths[ min_index ]
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


// func color_equals( c1 color.RGBA , c2 color.RGBA , tolerance uint8 ) ( result bool ) {
// 	r_test := abs( int( c1.R ) - int( c2.R ) ) <= int( tolerance )
// 	g_test := abs( int( c1.G ) - int( c2.G ) ) <= int( tolerance )
// 	b_test := abs( int( c1.B ) - int( c2.B ) ) <= int( tolerance )
// 	a_test := abs( int( c1.A ) - int( c2.A ) ) <= int( tolerance )
// 	if r_test && g_test && b_test && a_test { result = true } else { result = false }
// 	return
// }

// func abs( x int ) int {
// 	if x < 0 { return -x }
// 	return x
// }