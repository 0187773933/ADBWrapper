package main

import (
	"fmt"
	"time"
	adb "github.com/0187773933/ADBWrapper/v1/wrapper"
)

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

	adb := adb.ConnectUSB(
		"/usr/local/bin/adb" ,
		"GCC0X8081307034C" ,
	)
	fmt.Println( "Connected" , time.Second )
	// fmt.Println( adb.GetTopWindowInfo() )

	// fmt.Println( adb.Exec( "shell" , "dumpsys" , "media_session" ) )
	// adb.OpenURI( "spotify:playlist:46CkdOm6pd6tsREVoIgZWw:play" )
	// time.Sleep( 1000 * time.Millisecond )
	// adb.PressButtonSequence( 21 , 21 , 23 )
	// time.Sleep( 1000 * time.Millisecond )
	// adb.PressButtonSequence( 22 , 22 , 22 )
	// time.Sleep( 400 * time.Millisecond )
	// adb.PressButtonSequence( 23 )
	// time.Sleep( 2000 * time.Millisecond )
	// adb.Tap( 500 , 50 )


	// adb.Screenshot( "screenshots/disney_profile_selection.png" )
	// same_screen := adb.IsSameScreen(  )
	// fmt.Pritnln( same )

	// closest := adb.ClosestScreen( "/Users/morpheous/WORKSPACE/GO/ADBWrapper/screenshots" )
	// fmt.Println( closest )

	// adb.Swipe( 908 , 293 , 894 , 107 )
	// fmt.Println( adb.Exec( "shell" , "getevent" , "-il" ) )
	// fmt.Println( adb.GetEventDevices() )

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