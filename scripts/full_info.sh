#!/bin/sh 
# Invoke:
# 	ssh $host < host_info.sh 2> /dev/null
# View colors with:
#	while read -r line; do echo -e $line; done < ./scripts/full_info.sh
u=$(uname)
case $u in
	Darwin)
		model="\033[97m \033[0m$(system_profiler SPHardwareDataType | sed -nE 's/.*Model Identifier: (.*)/\1/p')"
	;;
	Linux)
		model=$(cat /sys/devices/virtual/dmi/id/board_{name,version} 2> /dev/null | tr '\n' ' ' 2> /dev/null)
		[ -z "$model" ] && model=$(cat /sys/firmware/devicetree/base/model)

		case "$(sed -nE 's/^ID=(.*)/\1/p' /etc/os-release)" in
			arch|archarm)     model="\033[94m \033[0m$model" ;;
			gentoo)   	  model="\033[95m \033[0m$model" ;;
			debian)   	  model="\033[91m \033[0m$model" ;;
			ubuntu)   	  model="\033[93m \033[0m$model" ;;
			alpine)   	  model="\033[90m \033[0m$model" ;;
			fedora)   	  model="\033[90m \033[0m$model" ;;
			manjaro)  	  model="\033[92m \033[0m$model" ;;
			raspbian) 	  model="\033[91m \033[0m$model" ;;
			*)	  	  model="\033[97m \033[0m$model" ;;
		esac

		uname -a | grep -iq microsoft && model="\033[96m \033[0m$model"
	;;
	*BSD)
		which doas &> /dev/null && elevate=doas || elevate=sudo
		model="$($elevate dmidecode --type system 2> /dev/null | sed -nE 's/.*Product Name: (.*)/\1/p')"

		case $u in
			FreeBSD) model="\033[91m \033[0m $model" ;;
			OpenBSD) model="\033[93m \033[0m $model" ;;
			NetBSD)  model="\033[93m \033[0m $model" ;;
		esac
	;;
esac
printf "$model $(uname -rms)"
