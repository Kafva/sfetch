#!/bin/sh 
die(){ echo -e "$1" >&2 ; exit 1; }

sedExec=sed
if ! $(sed --version 2>&1 | grep -q "GNU sed"); then
	which gsed &> /dev/null || 
		die 'Install GNU `sed`' && 
		sedExec=gsed
fi

minifyBash(){
	$sedExec -E '/^#/d; s/\[ /; \[ /g' $1 | tr '\n' ' ' | 
		$sedExec 's/\s+/ /g;' 	| # Remove unnecessary whitespace
		$sedExec -E 's@\\@\\\\@g;' 	| # Escape backslashes
		$sedExec 's/\&/\\\&/g' 	| # `&` needs to be escaped to avoid issues in the `sed -i` statement
		$sedExec 's/case/;case/g'	| # Hacks for case...esac
		$sedExec 's/^;case/case/' |
		$sedExec 's/esac\s*printf/esac ; printf/' 

}


project_dir="${PWD%%sfetch*}sfetch"
info_script=$(minifyBash $project_dir/scripts/info.sh)
full_info_script=$(minifyBash $project_dir/scripts/full_info.sh)
config="$project_dir/lib/config.go"

cp $config /tmp/config.go

$sedExec -i "s/RELEASE = false/RELEASE = true/; s@./scripts/info.sh@${info_script}@; s@./scripts/full_info.sh@${full_info_script}@" $config

goExec=go
uname -a | grep -q microsoft && goExec=/mnt/c/Users/$USER/scoop/shims/go.exe

$goExec build && $goExec install
cp /tmp/config.go $config
