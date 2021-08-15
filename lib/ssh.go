package lib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetUnameMapping(tree_map map[string][]string, config_file string, verbosity int) (uname_mapping map[string]string) {

	uname_mapping = make(map[string]string, len(tree_map))

	for host, jump_hosts := range tree_map {

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

/// Creates a mapping where key == value, enables tree printing without any network connectivity
func GetHostnameMapping(tree_map map[string][]string) (uname_mapping map[string]string) {
	uname_mapping = make(map[string]string, len(tree_map))

	for host, jump_hosts := range tree_map {

		if uname_mapping[host] == "" {
			uname_mapping[host] = host
		}

		for _, jump_to := range jump_hosts {
			if uname_mapping[jump_to] == "" {
				uname_mapping[jump_to] = jump_to 
			}
		}
	}

	return uname_mapping
}

/// Returns information regarding the provided host using `ssh` (if not "localhost") and `uname`
/// if a verbosity level >0 is provided a custom script is passed as stdin to the process
/// instead of running uname 
func GetHostInfo(host, config_file string, verbosity int) string {
	
	Debug("=>", host)
	
	cmd := exec.Cmd{}
	
	if host != "localhost" {
		cmd = *exec.Command(
			SSH_PATH,
			"-F",
			config_file,
			"-o",
			fmt.Sprintf("ConnectTimeout=%d", CONNECTION_TIMEOUT),
			host,
		)
	}
	
	script_path := ""

	switch {
		case verbosity == 1:
			script_path = INFO_SCRIPT
		case verbosity >= 2:
			script_path = FULL_INFO_SCRIPT
		default:
			if host == "localhost" {
				cmd = *exec.Command("uname", "-rms")
			} else {
				cmd.Args = append(cmd.Args, "uname", "-rms") 
			}
	}

	if script_path != "" {
		if host == "localhost" {
			cmd = *exec.Command(script_path) 
		} else {
			f, err := os.Open(script_path)
			if err != nil { 
				Die("Missing: ", script_path) 
			}
			defer f.Close()
			cmd.Stdin = f 
		}
	}

	result, err := cmd.Output()
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "[%s] Command failed: %s\n", host, err.Error())
		return ""
	}

	return strings.TrimSuffix(string(result), "\n") 
}