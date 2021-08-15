package lib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetUnameMapping(hosts_map map[string][]string, config_file string, verbosity int) (uname_mapping map[string]string) {

	uname_mapping = make(map[string]string, len(hosts_map))

	for host, jump_hosts := range hosts_map {

		if uname_mapping[host] == "" {
			// If their uname_mapping is empty, fetch the `uname` output
			uname_mapping[host] = GetHostInfo(host, config_file, verbosity)
		}

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
			if uname_mapping[jump_to] == "" {
				uname_mapping[jump_to] = GetHostInfo(jump_to, config_file, verbosity)
			}
		}
	}

	return uname_mapping
}

func GetHostInfo(host, config_file string, verbosity int) string {
	
	Debug("=>", host)
	
	cmd := exec.Cmd{}
	
	if host == "localhost" {
		cmd = exec.Cmd {
			Path: "/bin/sh",
			Args: []string {},
		}
	} else {
		cmd = exec.Cmd {
			Path: SSH_PATH,
			Args: []string {
				"-F",
				config_file,
				"-o",
				fmt.Sprintf("ConnectTimeout=%d", CONNECTION_TIMEOUT),
				host,
			},
		}
	}
	
	script_path := ""

	switch {
		case verbosity == 1:
			script_path = INFO_SCRIPT
		case verbosity >= 2:
			script_path = FULL_INFO_SCRIPT
		default:
			if host == "localhost" {
				cmd.Path = "uname"
				cmd.Args = append(cmd.Args, "-rms") 	
			} else {
				cmd.Args = append(cmd.Args, "uname", "-rms") 
			}
	}

	if script_path != "" {
		if host == "localhost" {
			cmd.Path = script_path 
		} else {
			f, err := os.Open(script_path)
			if err != nil { 
				Die("Missing: ", script_path) 
			}
			defer f.Close()
			cmd.Stdin = f 
			cmd.Args = append(cmd.Args, "--")
		}
	}

	Debug(cmd)
	result, err := cmd.CombinedOutput()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s] Command failed: %s\n", host, err.Error())
		return ""
	}

	return strings.TrimSuffix(string(result), "\n") 
}