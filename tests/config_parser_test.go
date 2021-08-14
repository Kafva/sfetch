package tests

import (
	"github.com/Kafva/sfetch/lib"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)


func Test_GetIgnoreHosts(t *testing.T) {
	ignore_hosts := lib.GetIgnoreHosts(".mocks/sfetchignore")
	assert.Contains(t, ignore_hosts, "github.com")
	assert.NotContains(t, ignore_hosts, "#vel")
	assert.NotContains(t, ignore_hosts, "vel")
}

func Test_GetHostMapping(t *testing.T) {
	hosts := lib.GetHostMapping(".mocks/ssh_config")

	assert.Contains(t,  hosts["loc10"], "loc20", "loc21", "opt1")
	assert.Contains(t,  hosts["loc20"], "loc30", "loc31")
	assert.Contains(t, hosts["loc31"], "loc40"  )
	assert.Contains(t, hosts["opt1"] , "opt2"    )
	assert.Contains(t, hosts["opt2"] , "opt_loc" )
	assert.Contains(t, hosts["streak1"] , "streak2") 
	assert.Contains(t, hosts["streak2"] , "streak3") 
	assert.Contains(t, hosts["streak3"] , "streak4") 
	assert.Contains(t, hosts["streak4"] , "end") 

	assert.Empty(t, hosts["loc21"])
	assert.Empty(t, hosts["loc40"])
	assert.Empty(t, hosts["opt_loc"])
	assert.Empty(t, hosts["opt_loc_2"])
	assert.Empty(t, hosts["end"])

	fmt.Println("---------------------------")
	fmt.Println(hosts)

}