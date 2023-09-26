package md5

import (
	"fmt"
	"os"
	"os/user"
	"testing"
)

func TestHexMd5(t *testing.T) {
	fmt.Println(HexMd5("1"))
}

func TestFileHexMd5(t *testing.T) {
	u, _ := user.Current()
	fmt.Println(FileHexMd5(u.HomeDir + string(os.PathSeparator) + ".zprofile"))
}
