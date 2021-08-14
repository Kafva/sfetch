package lib

import (
	"os/exec"
	"strings"
)

func GetUname() string {
	uname, _ := exec.Command("uname", "-m").Output()

	return  strings.TrimSuffix(string(uname), "\n") 
}