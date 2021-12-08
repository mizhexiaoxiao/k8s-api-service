package cluster

import (
	"github.com/mizhexiaoxiao/k8s-api-service/app"
	"github.com/mizhexiaoxiao/k8s-api-service/models"
	"gorm.io/gorm"
)

func Create(data models.ClusterModel) (err error) {
	err = models.DB.Model(&models.ClusterModel{}).Create(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func Update(id int, data models.ClusterModel) (err error) {
	err = models.DB.Model(&models.ClusterModel{}).Where("id = ?", id).Updates(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func List(pageInfo app.PageInfo) (clusters []*models.ClusterModel, err error) {
	err = models.DB.Model(&models.ClusterModel{}).Offset((pageInfo.Page - 1) * pageInfo.PageSize).Limit(pageInfo.PageSize).Find(&clusters).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return clusters, nil
}

func Get(id int) (cluster models.ClusterModel, err error) {
	err = models.DB.Model(&models.ClusterModel{}).Where("id = ?", id).First(&cluster).Error
	return
}

func Delete(id int) (err error) {
	//软删除
	//err = models.DB.Model(&models.ClusterModel{}).Delete("id = ?", id).Error
	//硬删除
	err = models.DB.Model(&models.ClusterModel{}).Unscoped().Delete("id = ?", id).Error
	return
}

func Count() (count int64, err error) {
	if err := models.DB.Model(&models.ClusterModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
