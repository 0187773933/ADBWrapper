package main

import (
	"fmt"
	"time"
	"strings"
	// "encoding/json"
	// color "image/color"
	// image_similarity "github.com/0187773933/ADBWrapper/v1/image-similarity"
	adb_wrapper "github.com/0187773933/ADBWrapper/v1/wrapper"
	utils "github.com/0187773933/ADBWrapper/v1/utils"
)

func fire_7_tablet_2019_unlock( adb *adb_wrapper.Wrapper ) {
	adb.Swipe( 522 , 562 , 518 , 230 )
}

func fire_7_tablet_2019_close_all_apps( adb *adb_wrapper.Wrapper ) {
	open_apps := adb.GetRunningPackages()
	if len( open_apps ) < 1 { return; }
	adb.Key( "KEYCODE_HOME" )
	time.Sleep( 500 * time.Millisecond )
	adb.Key( "KEYCODE_APP_SWITCH" )
	time.Sleep( 1 * time.Second )
	for _ , app := range open_apps {
		fmt.Println( "Closing" , app )
		adb.Swipe( 528 , 302 , 528 , 43 )
		time.Sleep( 1 * time.Second )
	}
	adb.Key( "KEYCODE_HOME" )
}

func example_disney_plus_sign_in( adb *adb_wrapper.Wrapper ) {
	// adb.SaveEvents( "disney_login_1.json" )
	adb.OpenPackage( "com.disney.disneyplus" )
	adb.WaitOnScreen( "screenshots/disney_login.png" , ( 5 * time.Second ) )
	// adb.PlaybackEvents( "disney_login_1.json" )
	adb.Swipe( 908 , 293 , 894 , 107 )
	time.Sleep( 1 * time.Second )
	adb.Tap( 510 , 516 )
	adb.Type( "email-part-1" )
	time.Sleep( 100 * time.Millisecond )
	adb.Type( "email-part-2" )
	adb.KeyInt( 61 ) // TAB
	adb.Tap( 508 , 495 )
	adb.Type( "password-part-1" )
	time.Sleep( 100 * time.Millisecond )
	adb.Type( "password-part-2" ) // First Char can't be a '#' ?
	adb.KeyInt( 66 ) // ENTER Key
	adb.WaitOnScreen( "screenshots/disney_profile_selection.png" , ( 5 * time.Second ) )
	adb.Tap( 741 , 245 )
}

func example_spotify_play_playlist( adb *adb_wrapper.Wrapper ) {

	// fmt.Println( adb.Exec( "shell" , "dumpsys" , "media_session" ) )
	adb.ClosePackage( "com.spotify.music" )
	time.Sleep( 1 * time.Second )
	adb.Shell( "am" , "start" , "-n" , "com.spotify.music/com.spotify.music.MainActivity" )
	time.Sleep( 10 * time.Second )
	adb.OpenURI( "spotify:playlist:46CkdOm6pd6tsREVoIgZWw:play" )
	time.Sleep( 1 * time.Second )
	adb.Key( "KEYCODE_MEDIA_PLAY" )

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
	// fmt.Println( adb.GetRunningPackages() )
	// adb.OpenPackage( "tv.twitch.android.viewer" )
	adb.ClosePackage( "tv.twitch.android.viewer" )
	// adb.Key( "KEYCODE_HOME" )
	// adb.OpenURI( fmt.Sprintf( "twitch://stream/%s" , "exbc" ) )
	adb.OpenURI( fmt.Sprintf( "twitch://stream/%s" , "gmhansn" ) )
	// time.Sleep( 10 * time.Second )
}

func example_youtube( adb *adb_wrapper.Wrapper ) {
	// fmt.Println( adb.GetRunningPackages() )
	// adb.OpenPackage( "tv.twitch.android.viewer" )
	adb.StopAllPackages()
	// adb.ClosePackage( "com.amazon.firetv.youtube" )
	// adb.Key( "KEYCODE_HOME" )
	// adb.OpenURI( fmt.Sprintf( "twitch://stream/%s" , "exbc" ) )
	adb.SetBrightness( 0 )
	adb.OpenURI( "https://www.youtube.com/watch?v=bOwsLtwa2Ts" )
	// time.Sleep( 10 * time.Second )
}

func example_disney( adb *adb_wrapper.Wrapper ) {
	adb.StopAllPackages()
	adb.SetBrightness( 0 )
	adb.Shell( "am" , "start" , "-a" , "android.intent.action.VIEW" , "-d" , "" )
	// adb.OpenURI( fmt.Sprintf( "https://www.disneyplus.com/video/%s" , "" ) )
}

func example_vlc( adb *adb_wrapper.Wrapper ) {
	adb.StopAllPackages()
	adb.SetBrightness( 0 )
	adb.OpenURI( fmt.Sprintf( "vlc://%s" , "" ) )
}

func example_spotify( adb *adb_wrapper.Wrapper ) {

	// TODO : TV Volume Off
	// adb.SetVolumePercent( 0 )
	adb.StopAllPackages()
	adb.SetBrightness( 0 )
	adb.ClosePackage( "com.spotify.tv.android" )
	time.Sleep( 1 * time.Second )
	playlist_uri := fmt.Sprintf( "spotify:playlist:%s:play" , "6ZFJWltDYCI0OVyXNleN9e" )
	adb.OpenURI( playlist_uri )

	// Enable Shuffle
	// time.Sleep( 10 * time.Second )
	// ( 0 , 0 ) = Top-Left
	// adb.Screenshot( "./screenshots/spotify/shuffle_off.png" , 735 , 925 , 35 , 45 )

	// adb.Key( "KEYCODE_ENTER" )
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
		adb.Key( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		adb.Key( "KEYCODE_DPAD_LEFT" )
		time.Sleep( 200 * time.Millisecond )
		adb.Key( "KEYCODE_ENTER" )
		time.Sleep( 200 * time.Millisecond )
		adb.Key( "KEYCODE_MEDIA_NEXT" )
		// adb.SetVolumePercent( 100 )
		time.Sleep( 200 * time.Millisecond )
		adb.Key( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 200 * time.Millisecond )
		adb.Key( "KEYCODE_DPAD_RIGHT" )
		time.Sleep( 200 * time.Millisecond )
		adb.Key( "KEYCODE_DPAD_RIGHT" )
	} else {
		// TODO : Turn TV Volume On
		// adb.SetVolumePercent( 100 )
	}

	time.Sleep( 10 * time.Second )
	adb.OpenURI( fmt.Sprintf( "spotify:playlist:%s:play" , "3UMDmO2YJb8DgUjpSBu8y9" ) )
	time.Sleep( 500 * time.Millisecond )
	adb.Key( "KEYCODE_MEDIA_NEXT" )

}


func fire_7_tablet_2019_netflix( adb *adb_wrapper.Wrapper ) {
	// netflix_package := "com.netflix.ninja"
	netflix_package := "com.netflix.mediaclient"
	// netflix_source_activity := "com.netflix.ninja/.MainActivity"
	netflix_source_activity := "com.netflix.mediaclient/.ui.launch.NetflixComLaunchActivity"
	uri := "https://www.netflix.com/watch/80223868?trackId=14315607"
	fmt.Println( uri )
	adb_status := adb.GetStatus()
	fmt.Println( adb_status )
	adb.StopAllPackages()
	adb.ClosePackage( netflix_package )
	time.Sleep( 500 * time.Millisecond )
	adb.OpenPackage( netflix_package )
	adb.Shell(
		"am" , "start" , "-c" , "android.intent.category.LEANBACK_LAUNCHER" ,
		"-a" , "android.intent.action.VIEW" , "-d" , uri ,
		"-f" , "0x10008000" ,
		"-e" , "source" , "30" , netflix_source_activity ,
	)

	// have to treat everything as nested / transparent windows on top of windows
	// players on top of players
	// so everything is in a map
	// mostly only firetablet ( older android ) has these proxy objects

	// players := adb.FindPlayers( "netflix" )
	fmt.Println( "waiting 20 seconds for netflix player to appear" )
	netflix_players := adb.WaitOnPlayers( "netflix" , 20 )
	if len( netflix_players ) < 1 {
		fmt.Println( "never started playing , we might have to try play button" )
	}
	fmt.Println( "netflix player should be ready" )
	utils.PrettyPrint( netflix_players )
	fmt.Println( "waiting 10 seconds to see if netflix auto starts playing" )
	playing := adb.WaitOnPlayersPlaying( "netflix" , 10 )
	if len( playing ) < 1 {
		fmt.Println( "never started playing , we might have to try play button" )
	}
	utils.PrettyPrint( playing )
	fmt.Println( "total now playing" , len( playing ) )
	// x := adb.GetNowPlaying( "netflix" , 60 )
	// x := adb.GetNowPlayingForce( "netflix" , 60 )
	for _ , player := range playing {
		if player.Updated > 0 {
			fmt.Println( "netflix autostarted playing on it's own" )
			return
		}
	}
	fmt.Println( "trying to force update adb playback state" )
	for _ , player := range playing {
		playing = adb.WaitOnPlayersUpdatedForce( "netflix" , player.Updated , 60 )
		utils.PrettyPrint( playing )
	}
	fmt.Println( "trying to force update adb playback state" )
	for _ , player := range playing {
		playing = adb.WaitOnPlayersUpdatedForce( "netflix" , player.Updated , 60 )
		utils.PrettyPrint( playing )
	}
}


// brew install opencv@4
// brew link --force opencv@4
// export PKG_CONFIG_PATH="/usr/local/opt/opencv@4/lib/pkgconfig:$PKG_CONFIG_PATH"
// ^^^ add to ~./bash_profile

// mFocusedApp=Token{d00eec4 ActivityRecord{ca7c5d7 u0 com.amazon.firetv.youtube/dev.cobalt.app.MainActivity t265}}

func main() {

	adb := adb_wrapper.ConnectIP(
		"/usr/local/bin/adb" ,
		"192.168.4.193" , // firecube
		// "192.168.4.56" , // firestick
		"5555" ,
	)

	// adb := adb_wrapper.ConnectUSB(
	// 	"/usr/local/bin/adb" ,
	// 	"GCC0X8081307034C" , // firetablet
	// )

	// status := adb.GetStatus()
	// utils.PrettyPrint( status )
	// fmt.Println( adb.IsSearchTermActivityOpen( "ProfileSelection" ) )

	fmt.Println( adb.ScreenshotToFile( "screenshots/netflix/profile-selection.png" ) )
	// fmt.Println( "waiting" )
	// adb.WaitOnPixelColor( 1694 , 96 , color.RGBA{ R: 28 , G: 231 , B: 131 , A: 255 } , 10 * time.Second )
	// fmt.Println( "done" )

	// adb.WaitOnPixelColor( x int , y int , x_color color.Color , timeout time.Duration )

	// pss := "/Users/morpheous/WORKSPACE/GO/FireC2Server/SAVE_FILES/screenshots/disney/profile_selection.png"
	// pss_features := image_similarity.GetFeatureVectorFromFilePath( pss )
	// screenshot_bytes := adb.ScreenshotToBytes()
	// screenshot_features := adb.ImageBytesToFeatures( &screenshot_bytes )
	// distance := image_similarity.CalculateDistancePoint( &screenshot_features , &pss_features )
	// utils.PrettyPrint( distance )
	// pixel_color := adb.GetPixelColorFromImageBytes( &screenshot_bytes , 896 , 469 )
	// utils.PrettyPrint( pixel_color )

	// start := time.Now()
	// status := adb.GetStatus()
	// elapsed := time.Since( start )
	// utils.PrettyPrint( status )
	// fmt.Println( "GetStatus() took" , elapsed )
	// start = time.Now()
	// status.ScreenShot = adb.ScreenshotToBytes()
	// fmt.Println( "ScreenshotToBytes()" , len( status.ScreenShot ) , "took an extra" , time.Since( start ) )
	// adb.ScreenshotToFile( "test.png" )


	// time.Sleep( 5 * time.Minute )
	// fire_7_tablet_2019_netflix( &adb )

	// adb.ScreenOn()
	// adb.SetBrightness( 100 )
	// adb.Home()
	// adb.PowerOff()

	// activity := adb.GetPackagesDefaultActivity( "tv.twitch.android.viewer" )
	// fmt.Println( activity )

	// packages := adb.GetInstalledPackages()
	// utils.PrettyPrint( packages )
	// for _ , x_package := range packages {
	// 	activities := adb.GetPackagesActivitiesSearch( x_package )
	// 	fmt.Println( x_package , activities )
	// }

	// activities := adb.GetPackagesActivitiesSearch( "tv.twitch.android.viewer" )
	// utils.PrettyPrint( activities )

	// packages := adb.GetInstalledPackages()
	// fmt.Println( packages )

	// twitch_apk_path := adb.GetPackagePath( "tv.twitch.android.viewer" )
	// fmt.Println( twitch_apk_path )
	// adb.PullPackageAPK( "tv.twitch.android.viewer" , "./twitch.apk" )

	// activities := adb.GetPackagesActivitiesPull( "tv.twitch.android.viewer" )
	// fmt.Println( activities )

	// activities := adb.GetPackagesActivities( "tv.twitch.android.viewer" )
	// fmt.Println( activities )

	// packages := adb.GetInstalledPackagesAndActivities()
	// // utils.PrettyPrint( packages )
	// utils.WriteJSON( "./packages-firecube.json" , packages )

	// windows := adb.GetWindowStack()
	// utils.PrettyPrint( windows )

	// current_window := adb.GetTopWindow()
	// utils.PrettyPrint( current_window )

	// adb.Screenshot( "screenshots/spotify/shuffle_off_new.png" , 735 , 957 , 35 , 15 )
	// adb.ScreenshotToFile( "screenshots/spotify/new_position_5.png" )

	// adb.GetCurrentPackage()

	// positions := adb.GetPlaybackPositions()
	// utils.PrettyPrint( positions )
	// time.Sleep( 1 * time.Second )
	// adb.Pause()
	// time.Sleep( 100 * time.Millisecond )
	// adb.Play()
	// new_positions := adb.GetPlaybackPositions()
	// for key , _ := range positions {
	// 	if new_positions[ key ] != positions[ key ] {
	// 		fmt.Println( key , "changed" , positions[ key ] , new_positions[ key ] )
	// 	}
	// }

	// updated := adb.GetUpdatedPlaybackPosition( positions[ "playbackmediasessionwrapper" ] )
	// updated := adb.WaitOnUpdatedPlaybackPosition( positions[ "playbackmediasessionwrapper" ] )
	// utils.PrettyPrint( updated )
	// updated_positions := adb.GetPlaybackPositions()
	// utils.PrettyPrint( updated_positions )
	// test , pos := adb.GetPlaybackPositionForce()
	// fmt.Println( test , pos )

	// fmt.Println( adb.GetPlaybackPosition() )
	// positions := adb.GetPlaybackPositions()
	// spotify := positions[ "spotify-android-tv-media-session" ]
	// spotify := positions[ "netflix" ]
	// fmt.Println( spotify )
	// adb.Shell( "input" , "keyevent" , "KEYCODE_MEDIA_PLAY_PAUSE" )
	// time.Sleep( 100 * time.Millisecond )
	// adb.Shell( "input" , "keyevent" , "KEYCODE_MEDIA_PLAY_PAUSE" )
	// spotify = adb.GetUpdatedPlaybackPosition( spotify )
	// fmt.Println( spotify )

	// fmt.Println( positions[ "netflix" ] )
	// updated := adb.WaitOnUpdatedPlaybackPosition( positions[ "netflix" ] )
	// fmt.Println( updated )

	// fmt.Println( adb.GetPlaybackPositionTest() )
	// fmt.Println( adb.GetPlaybackPositionTest() )
	// fmt.Println( adb.GetPlaybackPositionTest() )
	// fmt.Println( adb.GetPlaybackPositionTest() )
	// fmt.Println( adb.GetPlaybackPositionTest() )
	// fmt.Println( adb.GetPlaybackPositionTest() )

	// status := adb.GetStatus()
	// status_json , _ := json.MarshalIndent( status , "", "    " )
	// fmt.Println( string( status_json ) )
	// fmt.Println( adb.GetCPUArchitecture() )

	// adb.Key( "KEYCODE_DPAD_LEFT" )
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
	// adb.StopAllPackages()
	// fire_7_tablet_2019_close_all_apps( &adb )
	// adb.Key( "KEYCODE_HOME" )
	// fmt.Println( "ready" )

	// example_twitch( &adb )
	// adb.ScreenOff()
	// fmt.Println( adb.GetRunningPackages() )
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