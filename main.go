package main

import (
	"flag"
	"fmt"
	"os"
	//"os/exec"
	"github.com/Kafva/sfetch/lib"
    //"github.com/disiqueira/gotree"
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
	
	ignore_hosts 	:= lib.GetIgnoreHosts(*ignore_file)
	fmt.Println("Ignore:", ignore_hosts)

	hosts_map 		:= lib.GetHostMapping(*config_file, ignore_hosts)
	fmt.Println("SSH: ", hosts_map)

	//uname_mapping 	:= lib.GetUnameMapping(hosts_map, *config_file)
	fmt.Println("-------------------------")

	//uname_mapping := map[string]string {
	//	"club": "Linux 5.13.9-arch1-1 x86_64", 
	//	"devi": "FreeBSD 13.0-RELEASE-p3 amd64", 
	//	"kafva.one": "Linux 5.11.4-1-ARCH aarch64", 
	//	"vel": "Linux 5.13.9-arch1-1 x86_64",
	//}
	
	//model_cmd  := lib.GetModelCommand("Darwin")
	//ssh kafva.one 'bash -c "uname -rms"'
	//THIS:https://stackoverflow.com/questions/39496572/running-command-with-pipe-in-golang-exec
	//model, _ := exec.Command( model_cmd.Path, model_cmd.Args... ).Output()
	
	//root := gotree.New("localhost")
	//sub := root.Add(uname_mapping["kafva.one"])
	//sub.Add(uname_mapping["vel"])
	//sub.Add(uname_mapping["club"])
	//fmt.Println(root.Print())

	//fmt.Println(model, model_cmd, uname_mapping)
}
