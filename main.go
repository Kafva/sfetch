package main

import (
	flag "github.com/spf13/pflag"
	"fmt"
	"os"
	"os/exec"
	. "github.com/Kafva/sfetch/lib"
)

/// 1. Enumerate all hosts in the SSH config in a map on the form:
///		host -> [jump_to_hosts]
/// 2. SSH into each one in parallell and construct another map
///		host -> uname 
/// 3. Use the two mappings to create a tree-form display
func main() {
	home, _ := os.UserHomeDir()

	HELP 				:= flag.BoolP("help", "h", false, "Show this help message and exit")
	VERBOSE 			:= flag.CountP("verbose", "v", "Increase verbosity")
	DEBUG 				=  flag.BoolP("debug", "d", false, "Print debug information")
	CONNECTION_TIMEOUT 	=  flag.IntP("timeout", "t", 2, "Connection timeout for ssh sessions")
	
	BASIC := flag.BoolP(
		"basic",
		"b",
		false,
		"Print tree structure without connecting to any hosts for OS information",
	)
	
	config_file := flag.StringP(
		"ssh_config",
		"c",
		fmt.Sprintf("%s/.ssh/config", home), 
		"Path to ssh config",
	)
	
	ignore_file := flag.StringP(
		"ignore",
		"i",
		"",
		"Path to a file with hosts to ignore",
	)
	
	flag.Usage = DetailUsage
	flag.Parse()
	
	if *HELP {
		DetailUsage()
		os.Exit(1)
	} 
	
	SSH_PATH, _ = exec.LookPath("ssh")
	if SSH_PATH == "" { 
		Die("No ssh executable found"); 
	}
	
	Debug("VERBOSE:", *VERBOSE)
	Debug("BASIC:", *BASIC)
	
	ignore_hosts 	:= GetIgnoreHosts(*ignore_file)
	Debug("ignore_hosts:", ignore_hosts)
	
	hosts_map, has_jump 	:= GetHostMapping(*config_file, ignore_hosts, false)
	Debug("hosts_map:", hosts_map)
	Debug("has_jump:", has_jump)

	var uname_mapping map[string]string
	var root_name string

	if !*BASIC {
		root_name 		= GetHostInfoSerial("localhost", *config_file, *VERBOSE)
		uname_mapping 	= GetUnameMapping(hosts_map, *config_file, *VERBOSE)
	} else {
		root_name,_ 	= os.Hostname()
		uname_mapping	= GetHostnameMapping(hosts_map)
	}
	
	Debug(uname_mapping)
	Debug("-------------------------")
	
	
	MakeTree(root_name, uname_mapping, hosts_map, has_jump)
}
