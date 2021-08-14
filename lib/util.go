package lib

import (
	"os"
	"fmt"
	"flag"
)

func Die(msg string){
	// Go does not support optional parameters (code=1)
	println(msg)
	os.Exit(1)
}

func DetailUsage(){
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("More info...")
}
