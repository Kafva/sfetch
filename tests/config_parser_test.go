package tests

import (
	. "github.com/Kafva/sfetch/lib"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)


func Test_GetIgnoreHosts(t *testing.T) {
	ignore_hosts := GetIgnoreHosts(".mocks/sfetchignore")
	assert.Contains(t, ignore_hosts, "github.com")
	assert.NotContains(t, ignore_hosts, "#vel")
	assert.NotContains(t, ignore_hosts, "vel")
}

/// A ProxyJump reads the config entries
/// from the original machine and not the current
/// node but the /etc/hosts on the remote machine matters
///
/// Even though opt2 is not ann explicit jump host
/// from opt1 it will still appear that way in the tree
/// It is possible for the same host to appear more than once
/// if one can reach it using several proxy paths, e.g. 
/// 	...streak3 -> streak4 -> end
///
/// localhost
/// └── streak1
/// │   ├── streak2
/// │       └── streak3
/// │           └── streak4
/// │               └── end
/// └── loc10
///     └── loc20
///     │   ├── loc31
///     │   │   ├── loc40
///     │   ├── loc30
///     └── loc21
///     └── opt1
///         └── opt2
///             └── opt_loc
///             │   ├── streak3
///             │       └── streak4
///             │           └── end
///             └── opt_loc_2
func Test_GetHostMapping(t *testing.T) {
	hosts, has_jump := GetHostMapping(".mocks/ssh_config", make(map[string]struct{}))

	assert.Contains(t,  hosts["loc10"], "loc20", "loc21", "opt1")
	assert.Contains(t,  hosts["loc20"], "loc30", "loc31")
	assert.Contains(t, hosts["loc31"], "loc40"  )
	assert.Contains(t, hosts["opt1"] , "opt2"    )
	assert.Contains(t, hosts["opt2"] , "opt_loc" )
	assert.Contains(t, hosts["opt_loc"] , "streak3") 
	assert.Contains(t, hosts["streak1"] , "streak2") 
	assert.Contains(t, hosts["streak2"] , "streak3") 
	assert.Contains(t, hosts["streak3"] , "streak4") 
	assert.Contains(t, hosts["streak4"] , "end") 

	assert.Empty(t, hosts["loc21"])
	assert.Empty(t, hosts["loc40"])
	assert.Empty(t, hosts["opt_loc_2"])
	assert.Empty(t, hosts["end"])

	
	assert.NotContains(t, has_jump, "loc10", "opt2")
	assert.NotContains(t, has_jump, "loc10", "opt2")
			
	fmt.Println("---------------------------")
	fmt.Println(hosts)

}