package lib

import (
	"fmt"
	"os/exec"
	"strings"
)

func GetUnameMapping(hosts_map map[string][]string, config_file string) (uname_mapping map[string]string) {

	ssh_path, err := exec.LookPath("ssh")
	if err != nil { 
		Die("No ssh executable found"); 
	}

	uname_mapping = make(map[string]string, len(hosts_map))

	for host, jump_hosts := range hosts_map {

		if uname_mapping[host] == "" {
			// If their uname_mapping is empty, fetch the `uname` output
			uname_mapping[host] = GetUname(host, config_file, ssh_path)
		}

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
			if uname_mapping[jump_to] == "" {
				uname_mapping[jump_to] = GetUname(jump_to, config_file, ssh_path)
			}
		}
	}

	return uname_mapping
}


func GetModelCommand(os_type string) string {
	cmd := ""

	switch os_type {
	case "Darwin":
		cmd = "system_profiler SPHardwareDataType | sed -nE 's/.*Model Identifier: (.*)/\\1/p'"
	case "Linux":
		cmd = "cat /sys/devices/virtual/dmi/id/board_{name,version} | tr '\n' ' '"
	case "FreeBSD":
		cmd = "doas dmidecode --type system | sed -nE 's/.*Product Name: (.*)/\\1/p'"
	default:
		Die(fmt.Sprintf("Unknown OS: %s", os_type))
	}
	
	return cmd
}

/// TODO error handling
func GetUname(host, config_file, ssh_path string) string {
	// Linux 5.11.4-1-ARCH aarch64
	// Linux 5.13.9-arch1-1 x86_64
	// Darwin 20.5.0 x86_64
	

	//fmt.Printf("=> %s\n", host)
	//return host

	uname_cmd := exec.Cmd {
		Path: ssh_path,
		Args: []string {
			"-F",
			config_file,
			host,
			"uname",
			"-rms",
		},
	}

	uname, err := exec.Command(uname_cmd.Path, uname_cmd.Args ...).Output()
	if err != nil {
		Die(err.Error())
	}

	return  strings.TrimSuffix(string(uname), "\n") 
}