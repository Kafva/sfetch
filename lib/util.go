package lib

import (
	"fmt"
	"os"
	tree "github.com/disiqueira/gotree"
	flag "github.com/spf13/pflag"
)

func Die(strs ... interface{}) {
	strs = append(strs, "\n")
	fmt.Fprint(os.Stderr, strs ...)
	os.Exit(EXIT_ERROR)
}

func Debug(strs ... interface{}) {
	if *DEBUG {
		fmt.Println(strs ...)
	}
}

func ErrMsg(format string, args ... interface{}) {
	if !*QUIET {
		fmt.Fprintf(os.Stderr, format, args ...)
	}
}

func DetailUsage(){
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	if RELEASE {
		fmt.Println(HELP_STR + "\n(Release version)")
	} else {
		fmt.Println(HELP_STR + "\n(Development version)")
	}
}

func addToTree(uname_mapping map[string]string, hosts_map map[string][]string, root tree.Tree, hostname string) {
	
	uname := uname_mapping[hostname]
	if uname != "" && uname != COMMAND_FAILED && uname != COMMAND_TIMEOUT {
		var current_node tree.Tree
		
		// If an error occured for a host its uname will be empty or "FAIL"
		if *INCLUDE_HOSTNAME {
			current_node = root.Add(uname + " " + HOSTNAME_ANSI_COLOR + hostname + "\033[0m")
		} else {
			current_node = root.Add(uname)
		}

		for _,host := range hosts_map[hostname] {
			// Traverse down all jumphosts
			addToTree(uname_mapping, hosts_map, current_node, host)	
		}
	}
}

func MakeTree(root_name string, uname_mapping map[string]string, hosts_map map[string][]string, has_jump map[string]struct{}) {
	
	root := tree.New(root_name)
	
	for hostname := range hosts_map {

		if _,found := has_jump[hostname]; found { 
			// Since the hosts_map is flat we need to ensure that we only iterate over the hosts
			// that are on the top level, i.e. those that do NOT have any proxies
			// Hosts with a proxy will appear implictly during the recursive calls
			continue
		}
		
		addToTree(uname_mapping, hosts_map, root, hostname)
	}

	fmt.Println(root.Print())
}