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
	fmt.Println("More info...")
}

func addToTree(uname_mapping map[string]string, tree_map map[string][]string, root tree.Tree, hostname string) {
	
	uname 	   	 := uname_mapping[hostname]
	if uname != "" {
		// If an error occured for a host its uname will be empty
		current_node := root.Add(uname)

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
			continue
		}
		
		addToTree(uname_mapping, tree_map, root, hostname)
	}

	fmt.Println(root.Print())
}