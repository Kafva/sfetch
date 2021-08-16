#!/bin/sh 
die(){ echo -e "$1" >&2 ; exit 1; }
minifyBash(){
	# Also escapes `"` andd `\`
	sed -E '/^#/d; s/\[/; \[/g' $1 | tr '\n' ' ' | sed 's/\s+/ /g;' | sed -E 's@\\@\\\\\\\\@g;  s/"/\\\\"/g'
}

if ! $(sed --version 2>&1 | grep -q "GNU sed"); then
	which gsed &> /dev/null || 
		die 'Install GNU `sed`' && 
		sedExec=gsed
fi

project_dir="${PWD%%sfetch*}sfetch"
info_script=$(minifyBash $project_dir/scripts/info.sh)
full_info_script=$(minifyBash $project_dir/scripts/full_info.sh)
config="$project_dir/lib/config.go"


cp $config /tmp/config.go
$sedExec -i "s/RELEASE = false/RELEASE = true/; s@./scripts/info.sh@${info_script}@; s@./scripts/full_info.sh@${full_info_script}@;" $config

go build && go install

cp /tmp/config.go $config
