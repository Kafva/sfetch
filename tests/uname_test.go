package tests

import (
	. "github.com/Kafva/sfetch/lib"
	"os"
	"io"
	"testing"
	"github.com/stretchr/testify/assert"
)


func Test_ParseSystemInfo(t *testing.T) {

	f, err := os.Open(".mocks/systeminfo") 
	if err != nil { 
		Die(err.Error()) 
	}
	defer f.Close()
	
	systemInfo, err := io.ReadAll(f)	
	if err != nil {
		Die(err.Error())
	}

	uname := ParseSystemInfo(string(systemInfo))
	
	assert.Equal(t, "Microsoft Windows 10 Education 10.0.19041 x64-based PC", uname)
}
