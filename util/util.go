package util

import "os"

func Die(msg string){
	// Go does not support optional parameters (code=1)
	println(msg)
	os.Exit(1)
}
