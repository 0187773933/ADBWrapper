<center><h1>ADB Wrapper</h1></center>

- apk's === arm-v7a , Android 5.0+
- https://www.adb-shell.com/android/media.html

## TODO

- better current state meta info

## Get Activites from APK

- https://stackoverflow.com/questions/6547703/list-all-activities-within-an-apk-from-the-shell

1. `brew install --cask android-sdk`
2. `export ANDROID_HOME=/usr/local/share/android-sdk`
3. `yes | sdkmanager --licenses`
4. `sdkmanager --list`
5. `sdkmanager "build-tools;34.0.0"`
6. `/usr/local/share/android-sdk/build-tools/34.0.0/aapt --version`
7. `sudo ln -s "/usr/local/share/android-sdk/build-tools/34.0.0/aapt" /usr/local/bin/aapt`
8. `aapt dump xmltree <apk-file> AndroidManifest.xml > apk-file.txt`
9. `aapt list -a /path/to/the/apk > apk-file.all.txt`

## ADB Info Links
- https://gist.github.com/mcfrojd/9e6875e1db5c089b1e3ddeb7dba0f304
- https://stackoverflow.com/questions/12698814/get-launchable-activity-name-of-package-from-adb
- https://github.com/Swind/pure-python-adb
- https://github.com/google/python-adb
- https://www.droidviews.com/adb-fastboot-commands-android/
- https://stackoverflow.com/questions/2789462/find-package-name-for-android-apps-to-use-intent-to-launch-market-app-from-web/7502519#7502519
- https://stackoverflow.com/questions/21034168/android-start-an-activity-from-command-line-using-intent-uri#21035826
- https://github.com/IntergalacticPenguin/mobile-toolkit
- https://paul.kinlan.me/deep-app-linking-on-android-and-chrome/
- http://adbcommand.com/awesome-adb
- http://www.dre.vanderbilt.edu/~schmidt/android/android-4.0/out/target/common/docs/doc-comment-check/guide/developing/tools/ddms.html
- https://www.zebra.com/us/en/support-downloads/knowledge-articles/evm/Using-ADB-to-Determine-Package-Name-and-Launchable-Activity.html
- https://stackoverflow.com/questions/34450789/android-deep-linking-with-custom-uri
- https://medium.com/airbnb-engineering/deeplinkdispatch-778bc2fd54b7
- https://github.com/airbnb/DeepLinkDispatch
- https://play.google.com/store/apps/details?id=com.manoj.dlt
- https://stackoverflow.com/questions/11142843/how-can-i-use-adb-to-send-a-longpress-key-event#20537408
- https://medium.com/@roshni.b.tiwari/adb-commands-for-music-video-and-fm-b02018219cd5
- https://www.guidingtech.com/15008/adb-control-keypress-broken-android-keys/
- https://stackoverflow.com/questions/7789826/adb-shell-input-events#8483797
- https://github.com/sncvikas/ADB_Commands
- https://stackoverflow.com/questions/27988069/query-android-content-provider-from-command-line-adb-shell
- https://stackoverflow.com/questions/27391326/how-to-change-the-volume-using-adb-shell-service-call-audio#27392847
- https://www.androidauthority.com/about-android-debug-bridge-adb-21510/
- https://www.apkmirror.com/apk/instagram/instagram-instagram/instagram-instagram-201-0-0-0-29-release/instagram-201-0-0-0-29-android-apk-download/download/#uploadAPK

## ADB Code Examples
- https://github.com/ceberous/ShzmTwitchBot/blob/master/main.js
- https://github.com/holmeszyx/goadb/blob/master/goadb.go
- https://gist.github.com/Pulimet/5013acf2cd5b28e55036c82c91bd56d8
- https://gist.github.com/bbsusheelkumar/357c62202bda3d905d589046b95d157b/revisions
- http://leesei.github.io/android-content/


## ADB Misc Commands

- `adb shell`

- `adb shell getprop ro.serialno`

- `pm list packages | sed -e "s/package://" | while read x; do cmd package resolve-activity --brief $x | tail -n 1 | grep -v "No activity found"; done`

- `adb shell am start tv.twitch.android.viewer/tv.twitch.android.apps.TwitchActivity`

- `adb shell input touchscreen tap 100 100`

- `adb shell am start -W -a android.intent.action.VIEW -d "example://gizmos" com.example.android`

- `adb shell pm list packages -f`

- `adb shell dumpsys window windows`

- `adb shell dumpsys window windows | find "mcurrentfocus"`

- `adb shell dumpsys window windows | grep -E 'mCurrentFocus|mFocusedApp'`

- `adb shell dumpsys display`

- `intent://#Intent;package=qualified.package.name;end`

- `adb shell am to-uri tv.twitch.android.viewer/tv.twitch.android.apps.TwitchActivity`

- `adb shell am start -W -a android.intent.action.VIEW -d "example://gizmos" com.example.android`

- `adb shell am start -W -a android.intent.action.VIEW -d "intent://#Intent;action=android.intent.action.MAIN;category=android.intent.category.LAUNCHER;component=tv.twitch.android.viewer/tv.twitch.android.apps.TwitchActivity;end" com.example.android`

- `content query --uri content://settings/global`

- `adb shell monkey -p <insert-package-name> -c android.intent.category.LAUNCHER 1`

- `adb shell input touchscreen text "Wow,\ it\ so\ cool\ feature"`

- `adb shell am start -n tv.twitch.android.viewer/tv.twitch.android.apps.TwitchActivity --es "toast" "hello, world"`

- `adb shell pm dump tv.twitch.android.viewer`

- `adb shell "cmd package resolve-activity tv.twitch.android.viewer"`

- `adb shell logcat -d | grep 'START u0' | tail -n 1 | sed 's/.*cmp=\(.*\)} .*/\1/g'`

- `adb shell dumpsys package | grep -Eo "^[[:space:]]+[0-9a-f]+[[:space:]]+tv.twitch.android.viewer/[^[:space:]]+" | grep -oE "[^[:space:]]+$"`

- `adb shell service call audio 3 i32 3 i32 0 i32`

- `adb exec-out screencap -p > test.png`

- `adb install instagram.apk`


## Get Info

- `adb shell dumpsys window windows | grep -E 'mCurrentFocus|mFocusedApp'`
- `adb shell dumpsys media_session | grep "state=PlaybackState"`
- `adb shell dumpsys media_session`
- `adb shell content query --uri content://call_log/calls`
- `adb shell getevent -il`

## Twitch

- https://dev.twitch.tv/docs/mobile-deeplinks
````
adb shell am start -a android.intent.action.VIEW -d "twitch://stream/chess24" && \
adb shell input keyevent 4
````

## Disney

- `adb shell am start -a android.intent.action.VIEW -d "https://www.disneyplus.com/video/74351ae5-f6cd-4464-aeb6-8f4b10ca2649"`


## Spotify

- Startup = `adb shell am start -a android.intent.action.VIEW -d "spotify:playlist:46CkdOm6pd6tsREVoIgZWw:play"`
- Startup-Random =
```bash
adb shell am start -a android.intent.action.VIEW -d "spotify:playlist:46CkdOm6pd6tsREVoIgZWw:play" && \
adb shell input keyevent 21 21 23 && sleep 1 && \
adb shell input keyevent 22 22 22 && sleep 0.4 && adb shell input keyevent 23 && sleep 2 \
adb shell input tap 500 500
```