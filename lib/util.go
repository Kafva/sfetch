package lib

import (
	"os"
	"fmt"
	tree "github.com/disiqueira/gotree"
	flag "github.com/spf13/pflag"
)

func Die(strs ... interface{}){
	// Go does not support optional parameters (code=1)
	strs = append(strs, "\n")
	fmt.Fprint(os.Stderr, strs ...)
	os.Exit(1)
}

func Debug(strs ... interface{}){
	if *DEBUG {
		fmt.Println(strs ...)
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

func addToTree(uname_mapping map[string]string, tree_map map[string][]string, root tree.Tree, hostname string) {
	
	uname := uname_mapping[hostname]
	if uname != "" && uname != FAILED {
		var current_node tree.Tree
		
		// If an error occured for a host its uname will be empty or "FAIL"
		if *INCLUDE_HOSTNAME {
			current_node = root.Add(uname + " " + HOSTNAME_ANSI_COLOR + hostname + "\033[0m")
		} else {
			current_node = root.Add(uname)
		}

		for _,host := range tree_map[hostname] {
			// Traverse down all jumphosts
			addToTree(uname_mapping, tree_map, current_node, host)	
		}
	}
}

func MakeTree(root_name string, uname_mapping map[string]string, tree_map map[string][]string, has_jump map[string]struct{}) {
	
	root := tree.New(root_name)
	
	for hostname := range tree_map {

		if _,found := has_jump[hostname]; found { 
			// Since the tree_map is flat we need to ensure that we only iterate over the hosts
			// that are on the top level, i.e. those that do NOT have any proxies
			// Other hosts will appear implictly during the recursive calls
			continue
		}
		
		addToTree(uname_mapping, tree_map, root, hostname)
	}

	fmt.Println(root.Print())
}