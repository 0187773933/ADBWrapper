package main

import (
	"fmt"
	"time"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
)

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
	fmt.Println( adb.Exec( "shell" , "dumpsys" , "media_session" ) )
	adb.OpenURI( "spotify:playlist:46CkdOm6pd6tsREVoIgZWw:play" )
	time.Sleep( 1000 * time.Millisecond )
	adb.PressButtonSequence( 21 , 21 , 23 )
	time.Sleep( 1000 * time.Millisecond )
	adb.PressButtonSequence( 22 , 22 , 22 )
	time.Sleep( 400 * time.Millisecond )
	adb.PressButtonSequence( 23 )
	time.Sleep( 2000 * time.Millisecond )
	adb.Tap( 500 , 50 )
}

// brew install opencv@4
// brew link --force opencv@4
// export PKG_CONFIG_PATH="/usr/local/opt/opencv@4/lib/pkgconfig:$PKG_CONFIG_PATH"
// ^^^ add to ~./bash_profile
func main() {
	// adb := adb.ConnectIP(
	// 	"/Users/morpheous/Library/Android/sdk/platform-tools/adb" ,
	// 	"192.168.1.120" ,
	// 	"5555" ,
	// )

	adb := adb_wrapper.ConnectUSB(
		"/usr/local/bin/adb" ,
		"GCC0X8081307034C" ,
	)
	fmt.Println( "Connected" , time.Second )
	// fmt.Println( adb.GetTopWindowInfo() )

	// screenshots/disney_login.png: PNG image data, 1024 x 600, 8-bit/color RGBA, non-interlaced
	// Rect(xMin, yMin, xMax, yMax)
	// but we changed it to be ( x , y , width , height )
	// adb.Screenshot( "screenshots/test_crop.png" , 28 , 337 , 188 , 115 )
	fmt.Println( adb.IsSameScreenV2( "screenshots/test_crop.png" , 28 , 337 , 188 , 115 ) )
}