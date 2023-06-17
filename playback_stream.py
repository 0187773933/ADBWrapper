#!/usr/bin/env python3
import json
import subprocess

def adbshell(commands):
	subprocess.run(["adb", "shell", ' '.join(commands)])

def playback_stream(stream):
	commands = []
	for event in stream:
		event_num, type_dec, code_dec, value_dec = event
		event_file = f"/dev/input/event{event_num}"
		command = f'S="sendevent {event_file}";$S {type_dec} {code_dec} {value_dec};'
		commands.append(command)
	adbshell(commands)

if __name__ == "__main__":
	# Read the streams from the JSON file
	with open('streams.json', 'r') as f:
		streams = json.load(f)

	for index , stream in enumerate(streams):
		playback_stream( stream )