package lib

import (
	//"fmt"
	//"os/exec"
	//"strings"
)

func GetUnameMapping(hosts_map map[string][]string, config_file string, ignore_hosts map[string]struct{} ) (uname_mapping map[string]string) {

	uname_mapping = make(map[string]string, len(hosts_map))

	for host, jump_hosts := range hosts_map {

		addToUnameMapping(host, config_file, uname_mapping, ignore_hosts)

		for _, jump_to := range jump_hosts {
			// We need to go through each jump host in case they don't have their own entry
			addToUnameMapping(jump_to, config_file, uname_mapping, ignore_hosts)
		}
	}

	return uname_mapping
}

func addToUnameMapping(host, config_file string, uname_mapping map[string]string, ignore_hosts map[string]struct{}){
	if _, ok := ignore_hosts[host]; !ok {
		// Continue to next if the host is in the ignore list
		if uname_mapping[host] == "" {
			// If their uname_mapping is empty, fetch the `uname` output
			uname_mapping[host] = GetUname(host, config_file)
		}
	} 
} 


/// TODO error handling
func GetUname(host, config_file string) string {
	// Linux 5.11.4-1-ARCH aarch64
	// Linux 5.13.9-arch1-1 x86_64
	// Darwin 20.5.0 x86_64
	
	//uname, err := exec.Command("ssh", "-F", config_file,  host,  "uname", "-rms").Output()
	//if err != nil {
	//	Die(err.Error())
	//}

	//return  strings.TrimSuffix(string(uname), "\n") 
	return host 
}