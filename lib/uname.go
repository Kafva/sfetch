package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

/// If a connection has failed once we don't want to re-run it
/// Prevent a new go routine from running if one is already waiting on output
func GetUnameMapping(hosts_map map[string][]string) (uname_mapping map[string]string) {

	// Each SSH session is given its own channel inside a map to enable concurrent execution
	uname_mapping = make(map[string]string, len(hosts_map))
	chan_mapping := make(map[string]chan string)
	
	// To avoid several go routines from being launched for the same host
	// we maintain a map of all hosts which have been / are being processed
	// Each host should only have one `=>` message in debug mode
	invoked_mapping := make(map[string]struct{}, len(hosts_map))
    
	for host, jump_hosts := range hosts_map {

		addUnameMapping(uname_mapping, chan_mapping, invoked_mapping, host)

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
			
			addUnameMapping(uname_mapping, chan_mapping, invoked_mapping, jump_to)
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

func addUnameMapping(uname_mapping map[string]string, chan_mapping map[string]chan string, invoked_mapping map[string]struct{}, host string) {
	
	if _, found := invoked_mapping[host]; !found {
		// Ensure that another go-routine hasn't been ran / is running for the host
		Debug("=>", host)
		invoked_mapping[host] = struct{}{}

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

	cmd := exec.Cmd{}
	
	if host != LOCALHOST {
		cmd = *exec.Command(SSH_PATH, "-F", *CONFIG_FILE, 
			"-o", fmt.Sprintf("ConnectTimeout=%d", *CONNECTION_TIMEOUT), 
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
	
	var out bytes.Buffer
	timeout_regex := regexp.MustCompile("exit status 255")

	// Using the .Output() functions hangs on hosts were the jump hosts is accessible
	// but the target is unavailable 
	cmd.Stdout = &out
	err := cmd.Run()
	
	if err != nil {
		if timeout_regex.Match([]byte(err.Error())) {
			ErrMsg("[%s] Connection failed: %s\n", host, err.Error())
			return COMMAND_TIMEOUT
		} else {
			ErrMsg("[%s] Command failed: %s\n", host, err.Error())
			return COMMAND_FAILED
		}
	}
	
	result, err := io.ReadAll(&out)
	if err != nil {
		ErrMsg("[%s] Read error: %s\n", host, err.Error())
	}

	return strings.TrimSuffix(string(result), "\n") 
}
