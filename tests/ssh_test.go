package tests
//
//import (
//	"fmt"
//	"os"
//	"testing"
//
//	"github.com/Kafva/sfetch/lib"
//	"github.com/stretchr/testify/assert"
//)
//
//var home, _ 	= os.UserHomeDir()
//var config_file = fmt.Sprintf("%s/.ssh/config", home)
//
//func Test_GetUnameMapping(t *testing.T) {
//	hosts_map 		:= lib.GetHostMapping(config_file)
//	ignore_hosts 	:= lib.GetIgnoreHosts(".mocks/sfetchignore")
//	uname_mapping 	:= lib.GetUnameMapping(hosts_map, config_file, ignore_hosts)
//	
//	fmt.Println(uname_mapping)
//	fmt.Println(hosts_map)
//}