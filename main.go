package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/Kafva/sfetch/lib"
)

/// 1. Enumerate all hosts in the SSH config in a map on the form:
///		host -> [jump_to_hosts]
/// 2. SSH into each one in parallell and construct another map
///		host -> [uname] 
/// 3. Use the two mappings to create a tree-form display
func main() {
	home, _ := os.UserHomeDir()

	config_file := flag.String(
		"ssh_config",
		fmt.Sprintf("%s/.ssh/config", home), 
		"Path to ssh config",
	)
	
	ignore_file := flag.String(
		"ignore",
		"",
		"Path to a file with hosts to ignore",
	)
	
	flag.Usage = lib.DetailUsage
	flag.Parse()
	
	hosts_map := lib.GetHostMapping(*config_file)

	fmt.Println("-------------------------")
	fmt.Println(hosts_map,"\n",ignore_file)
}
