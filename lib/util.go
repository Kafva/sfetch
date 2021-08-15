package lib

import (
	"os"
	"fmt"
	"github.com/disiqueira/gotree"
	flag "github.com/spf13/pflag"
)

func Die(strs ... interface{}){
	// Go does not support optional parameters (code=1)
	strs = append(strs, "\n")
	fmt.Fprint(os.Stderr, strs ...)
	os.Exit(1)
}

func Debug(strs ... interface{}){
	if DEBUG {
		fmt.Println(strs ...)
	}
}

func DetailUsage(){
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("More info...")
}

func MakeTree(uname_mapping map[string]string, hosts_map map[string][]string, config_file string, verbosity int) {
	localhost := GetHostInfo("localhost", config_file, verbosity) 

	root := gotree.New(localhost)
	//sub := root.Add(uname_mapping["kafva.one"])
	//sub.Add(uname_mapping["vel"])
	//sub.Add(uname_mapping["club"])
	fmt.Println(root.Print())

}