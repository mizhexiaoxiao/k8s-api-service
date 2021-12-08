package models

import (
	"gorm.io/datatypes"
)

type ClusterModel struct {
	Model
	Cluster
}

type Cluster struct {
	Name    string         `json:"name" gorm:"unique"`
	Desc    string         `json:"desc"`
	Context datatypes.JSON `json:"context"`
	State   bool           `json:"state"`
}
