#!/bin/sh 
case "$(uname)" in
	Darwin)
		model=$(system_profiler SPHardwareDataType | sed -nE 's/.*Model Identifier: (.*)/\1/p')
	;;
	Linux)
		model=$(cat /sys/devices/virtual/dmi/id/board_{name,version} 2> /dev/null | tr '\n' ' ' 2> /dev/null)
		[ -z "$model" ] && model=$(cat /sys/firmware/devicetree/base/model)
	;;
	FreeBSD)
		which doas &> /dev/null && elevate=doas || elevate=sudo
		model=$($elevate dmidecode --type system 2> /dev/null | sed -nE 's/.*Product Name: (.*)/\1/p')
	;;
esac
printf "$model $(uname -rms)"

