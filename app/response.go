package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ResponseExtra struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Response
}

// Response setting gin.JSON
func (g *Gin) Success(httpCode int, msg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  msg,
		Data: data,
	})
}

func (g *Gin) SuccessExtra(total int64, page int, pageSize, httpCode int, msg string, data interface{}) {
	g.C.JSON(httpCode, ResponseExtra{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Response: Response{
			Code: httpCode,
			Msg:  msg,
			Data: data},
	})
}

func (g *Gin) Fail(httpCode int, err error, data interface{}) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		g.C.JSON(httpCode, Response{
			Code: httpCode,
			Msg:  err.Error(),
			Data: data,
		})
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	g.C.JSON(httpCode, Response{
		Code: httpCode,
		Msg:  "validate error",
		Data: removeTopStruct(errs.Translate(trans)),
	})
}

//去掉结构体名称标识
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
