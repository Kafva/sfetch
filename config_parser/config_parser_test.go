package config_parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseSshConfig(t *testing.T){
	hosts := ParseSshConfig("../.tests/ssh_config_1")

	assert.Contains(t,  hosts["loc10"], "loc20", "loc21", "opt1")
	assert.Contains(t,  hosts["loc20"], "loc30", "loc31")
	assert.Contains(t, hosts["loc31"], "loc40"  )
	assert.Contains(t, hosts["opt1"] , "opt2"    )
	assert.Contains(t, hosts["opt2"] , "opt_loc" )

	assert.Empty(t, hosts["loc21"])
	assert.Empty(t, hosts["loc40"])
	assert.Empty(t, hosts["opt_loc"])
	assert.Empty(t, hosts["opt_loc_2"])

	fmt.Println("---------------------------")
	fmt.Println(hosts)

}