package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kevinburke/ssh_config"
)

// 1. Enumerate all hosts in the SSH config and store them in a struct list
// 2. SSH into each one in parallell and save `uname` output
// 3. Display in tree format
type SSHHost struct {
	Hostname string
	Uname string
	JumpHosts[] SSHHost
}


func MustGet(cfg ssh_config.Config, host string, key string) string {
	
	value, err := cfg.Get(host, key)

	if err != nil { 
		die(err.Error()) 
	} else if value == "" {
		die(fmt.Sprintf("Missing key for %s: %s", host, key))
	}
	
	return value
}

func die(msg string){
	// Go does not support optional parameters (code=1)
	println(msg)
	os.Exit(1)
}

func DetailUsage(){
	fmt.Printf("Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("More info...")
}

func main() {
	home, _ := os.UserHomeDir()
	default_ssh_config := fmt.Sprintf("%s/.ssh/config", home)	

	configFile := flag.String(
		"config", 
		fmt.Sprintf("%s/.config/sfetch/config.json", home), 
		"Path to an optional JSON config file",
	)
	
	sshConfigPath := flag.String(
		"ssh_config",
		default_ssh_config, 
		"Path to ssh config",
	)
	
	flag.Usage = DetailUsage
	flag.Parse()
	
	// Work around to ignore errors
	//	https://stackoverflow.com/questions/28227095/multiple-values-in-single-value-context

	f, err := os.Open(*sshConfigPath) 
	if err != nil { die(err.Error()) }

	sshConfig, err := ssh_config.Decode(f)
	if err != nil { die(err.Error()) }

	f.Close()

	fmt.Printf("Port: %s\n", MustGet(*sshConfig, "kafva.one", "Port")  )
	fmt.Printf("config: %s\n", *configFile)

}
