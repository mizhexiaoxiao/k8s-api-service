package metadata

type CommonQueryParameter struct {
	NameSpace     string `form:"namespace" binding:"required"`
	LabelSelector string `form:"labelSelector" binding:"required"`
}
