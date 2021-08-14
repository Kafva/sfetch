package lib

import (
	"os/exec"
	"strings"
)

func GetUname(host, config_file string) string {
	// Linux 5.11.4-1-ARCH aarch64
	// Linux 5.13.9-arch1-1 x86_64
	// Darwin 20.5.0 x86_64
	uname, err := exec.Command("ssh", "-F", config_file,  host,  "uname", "-rms").Output()
	if err != nil {
		Die(err.Error())
	}

	return  strings.TrimSuffix(string(uname), "\n") 
}