#!/usr/bin/env python3
import subprocess
import re
import threading
import time
import json

def write_json( file_path , python_object ):
	with open( file_path , 'w', encoding='utf-8' ) as f:
		json.dump( python_object , f , ensure_ascii=False , indent=4 )

def read_json( file_path ):
	with open( file_path ) as f:
		return json.load( f )

# https://stackoverflow.com/questions/8647826/simulating-touch-using-adb
# https://web.archive.org/web/20130314094640/http://softteco.blogspot.com/2011/03/android-low-level-shell-click-on-screen.html
# // http://ktnr74.blogspot.com/2013/06/emulating-touchscreen-interaction-with.html
# adb shell getevent -il
# adb shell getevent | grep event2
# adb shell getevent | perl -pe 's/([0-9a-f]{4,8})/sprintf("%d", hex($1))/ge'

# def handle_output(process):
#     while True:
#         output = process.stdout.readline()
#         if output == '' and process.poll() is not None:
#             break
#         if output:
#             line = output.strip().decode('utf-8')
#             match = re.match(r'^(/dev/input/event\d+): ([0-9a-f]+) ([0-9a-f]+) ([0-9a-f]+)$', line)
#             if match:
#                 event_file, type_hex, code_hex, value_hex = match.groups()
#                 type_dec, code_dec, value_dec = [int(x, 16) for x in (type_hex, code_hex, value_hex)]
#                 print(f'{event_file}: {type_dec} {code_dec} {value_dec}')

 # print( event_num ,type_dec, code_dec, value_dec )


# Set the pause threshold (in seconds)
PAUSE_THRESHOLD = 0.5  # for example, 0.5 seconds

streams = []
current_stream = []

def handle_output(process):
	global streams
	global current_stream
	last_event_time = time.time()
	try:
		while True:
			output = process.stdout.readline()
			if output == '' and process.poll() is not None:
				break
			if output:
				line = output.strip().decode('utf-8')
				match = re.match(r'^(/dev/input/event\d+): ([0-9a-f]+) ([0-9a-f]+) ([0-9a-f]+)$', line)
				if match:
					event_file, type_hex, code_hex, value_hex = match.groups()
					type_dec, code_dec, value_dec = [int(x, 16) for x in (type_hex, code_hex, value_hex)]
					event_num = int(re.match(r'/dev/input/event(\d+)', event_file).group(1))

					# If the time since the last event is greater than the pause threshold, start a new stream
					current_time = time.time()
					if current_time - last_event_time > PAUSE_THRESHOLD:
						if current_stream:  # if the current stream is not empty
							streams.append(current_stream)
							print(f'Stream {len(streams)}:')
							for event in current_stream:
								print(event)
							current_stream = []
					last_event_time = current_time

					# Add the event to the current stream
					current_stream.append((event_num, type_dec, code_dec, value_dec))
	except KeyboardInterrupt:
		# On keyboard interrupt, stop the process and exit the loop
		process.terminate()
		return

# Start the process
process = subprocess.Popen(['adb', 'shell', 'getevent'], stdout=subprocess.PIPE)
thread = threading.Thread(target=handle_output, args=(process,))
thread.start()

try:
	thread.join()  # wait for the handle_output function to finish
except KeyboardInterrupt:
	process.terminate()

# Add the last stream if it's not empty
if current_stream:
	streams.append(current_stream)
	print(f'Stream {len(streams)}:')
	for event in current_stream:
		print(event)

write_json( "streams.json" , streams )