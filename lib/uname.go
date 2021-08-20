package lib

import (
	"bytes"
	"fmt"
	"bufio"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"runtime"
)

// If a connection has failed once we don't want to re-run it
// Prevent a new go routine from running if one is already waiting on output
func GetUnameMapping(hosts_map map[string][]string) (uname_mapping map[string]string) {

	// Each SSH session is given its own channel inside a map to enable concurrent execution
	uname_mapping = make(map[string]string, len(hosts_map))
	chan_mapping := make(map[string]chan string)
	 
	for host, jump_hosts := range hosts_map {
		addUnameMapping(uname_mapping, chan_mapping, host)

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
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
		// Ensure that another go-routine hasn't been ran / is running for the host
		Debug("=>", host)

		// To avoid several go routines from being launched for the same host
		// we set each host entry as 'IN_PROGRESS' before launching a go routine 
		// Each host should only have one `=>` message in debug mode
		uname_mapping[host] = COMMAND_IN_PROGRESS

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

// Creates a mapping where key == value, enables tree printing without any network connectivity
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

// Returns information regarding the provided host using `ssh` (if not "localhost") and `uname`
// if a VERBOSE level >0 is provided a custom script is passed as stdin to the process
// instead of running uname 
func GetHostInfo(host string) string {

	cmd := exec.Cmd{}
	
	if host != LOCALHOST {
		cmd = *exec.Command(SSH_PATH, "-F", *CONFIG_FILE, 
			"-o", fmt.Sprintf("ConnectTimeout=%d", *CONNECTION_TIMEOUT), 
			host,
		)
	}
	
	script := ""
	if runtime.GOOS == "windows" {

	} else {
	}

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

	// In the release build the INFO_SCRIPT values will
	// contain the actual code to be ran and in dev mode
	// it will be the path to the script 
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
	
	// Using the .Output() functions hangs on hosts were the jump host is accessible
	// but the target is unavailable, .Run() is therefore used instead 
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	
	if err != nil {
		timeout_regex := regexp.MustCompile("exit status 255")
		
		if timeout_regex.Match([]byte(err.Error())) {
			ErrMsg("[%s] Connection failed: %s\n", host, err.Error())
			return COMMAND_TIMEOUT
		} else {
			// Retry using a Windows compatible command
			return GetWindowsHostInfo(host)
		}
	}
	
	result, err := io.ReadAll(&out)
	if err != nil {
		ErrMsg("[%s] Read error: %s\n", host, err.Error())
	}
	
	if len(result) > 200 {
		// In Release mode on Windows the captured output becomes the entire
		// script (i.e. over 200 chars)
		return GetWindowsHostInfo(host)
	}

	ret := strings.TrimSuffix(string(result), "\n")
	if ret == "" {
		ret = COMMAND_FAILED
	}
	return ret  
}

func GetWindowsHostInfo(host string) string {
	
	cmd := exec.Cmd{}

	if host != LOCALHOST {
		cmd = *exec.Command(SSH_PATH, "-F", *CONFIG_FILE, 
			"-o", fmt.Sprintf("ConnectTimeout=%d", *CONNECTION_TIMEOUT), 
			host,
		)
	} else {
		cmd = *exec.Command("powershell.exe", "-c")
	}

	prefix := ""

	switch  {
	case *VERBOSE >= 2:
		prefix = WINDOWS_PREFIX 
		cmd.Stdin = strings.NewReader(WINDOWS_FULL_INFO)
	case *VERBOSE >= 1:
		cmd.Stdin = strings.NewReader(WINDOWS_FULL_INFO)
	default:
		cmd.Stdin = strings.NewReader(WINDOWS_INFO) 
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	
	if err != nil {
		ErrMsg("[%s] Command failed: %s\n", host, err.Error())
		return COMMAND_FAILED
	}
	
	result, err := io.ReadAll(&out)
	if err != nil {
		ErrMsg("[%s] Read error: %s\n", host, err.Error())
	}

	return prefix + ParseSystemInfo(string(result))
}

/// Return a string comprised of `(Motherboard) - OS name - OS Version` given output from wmic commands
func ParseSystemInfo(systemInfo string) string {
	scanner := bufio.NewScanner( strings.NewReader(systemInfo) )
	
	header_regex := regexp.MustCompile(`(?i)^\s*Name|Version|Product\s*`)
	output_regex := regexp.MustCompile(`(?i)[-_a-zA-Z0-9.() ]+`)
	spaces_regex := regexp.MustCompile(`\s+`)
	
	output_line := false
	uname := ""

	for scanner.Scan() {
		line := scanner.Text()
		
		if output_line {
			uname = uname + " " + output_regex.FindString(line)	
			output_line = false
		} else {
			if header_regex.Match([]byte(line)) { 
				// The next line will hold a value we want to save
				output_line = true
			}
		}
	}
	
	if uname == "" {
		return COMMAND_FAILED
	} else {
		// Remove all trailing whitespaces and replace sequences of more than one
		// space with a single space
		return strings.TrimSpace( spaces_regex.ReplaceAllString(uname, " "))
	}
}