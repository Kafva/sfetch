package lib

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

/// Returns a map with keys for all hosts to allow for an
/// easy way to determine if a host should be ignored
///		if _, ok := ignore_hosts["name"]; ok 
/// The values in the map are 'nil' since we only care about the keys
func GetIgnoreHosts(ignore_file string) map[string]struct{} {
	ignore_hosts := make(map[string]struct{})
	
	if ignore_file != "" {
		f, err := os.Open(ignore_file) 
		if err != nil { 
			Die(err.Error()) 
		}
		defer f.Close()

		scanner 	 := bufio.NewScanner(f)
		host_regex 	 := regexp.MustCompile(`^[^#]`)

		for scanner.Scan() {
			line := scanner.Text()
			if host_regex.Match( []byte(line) ) {

				ignore_hosts[strings.TrimSpace(line)] = struct{}{}
			}
		}
	}

	return ignore_hosts
}

/// Returns a mapping on the form `host -> []jumpHosts` for each host
/// in the provided ssh_config, excluding hosts provided in the `ignore_hosts` map
/// NOTE: we assume that each specified proxy has a corresponding 'Host' entry and
/// dont look at 'Hostname'
/// If `tree` is set to true dedicated entries for hosts that have a proxy wont be returned, i.e.
///		[ loc1:[], loc2:[], loc3:[loc1, loc2], loc4:[] ]
///	becomes
///		[ loc3: [loc1, loc2], loc4:[] ] 
func GetHostMapping(config_file string, ignore_hosts map[string]struct{}, tree bool) (hosts_map map[string][]string, has_jump map[string]struct{}) {
	// ssh_config format has:
	// 	Host dst
	//			ProxyJump proxy[,proxy2...]
	// Internally we want the opposite map so that we get a tree structure:
	// 	(Proxy)Host proxy
	//			Hosts [dst]

	f, err := os.Open(config_file) 
	if err != nil { 
		Die(err.Error()) 
	}
	defer f.Close()

	scanner 	:= bufio.NewScanner(f)
	
	// (?i) provides case-Insensitive matching
	host_regex 			:= regexp.MustCompile(`(?i)^\s*Host\s+([^\s]*)`)
	proxyjump_regex 	:= regexp.MustCompile(`(?i)^\s*ProxyJump\s+([^\s]*)`)
	
	// Maps hostname -> [jump_to_hosts]
	hosts_map 		 	= make(map[string][]string)
	// List of all hosts that appear after as a proxy
	has_jump = make(map[string]struct{})
	current_host := ""

	for scanner.Scan() {
		
		line := scanner.Text()

		// The first match will be the entire string and the 
		// remaining indices will hold capture groups i.e. the hostname 
		matches := host_regex.FindStringSubmatch(line)
		
		if len(matches) > 0 { 
			// Add the identified host to the map with an empty array as its value unless 
			//	* it is already present 
			//	* it should be ignored
			// and continue parsing the next line
			current_host = matches[1]
			_, found := ignore_hosts[current_host]

			if hosts_map[current_host] == nil && !found {
				hosts_map[current_host] = make([]string, 0) 
			}
			continue
		}

		// Look for a ProxyJump if the line did not have a `Host` value
		matches = proxyjump_regex.FindStringSubmatch(line)

		if len(matches) > 0 { 
			jump_hosts 	  :=  strings.Split(matches[1], ",")
			
			// Set the inital jump location as the 'main' host
			main_host := jump_hosts[0]
			
			// The `current_host` will be set to the most recently read `Host` line
			// and should be set as the exit node when more than one proxy exists
			//  main_host				=  jump_to_1
			//	hosts_map[main_host]	+= jump_to_2
			//	hosts_map[jump_to_2]	+= jump_to_3
			//	...
			//	hosts_map[jump_to_n]	+= current_host

			jump_hosts = append(jump_hosts[1:], current_host)
			
			for _, jump_to := range jump_hosts {

				_ , main_found  := ignore_hosts[main_host]
				_ , jump_found  := ignore_hosts[jump_to]

				if !main_found && !jump_found {
					// Only add hosts if neither one should be ignored
					if hosts_map[main_host] == nil {
						// If the main_host doesn't have a key, initalise its array with the jump_to host
						hosts_map[main_host] = []string { jump_to } 	
					} else {
						// Otherwise append the jump_to host
						hosts_map[main_host] = append(hosts_map[main_host], jump_to) 
					}
					
					// Record each host that has a jump_to entry
					has_jump[jump_to] = struct{}{}
				}

				main_host = jump_to
			}
		}
	} 
	
	if tree {
		for host := range has_jump {
			// Remove all top-level hosts that have a jump_to entry
			Debug("deleting", host)
			delete(hosts_map, host)
		}
	}

	return hosts_map, has_jump
}


//func GetTreeMapping(config_file string, ignore_hosts map[string]struct{}) (host_map map[string][]string) {
//
//	f, err := os.Open(config_file) 
//	if err != nil { 
//		Die(err.Error()) 
//	}
//	defer f.Close()
//
//	scanner 	:= bufio.NewScanner(f)
//	
//	// (?i) provides case-Insensitive matching
//	host_regex 			:= regexp.MustCompile(`(?i)^\s*Host\s+([^\s]*)`)
//	proxyjump_regex 	:= regexp.MustCompile(`(?i)^\s*ProxyJump\s+([^\s]*)`)
//	
//	// Maps hostname -> [jump_to_hosts]
//	hosts_map 		 	:= make(map[string][]string)
//	current_host := ""
//	
//	for scanner.Scan() {
//		
//		line := scanner.Text()
//
//		// The first match will be the entire string and the 
//		// remaining indices will hold capture groups i.e. the hostname 
//		matches := host_regex.FindStringSubmatch(line)
//		
//		if len(matches) > 0 { 
//			// Add the identified host to the map with an empty array as its value
//			// unless it has already been encountered or should be ignored 
//			// and continue parsing the next line
//			current_host = matches[1]
//			_, found := ignore_hosts[current_host]
//
//			if hosts_map[current_host] == nil && !found {
//				hosts_map[current_host] = make([]string, 0) 
//			}
//			continue
//		}
//
//		// Look for a ProxyJump if the line did not have a `Host` value
//		matches = proxyjump_regex.FindStringSubmatch(line)
//
//		if len(matches) > 0 { 
//			jump_hosts 	  :=  strings.Split(matches[1], ",")
//			
//			// Set the inital jump location as the 'main' host
//			main_host := jump_hosts[0]
//			
//			// The `current_host` will be set to the most recently read `Host` line
//			// and should be set as the exit node when more than one proxy exists
//			//  main_host				=  jump_to_1
//			//	hosts_map[main_host]	+= jump_to_2
//			//	hosts_map[jump_to_2]	+= jump_to_3
//			//	...
//			//	hosts_map[jump_to_n]	+= current_host
//
//			jump_hosts = append(jump_hosts[1:], current_host)
//			
//			for _, jump_to := range jump_hosts {
//
//				_ , main_found  := ignore_hosts[main_host]
//				_ , jump_found  := ignore_hosts[jump_to]
//
//				if !main_found && !jump_found {
//					// Only add hosts if neither one should be ignored
//					if hosts_map[main_host] == nil {
//						// If the main_host doesn't have a key, initalise its array with the jump_to host
//						hosts_map[main_host] = []string { jump_to } 
//						
//						// Remove the root-level entry for the jump_to host if one exists
//						hosts_map[jump_to] = nil
//					} else {
//						// Otherwise append the jump_to host
//						hosts_map[main_host] = append(hosts_map[main_host], jump_to) 
//						
//						// Remove the root-level entry for the jump_to host if one exists
//						hosts_map[jump_to] = nil
//					}
//				}
//
//				main_host = jump_to
//			}
//		}
//	} 
//
//	return hosts_map
//}
//