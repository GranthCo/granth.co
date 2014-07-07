#!/usr/bin/env python 

import subprocess
import sys


def main(argv):
	if len(argv) < 2:
		print "Need two arguments"
		return -1

	primary = argv[1]
	secondary = argv[2]

	if primary == secondary or secondary == "granth":
		print "invalid arguments"
		return -1

	command = "heroku config --app " + primary
	config = subprocess.check_output(command.split())
	for a in config.split('\n'):
		x = a.split(":")
		if len(x) == 2:
			key = x[0].strip()
			val = x[1].strip()
			print "Setting ", key
			command = "heroku config:set " + key + "="+val+" --app " + secondary
			subprocess.check_output(command.split())


if __name__=="__main__":
	sys.exit(main(sys.argv))
