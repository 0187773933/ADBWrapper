package main

import (
	"fmt"
	"time"
	"strings"
	"encoding/json"
	// color "image/color"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
)

func fire_7_tablet_2019_unlock( adb *adb_wrapper.Wrapper ) {
	adb.Swipe( 522 , 562 , 518 , 230 )
}

func fire_7_tablet_2019_close_all_apps( adb *adb_wrapper.Wrapper ) {
	open_apps := adb.GetRunningApps()
	if len( open_apps ) < 1 { return; }
	adb.PressKeyName( "KEYCODE_HOME" )
	time.Sleep( 500 * time.Millisecond )
	adb.PressKeyName( "KEYCODE_APP_SWITCH" )
	time.Sleep( 1 * time.Second )
	for _ , app := range open_apps {
		fmt.Println( "Closing" , app )
		adb.Swipe( 528 , 302 , 528 , 43 )
		time.Sleep( 1 * time.Second )
	}
	adb.PressKeyName( "KEYCODE_HOME" )
}

func example_disney_plus_sign_in( adb *adb_wrapper.Wrapper ) {
	// adb.SaveEvents( "disney_login_1.json" )
	adb.OpenAppName( "com.disney.disneyplus" )
	adb.WaitOnScreen( "screenshots/disney_login.png" , ( 5 * time.Second ) )
	// adb.PlaybackEvents( "disney_login_1.json" )
	adb.Swipe( 908 , 293 , 894 , 107 )
	time.Sleep( 1 * time.Second )
	adb.Tap( 510 , 516 )
	adb.Type( "email-part-1" )
	time.Sleep( 100 * time.Millisecond )
	adb.Type( "email-part-2" )
	adb.PressKey( 61 ) // TAB
	adb.Tap( 508 , 495 )
	adb.Type( "password-part-1" )
	time.Sleep( 100 * time.Millisecond )
	adb.Type( "password-part-2" ) // First Char can't be a '#' ?
	adb.PressKey( 66 ) // ENTER Key
	adb.WaitOnScreen( "screenshots/disney_profile_selection.png" , ( 5 * time.Second ) )
	adb.Tap( 741 , 245 )
}

func example_spotify_play_playlist( adb *adb_wrapper.Wrapper ) {

	// fmt.Println( adb.Exec( "shell" , "dumpsys" , "media_session" ) )
	adb.CloseAppName( "com.spotify.music" )
	time.Sleep( 1 * time.Second )
	adb.Shell( "am" , "start" , "-n" , "com.spotify.music/com.spotify.music.MainActivity" )
	time.Sleep( 10 * time.Second )
	adb.OpenURI( "spotify:playlist:46CkdOm6pd6tsREVoIgZWw:play" )
	time.Sleep( 1 * time.Second )
	adb.PressKeyName( "KEYCODE_MEDIA_PLAY" )

	// time.Sleep( 1000 * time.Millisecond )
	// adb.PressButtonSequence( 21 , 21 , 23 )
	// time.Sleep( 1000 * time.Millisecond )
	// adb.PressButtonSequence( 22 , 22 , 22 )
	// time.Sleep( 400 * time.Millisecond )
	// adb.PressButtonSequence( 23 )
	// time.Sleep( 2000 * time.Millisecond )
	// adb.Tap( 500 , 50 )

	// adb.SetVolumePercent( 70 )
	adb.SetVolumePercent( 56 )
	// adb.SetVolume( 15 )
}

func example_twitch( adb *adb_wrapper.Wrapper ) {
	// fmt.Println( adb.GetRunningApps() )
	// adb.OpenAppName( "tv.twitch.android.viewer" )
	adb.CloseAppName( "tv.twitch.android.viewer" )
	// adb.PressKeyName( "KEYCODE_HOME" )
	// adb.OpenURI( fmt.Sprintf( "twitch://stream/%s" , "exbc" ) )
	adb.OpenURI( fmt.Sprintf( "twitch://stream/%s" , "gmhansn" ) )
	// time.Sleep( 10 * time.Second )
}

func example_youtube( adb *adb_wrapper.Wrapper ) {
	// fmt.Println( adb.GetRunningApps() )
	// adb.OpenAppName( "tv.twitch.android.viewer" )
	adb.StopAllApps()
	// adb.CloseAppName( "com.amazon.firetv.youtube" )
	// adb.PressKeyName( "KEYCODE_HOME" )
	// adb.OpenURI( fmt.Sprintf( "twitch://stream/%s" , "exbc" ) )
	adb.Brightness( 0 )
	adb.OpenURI( "https://www.youtube.com/watch?v=bOwsLtwa2Ts" )
	// time.Sleep( 10 * time.Second )
}

func example_disney( adb *adb_wrapper.Wrapper ) {
	adb.StopAllApps()
	adb.Brightness( 0 )
	adb.Shell( "am" , "start" , "-a" , "android.intent.action.VIEW" , "-d" , "" )
	// adb.OpenURI( fmt.Sprintf( "https://www.disneyplus.com/video/%s" , "" ) )
}

func example_vlc( adb *adb_wrapper.Wrapper ) {
	adb.StopAllApps()
	adb.Brightness( 0 )
	adb.OpenURI( fmt.Sprintf( "vlc://%s" , "" ) )
}

func example_spotify( adb *adb_wrapper.Wrapper ) {

	// TODO : TV Volume Off
	// adb.SetVolumePercent( 0 )
	adb.StopAllApps()
	adb.Brightness( 0 )
	adb.CloseAppName( "com.spotify.tv.android" )
	time.Sleep( 1 * time.Second )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , "6ZFJWltDYCI0OVyXNleN9e" )
	adb.OpenURI( playlist_uri )

	// Enable Shuffle
	// time.Sleep( 10 * time.Second )
	// ( 0 , 0 ) = Top-Left
	// adb.Screenshot( "./screenshots/spotify/shuffle_off.png" , 735 , 925 , 35 , 45 )

	// adb.PressKeyName( "KEYCODE_ENTER" )
	adb.WaitOnScreen( "./screenshots/spotify/playing.png" , ( 10 * time.Second ) , 945 , 925 , 30 , 30 )
	fmt.Println( "Ready" )
	time.Sleep( 1 * time.Second )
	shuffle_test := adb.ClosestScreenInList( []string{
			"./screenshots/spotify/shuffle_off.png" ,
			"./screenshots/spotify/shuffle_on.png" ,
		} ,
		735 , 925 , 35 , 45 ,
	)
	if strings.Contains( shuffle_test , "off" ) {
		adb.PressKeyName( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_ENTER" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_MEDIA_NEXT" )
		// adb.SetVolumePercent( 100 )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 200 * time.Millisecond )
		adb.PressKeyName( "KEYCODE_DPAD_RIGHT" )
	} else {
		// TODO : Turn TV Volume On
		// adb.SetVolumePercent( 100 )
	}

	time.Sleep( 10 * time.Second )
	adb.OpenURI( fmt.Sprintf( "spotify:playlist:%s:play" , "3UMDmO2YJb8DgUjpSBu8y9" ) )
	time.Sleep( 500 * time.Millisecond )
	adb.PressKeyName( "KEYCODE_MEDIA_NEXT" )

}


// brew install opencv@4
// brew link --force opencv@4
// export PKG_CONFIG_PATH="/usr/local/opt/opencv@4/lib/pkgconfig:$PKG_CONFIG_PATH"
// ^^^ add to ~./bash_profile

// mFocusedApp=Token{d00eec4 ActivityRecord{ca7c5d7 u0 com.amazon.firetv.youtube/dev.cobalt.app.MainActivity t265}}

func main() {

	adb := adb_wrapper.ConnectIP(
		"/usr/local/bin/adb" ,
		"192.168.4.193" ,
		"5555" ,
	)

	// adb.Screenshot( "screenshots/spotify/shuffle_off_new.png" , 735 , 957 , 35 , 15 )
	// adb.ScreenshotToFile( "screenshots/spotify/new_position_5.png" )

	status := adb.GetStatus()
	status_json , _ := json.MarshalIndent( status , "", "    " )
	fmt.Println( string( status_json ) )
	fmt.Println( adb.GetCPUArchitecture() )

	// adb.PressKeyName( "KEYCODE_DPAD_LEFT" )
	// white := color.RGBA{ R: 255 , G: 255 , B: 255 , A: 255 }
	// // shuffle_pixel_color := adb.GetPixelColor( 752 , 964 )
	// if adb.IsPixelTheSameColor( 752 , 964 , white ) == true {
	// 	fmt.Println( "Shuffle === ON" )
	// } else {
	// 	fmt.Println( "Shuffle === OFF" )
	// }

	// shuffle_test := adb.ClosestScreenInList( []string{
	// 		"./screenshots/spotify/shuffle_off.png" ,
	// 		"./screenshots/spotify/shuffle_on.png" ,
	// 	} ,
	// 	735 , 957 , 35 , 15 ,
	// )
	// fmt.Println( shuffle_test )

	// adb.WaitOnScreen( "./screenshots/spotify/playing.png" , ( 10 * time.Second ) , 945 , 925 , 30 , 30 )
	// adb.WaitOnScreen( "./screenshots/spotify/playing.png" , ( 10 * time.Second ) )

	// adb := adb_wrapper.ConnectUSB(
	// 	"/usr/local/bin/adb" ,
	// 	"GCC0X8081307034C" ,
	// )
	// if adb.ForceScreenOn() == true {
	// 	fire_7_tablet_2019_unlock( &adb )
	// }
	// adb.DisableScreenTimeout()
	// adb.StopAllApps()
	// fire_7_tablet_2019_close_all_apps( &adb )
	// adb.PressKeyName( "KEYCODE_HOME" )
	// fmt.Println( "ready" )

	// example_twitch( &adb )
	// adb.ScreenOff()
	// fmt.Println( adb.GetRunningApps() )
	// example_youtube( &adb )
	// example_spotify( &adb )

	// fmt.Println( adb.GetWindowStack() )



	// fmt.Println( adb.GetTopWindowInfo() )

	// example_spotify_play_playlist( &adb )

	// black-screen
	// adb.Screenshot( "spotify-login.png" )


	// screenshots/disney_login.png: PNG image data, 1024 x 600, 8-bit/color RGBA, non-interlaced
	// Rect(xMin, yMin, xMax, yMax)
	// but we changed it to be ( x , y , width , height )
	// adb.Screenshot( "screenshots/test_crop.png" , 28 , 337 , 188 , 115 )
	// fmt.Println( adb.IsSameScreenV2( "screenshots/test_crop.png" , 28 , 337 , 188 , 115 ) )

	// adb.WaitOnScreen( "screenshots/test_crop.png" , ( 30 * time.Second ) , 28 , 337 , 188 , 115 )
	// fmt.Println( "found" )
}