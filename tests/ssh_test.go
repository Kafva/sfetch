package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Kafva/sfetch/lib"
	"github.com/stretchr/testify/assert"
)

var home, _ = os.UserHomeDir()

func Test_GetUname(t *testing.T) {
	assert.Equal(t, 
		"Linux 5.11.4-1-ARCH aarch64", 
		lib.GetUname("kafva.one", fmt.Sprintf("%s/.ssh/config", home) ),
	)
}
