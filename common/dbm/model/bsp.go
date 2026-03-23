package model

import (
	"time"

	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type Bsp struct {
	ID       string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Filename string `gorm:"column:filename; not null; type:varchar(256);" json:"filename"`
	// Boardname       string `gorm:"column:boardname; not null; type:varchar(256);" json:"boardname"`
	// Version         string `gorm:"column:version; type:varchar(256);" json:"version"`
	Path            string `gorm:"column:path; type:varchar(256);" json:"path"`
	Md5             string `gorm:"column:md5; type:varchar(256);" json:"md5"`
	Size            int64  `gorm:"column:size;" json:"size"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (Bsp) TableName() string {
	return "bsp"
}
func GetBsps() ([]*Bsp, error) {
	var Bsps []*Bsp
	err := global.DBAccess.Order("update_time_stamp desc").Find(&Bsps).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return Bsps, err
}

func GetBspByPage(page int, limit int) ([]*Bsp, error) {
	var Bsps []*Bsp
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&Bsps).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return Bsps, err
}

func GetBspByPageAndKeywords(page int, limit int, keywords string) ([]*Bsp, error) {
	var Bsps []*Bsp
	err := global.DBAccess.Where("boardname LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&Bsps).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return Bsps, err
}
func GetBspCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&Bsp{}).Where("boardname LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func IsExistByMd5(filename, md5 string) bool {
	if filename == "" || md5 == "" {
		return false
	}
	var count int64
	err := global.DBAccess.Model(&Bsp{}).Where("filename = ? and md5 = ?", filename, md5).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}
func GetBspCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&Bsp{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func GetBspById(id string) (*Bsp, error) {
	var Bsp *Bsp
	err := global.DBAccess.First(&Bsp, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return Bsp, err
	}
	return Bsp, err
}

func AddBsp(bsp *Bsp) error {
	bsp.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&bsp).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}
func DeleteBsp(id string) error {
	err := global.DBAccess.Where("id = ?", id).Delete(&Bsp{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil

}
