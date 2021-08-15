package lib

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

/// If a connection has failed once we don't want to re-run it
/// Prevent a new go routine from running if one is already waiting on output
func GetUnameMapping(hosts_map map[string][]string, config_file string, VERBOSE int) (uname_mapping map[string]string) {

	// Each SSH session is given its own channel inside a map to enable concurrent execution
	uname_mapping = make(map[string]string, len(hosts_map))
	chan_mapping := make(map[string]chan string)
    
	// TODO
	var wg sync.WaitGroup
	wg.Add(len(hosts_map))

	for host, jump_hosts := range hosts_map {
		addUnameMapping(uname_mapping, chan_mapping, &wg, host, config_file, VERBOSE)

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
			addUnameMapping(uname_mapping, chan_mapping, &wg, jump_to, config_file, VERBOSE)
		}
	}
	
	Debug("Wait starting...")
	wg.Wait()

	return uname_mapping
}

func addUnameMapping(uname_mapping map[string]string, chan_mapping map[string]chan string, wg *sync.WaitGroup, host, config_file string, VERBOSE int) {

	if uname_mapping[host] == "" { 
		// Only fetch `uname` output if none exists in the map 
		// if a previous command has been executed and failed the mapping wont be empty
		
		if chan_mapping[host] == nil { 
			chan_mapping[host] = make(chan string)
		}
		go GetHostInfo(host, config_file, VERBOSE, chan_mapping[host], wg)
		uname_mapping[host] = <- chan_mapping[host] 
		
		//uname_mapping[host] = GetHostInfoSerial(host, config_file, VERBOSE)
	}
}

/// Creates a mapping where key == value, enables tree printing without any network connectivity
func GetHostnameMapping(hosts_map map[string][]string) (uname_mapping map[string]string) {
	uname_mapping = make(map[string]string, len(hosts_map))

	for host, jump_hosts := range hosts_map {

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
/// if a VERBOSE level >0 is provided a custom script is passed as stdin to the process
/// instead of running uname 
func GetHostInfo(host, config_file string, VERBOSE int, info chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	Debug("=>", host)
	
	cmd := exec.Cmd{}
	
	if host != "localhost" {
		cmd = *exec.Command(
			SSH_PATH,
			"-F",
			config_file,
			"-o",
			fmt.Sprintf("ConnectTimeout=%d", *CONNECTION_TIMEOUT),
			host,
		)
	}
	
	script_path := ""

	switch {
		case VERBOSE == 1:
			script_path = INFO_SCRIPT
		case VERBOSE >= 2:
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
		info <- FAILED
	}

	info <- strings.TrimSuffix(string(result), "\n") 
}


/// Returns information regarding the provided host using `ssh` (if not "localhost") and `uname`
/// if a VERBOSE level >0 is provided a custom script is passed as stdin to the process
/// instead of running uname 
func GetHostInfoSerial(host, config_file string, VERBOSE int) string {
	
	Debug("=>", host)
	
	cmd := exec.Cmd{}
	
	if host != "localhost" {
		cmd = *exec.Command(
			SSH_PATH,
			"-F",
			config_file,
			"-o",
			fmt.Sprintf("ConnectTimeout=%d", *CONNECTION_TIMEOUT),
			host,
		)
	}
	
	script_path := ""

	switch {
		case VERBOSE == 1:
			script_path = INFO_SCRIPT
		case VERBOSE >= 2:
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
		return FAILED
	}

	return strings.TrimSuffix(string(result), "\n") 
}
