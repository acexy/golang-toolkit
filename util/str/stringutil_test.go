package str

import (
	"fmt"
	"testing"
)

func TestSubstring(t *testing.T) {
	fmt.Println(Substring("1234", 1, 3))
	fmt.Println(Substring("你好 w", 1))

	str := "&莫=1"
	fmt.Println(Substring(str, 1))
}

func TestCharLength(t *testing.T) {
	var str = "hello 你好"
	fmt.Println(len(str))
	fmt.Println(CharLength(str))
}

func TestBuilder(t *testing.T) {
	b := NewBuilder("你好").WriteByte(97).WriteRune(100).WriteString("1")
	fmt.Println(b.ToString())
}

func TestHasText(t *testing.T) {
	fmt.Println(HasText(" a b "))
	fmt.Println(HasText("   a  "))
	fmt.Println(HasText("     "))
}
