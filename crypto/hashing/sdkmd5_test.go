package hashing

import (
	"fmt"
	"os"
	"os/user"
	"testing"
)

func TestHexMd5(t *testing.T) {
	fmt.Println(Md5Hex("1"))
}

func TestFileHexMd5(t *testing.T) {
	u, _ := user.Current()
	fmt.Println(Md5FileHex(u.HomeDir + string(os.PathSeparator) + ".zprofile"))
}
