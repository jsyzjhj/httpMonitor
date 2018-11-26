/**
请求结果验证类
 */
package lib

import (
	"errors"
	"regexp"
	"strconv"
)

/**
自动验证方法
0为状态码匹配
1为内容全等匹配
2为正则匹配
 */
func verify(vtype, code int, content string, vval string) (bool, error) {
	switch vtype {
	case 0:
		if vvals, _ := strconv.Atoi(vval); vvals == code {
			return true, nil
		} else {
			return false, nil
		}
	case 1:
		if content == vval {
			return true, nil
		} else {
			return false, nil
		}
	case 2:
		is, err := regexp.Match(vval, []byte(content))
		if err != nil {
			return false, err
		}
		if is {
			return true, nil
		} else {
			return false, nil
		}
	default:
		return false, errors.New("判断方式不正确！")
	}
}
