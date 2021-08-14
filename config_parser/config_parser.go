package config_parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Kafva/sfetch/util"
)

/// Returns a mapping on the form `host -> []jumpHosts` for each host
/// in the provided ssh_config
/// NOTE: we assume that each specified proxy has a corresponding 'Host' entry!
func ParseSshConfig(filepath string) (host_map map[string][]string) {
	// ssh_config format has:
	// 	Host dst
	//			ProxyJump proxy[,proxy2...]
	// Internally we want the opposite map:
	// 	ProxyHost proxy
	//			Hosts [dst]
	// Go through line by line and insert hosts into a map `host -> [jumpHosts]` 
	// If we were to use an array  we would haft to perform a linear search
	// each time we find a ProxyJump to save it in the correct position


	// The 'defer' keyword will push the given statement onto
	// a stack which is executed on exit of the current function,
	// we can thus ensure that the opened file will be closed
	// regardless of what happens
	f, err := os.Open(filepath) 
	if err != nil { util.Die(err.Error()) }
	defer f.Close()

	scanner 	:= bufio.NewScanner(f)
	
	// (?i) provides case-Insensitive matching
	host_regex 			:= regexp.MustCompile(`(?i)^\s*Host\s+([^\s]*)`)
	proxyjump_regex 	:= regexp.MustCompile(`(?i)^\s*ProxyJump\s+([^\s]*)`)
	
	// Maps hostname -> [jumpHosts]
	hosts_map 		 	:= make(map[string][]string)
	current_host := ""
	
	for scanner.Scan() {
		
		line := scanner.Text()

		// The first match will be the entire string and the 
		// remaining indices will hold capture groups i.e. the hostname 
		matches := host_regex.FindStringSubmatch(line)
		
		if len(matches) > 0 { 
			// Add the identified host to the map with an empty array as its value
			// (unless it is already present) and continue parsing the next line
			current_host = matches[1]
			if hosts_map[current_host] == nil {
				hosts_map[current_host] = make([]string, 0) 
			}
			continue
		}

		// Look for a ProxyJump if the line did not have a `Host` value
		matches = proxyjump_regex.FindStringSubmatch(line)

		if len(matches) > 0 { 
			jump_hosts 	  :=  strings.Split(matches[1], ",")

			if len(jump_hosts) == 1 {

				fmt.Println(current_host,jump_hosts)

				jump_host := jump_hosts[0]

				if hosts_map[jump_host] == nil {
					// If the jump host doesn't have a key, initalise the array
					hosts_map[jump_host] = []string { current_host } 
				} else {
					// Otherwise append the current_host
					hosts_map[jump_host] = append(hosts_map[jump_host], current_host) 
				}
			} else {
				// The `current_host` will be set to the most recently read `Host` line
				// and will be added to the exit node when more than one proxy exists
				//	hosts_map[jump_host_1]	+= jump_host_2
				//	hosts_map[jump_host_2]	+= jump_host_3
				//	hosts_map[jump_host_3]	+= end_host
				//	...
				main_host := jump_hosts[0]
				jump_hosts = append(jump_hosts[1:], current_host)
				
				fmt.Println(main_host, jump_hosts)

				for _, jump_host := range jump_hosts {

					if hosts_map[main_host] == nil {
						// If the jump host doesn't have a key, initalise the array
						hosts_map[main_host] = []string { jump_host } 
					} else {
						// Otherwise append the next_host
						hosts_map[main_host] = append(hosts_map[main_host], jump_host) 
					}

					main_host = jump_host
				}

			}

		}
	} 

	return hosts_map
}