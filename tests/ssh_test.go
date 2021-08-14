package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Kafva/sfetch/lib"
	"github.com/stretchr/testify/assert"
)

var home, _ 	= os.UserHomeDir()
var config_file = fmt.Sprintf("%s/.ssh/config", home)

func Test_GetUname(t *testing.T) {
	assert.Equal(t, 
		"kafva.one",
		//"Linux 5.11.4-1-ARCH aarch64", 
		lib.GetUname("kafva.one",  config_file),
	)
}

func Test_GetUnameMapping(t *testing.T) {
	hosts_map 		:= lib.GetHostMapping(config_file)
	uname_mapping 	:= lib.GetUnameMapping(hosts_map, config_file, make(map[string]struct{}))
	
	fmt.Println(uname_mapping)
}