package tests

import (
	"github.com/Kafva/sfetch/lib"
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_GetUname(t *testing.T) {
	assert.Equal(t, 
		"x86_64", 
		lib.GetUname(),
	)
}
