#!/bin/sh
case "$(uname)" in
	Darwin)
		model=$(system_profiler SPHardwareDataType | sed -nE 's/.*Model Identifier: (.*)/\1/p')
	;;
	Linux)
		model=$(cat /sys/devices/virtual/dmi/id/board_{name,version} 2> /dev/null | tr '\n' ' ' | sed 's/None//g')
		[ -z "$model" ] && model=$(cat /sys/firmware/devicetree/base/model 2> /dev/null)
	;;
	FreeBSD)
		which doas &> /dev/null && elevate=doas || elevate=sudo
		model=$($elevate dmidecode --type baseboard 2> /dev/null | sed -nE 's/.*Product Name: (.*)/\1/p')
	;;
esac

[ -z "$model" ] && 
	uname -rms ||
	printf "$model $(uname -rms)"
