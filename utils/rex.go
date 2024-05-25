package utils

import (
	"fmt"
	"regexp"
)

type Function struct {
	Name          string
	ParameterName string
}

func GetFunction(s string) ([]Function, error) {
	// 定义正则表达式
	// 使用捕获组分别捕获函数名、参数类型和返回值类型，包括流类型的支持
	regex := `rpc\s+(\w+)\s*\(\s*(stream\s+)?\.*([\w\.]+)\s*\)\s+returns\s+\(\s*(stream\s+)?\.*([\w\.]+)\s*\);`

	// 编译正则表达式
	re := regexp.MustCompile(regex)

	// 执行全局匹配
	matches := re.FindAllStringSubmatch(s, -1)

	var funcs []Function

	// 检查匹配结果
	if len(matches) > 0 {
		for _, match := range matches {
			if len(match) == 6 {
				funcs = append(funcs, Function{
					Name:          match[1],
					ParameterName: match[3],
				})
			}
		}
	} else {
		return nil, fmt.Errorf("no match found")
	}

	return funcs, nil
}
