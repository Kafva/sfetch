package lib

import (
	"fmt"
	"context"
	"time"
	"os"
	"os/exec"
	"strings"
)

/// If a connection has failed once we don't want to re-run it
/// Prevent a new go routine from running if one is already waiting on output
func GetUnameMapping(hosts_map map[string][]string) (uname_mapping map[string]string) {

	// Each SSH session is given its own channel inside a map to enable concurrent execution
	uname_mapping = make(map[string]string, len(hosts_map))
	chan_mapping := make(map[string]chan string)
    
	for host, jump_hosts := range hosts_map {
		Debug("=> (Host)", host)
		addUnameMapping(uname_mapping, chan_mapping, host)

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
			Debug("=> (Jump)", jump_to)
			addUnameMapping(uname_mapping, chan_mapping, jump_to)
		}
	}
	
	if !*SLOW {
		// Once every SSH operation has started begin waiting for each one to complete
		for host,ch := range chan_mapping {
			Debug("=> (Waiting on)", host)
			uname_mapping[host] = <- ch
		}
	}

	return uname_mapping
}

func addUnameMapping(uname_mapping map[string]string, chan_mapping map[string]chan string, host string) {

	if uname_mapping[host] == "" { 
		// Only fetch `uname` output if none exists in the map 
		// if a previous command has been executed and failed the mapping wont be empty
		
		if !*SLOW {
			if chan_mapping[host] == nil { 
				chan_mapping[host] = make(chan string)
			}
			go GetHostInfoChannel(host, chan_mapping[host])
		} else {
			uname_mapping[host] = GetHostInfo(host)
		}
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

func GetHostInfoChannel(host string, info chan string) {
	info <- GetHostInfo(host)	
}

/// Returns information regarding the provided host using `ssh` (if not "localhost") and `uname`
/// if a VERBOSE level >0 is provided a custom script is passed as stdin to the process
/// instead of running uname 
func GetHostInfo(host string) string {
	
	// If the host requires more than one proxy jump and can't be reached the process will hang
	// we therefore need a maximum execution time after which a connection is considered failed
	// Using `-o ConnectTimeout` was insufficient
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*CONNECTION_TIMEOUT) )
	defer cancel()

	cmd := exec.Cmd{}
	
	if host != LOCALHOST {
		cmd = *exec.Command(
			SSH_PATH,
			"-F",
			*CONFIG_FILE,
			host,
		)
	}
	
	script := ""

	switch {
		case *VERBOSE == 1:
			script = INFO_SCRIPT
		case *VERBOSE >= 2:
			script = FULL_INFO_SCRIPT
		default:
			if host == LOCALHOST {
				cmd = *exec.Command("uname", "-rms")
			} else {
				cmd.Args = append(cmd.Args, "uname", "-rms") 
			}
	}

	if script != "" {
		if host == LOCALHOST {
			if RELEASE {
				cmd = *exec.Command("/bin/sh", "-c", script)
			} else {
				cmd = *exec.Command(script) 
			}
		} else {
			if RELEASE {
				f := strings.NewReader(script)
				cmd.Stdin = f 
			} else {
				f, err := os.Open(script)
				if err != nil { 
					Die("Missing: ", script) 
				}
				defer f.Close()
				cmd.Stdin = f 
			}
		}
	}

	result, err := cmd.Output()
	
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "[%s] Connection timeout: %s\n", host, err.Error())
		} else {
			fmt.Fprintf(os.Stderr, "[%s] Command failed: %s\n", host, err.Error())
		}
		return FAILED
	}

	return strings.TrimSuffix(string(result), "\n") 
}
