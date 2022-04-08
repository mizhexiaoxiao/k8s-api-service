package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetPathParameterInt
// 获取请求路径上的id参数
func GetPathParameterInt(context *gin.Context, args ...string) (map[string]int, error) {
	m := make(map[string]int)
	for _, arg := range args {
		paramStr := context.Param(arg)
		param, err := strconv.Atoi(paramStr)
		if err != nil {
			return nil, err
		}
		m[arg] = param
	}
	return m, nil
}

func GetPathParameterString(context *gin.Context, args ...string) (map[string]string, error) {
	m := make(map[string]string)
	for _, arg := range args {
		paramStr := context.Param(arg)
		if paramStr == "" {
			return nil, errors.New(fmt.Sprintf("%s path param not valid", arg))
		}
		m[arg] = paramStr
	}
	return m, nil
}
