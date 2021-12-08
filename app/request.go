package app

type PageInfo struct {
	Page     int `json:"page" form:"page" binding:"required,gte=1"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize" binding:"required,gte=1"` // 每页大小
}

type GetById struct {
	ID int `json:"id" uri:"id"` // 主键ID
}
