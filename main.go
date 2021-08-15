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
///		host -> [uname] 
/// 3. Use the two mappings to create a tree-form display
func main() {
	home, _ := os.UserHomeDir()

	help := flag.BoolP("help", "h", false, "Show this help message and exit")

	verbosity := flag.CountP("verbose", "v", "Increase verbosity")
	
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
	
	if *help {
		DetailUsage()
		os.Exit(1)
	} 
	
	SSH_PATH, _ = exec.LookPath("ssh")
	if SSH_PATH == "" { 
		Die("No ssh executable found"); 
	}
	
	Debug("verbosity:", *verbosity)
	
	ignore_hosts 	:= GetIgnoreHosts(*ignore_file)
	Debug("ignore_hosts:", ignore_hosts)

	hosts_map 		:= GetHostMapping(*config_file, ignore_hosts)
	Debug("hosts_map:", hosts_map)

	uname_mapping 	:= GetUnameMapping(hosts_map, *config_file, *verbosity)
	//uname_mapping := map[string]string {
	//	"club": "Linux 5.13.9-arch1-1 x86_64", 
	//	"devi": "FreeBSD 13.0-RELEASE-p3 amd64", 
	//	"kafva.one": "Linux 5.11.4-1-ARCH aarch64", 
	//	"vel": "Linux 5.13.9-arch1-1 x86_64",
	//}
	
	Debug("-------------------------")
	Debug(uname_mapping)
	
	MakeTree(uname_mapping, hosts_map, *config_file, *verbosity)
}
