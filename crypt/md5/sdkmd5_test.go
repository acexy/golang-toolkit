package md5

import (
	"fmt"
	"os"
	"os/user"
	"testing"
)

func TestMd5Hex(t *testing.T) {
	fmt.Println(Md5Hex("1"))
}

func TestMd5FileHex(t *testing.T) {
	u, _ := user.Current()
	fmt.Println(Md5FileHex(u.HomeDir + string(os.PathSeparator) + ".zprofile"))
}
