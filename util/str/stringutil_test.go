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

func TestIsCamelCase(t *testing.T) {
	// 测试字符串
	testCases := []string{"camelCase", "PascalCase", "snake_case", "camel_case", "camelCase123", "camelCaseExample"}
	for _, testCase := range testCases {
		fmt.Printf("%s: %v\n", testCase, IsCamelCase(testCase))
	}
}

func TestCamelToSnake(t *testing.T) {
	testCases := []string{"camelCase", "CamelCase", "simpleTestExample", "OneMoreExample", "like"}

	for _, testCase := range testCases {
		fmt.Printf("%s -> %s\n", testCase, CamelToSnake(testCase))
	}
}

func TestSnakeToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},                      // 空字符串
		{"hello", "hello"},            // 单词
		{"hello_world", "helloWorld"}, // 下划线分隔的小写
		{"snake_case_example", "snakeCaseExample"},          // 多个单词
		{"_leading_underscore", "leadingUnderscore"},        // 前置下划线
		{"trailing_underscore_", "trailingUnderscore"},      // 后置下划线
		{"mixed_case_example_TEST", "mixedCaseExampleTEST"}, // 混合大小写
	}

	for _, test := range tests {
		result := SnakeToCamel(test.input)
		if result != test.expected {
			t.Errorf("SnakeToCamel(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
