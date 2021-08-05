package main

import (
	// "fmt"
	// "time"
	adb "github.com/0187773933/ADBWrapper/v1/wrapper"
)

func main() {
	adb := adb.Connect(
		"/Users/morpheous/Library/Android/sdk/platform-tools/adb" ,
		"192.168.1.120" ,
		"5555" ,
	)
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

	// adb.Screenshot()
	adb.CurrentScreenSimilarityToReferenceImage( "/Users/morpheous/WORKSPACE/GO/ADBWrapper/profile_selection.png" )
	adb.CurrentScreenSimilarityToReferenceImage( "/Users/morpheous/WORKSPACE/GO/ADBWrapper/homescreen.png" )
}