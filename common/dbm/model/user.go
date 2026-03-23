package model

import (
	"time"

	"github.com/edgehook/ithings/common/global"
	"k8s.io/klog/v2"
)

type User struct {
	ID              string `gorm:"column:id; type:varchar(36); primary_key;" json:"id"`
	Name            string `gorm:"column:name; not null; unique; type:varchar(256);" json:"name"`
	Password        string `gorm:"column:password; not null; type:varchar(256);" json:"password"`
	Email           string `gorm:"column:email; type:varchar(256);" json:"email"`
	Rule            string `gorm:"column:rule; type:varchar(256);" json:"rule"`
	CreateTimeStamp int64  `gorm:"column:create_time_stamp;" json:"createTimeStamp"`
	UpdateTimeStamp int64  `gorm:"column:update_time_stamp;autoUpdateTime:milli" json:"updateTimeStamp"`
}

func (User) TableName() string {
	return "user"
}
func GetUsers() ([]*User, error) {
	var users []*User
	err := global.DBAccess.Order("update_time_stamp desc").Find(&users).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return users, err
}

func GetUserByPage(page int, limit int) ([]*User, error) {
	var users []*User
	err := global.DBAccess.Offset((page - 1) * limit).Limit(limit).Order("create_time_stamp desc").Find(&users).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return users, err
}

func GetUserByPageAndKeywords(page int, limit int, keywords string) ([]*User, error) {
	var Users []*User
	err := global.DBAccess.Where("name LIKE ?", "%"+keywords+"%").Offset((page - 1) * limit).Order("create_time_stamp desc").Limit(limit).Find(&Users).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}
	return Users, err
}
func GetUserCountByKeywords(keywords string) (int64, error) {
	var count int64
	err := global.DBAccess.Model(&User{}).Where("name LIKE ?", "%"+keywords+"%").Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}
func GetUserByName(name string) (*User, error) {
	User := &User{}
	err := global.DBAccess.Where("name = ?", name).First(User).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return nil, err
	}

	return User, err
}

func GetUserCount() (int64, error) {
	var count int64
	err := global.DBAccess.Model(&User{}).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return -1, err
	}
	return count, err
}

func IsExistUserByName(name string) bool {
	if name == "" {
		return false
	}
	var count int64
	err := global.DBAccess.Model(&User{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

func GetUserById(id string) (*User, error) {
	var User *User
	err := global.DBAccess.First(&User, id).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return User, err
	}
	return User, err
}

func AddUser(User *User) error {
	User.CreateTimeStamp = time.Now().UnixNano() / 1e6
	err := global.DBAccess.Create(&User).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SaveUser(id, password, rule, email string) error {
	vals := make(map[string]interface{})
	if password != "" {
		vals["password"] = password
	}
	if rule != "" {
		vals["rule"] = rule
	}

	if email != "" {
		vals["email"] = email
	}
	err := global.DBAccess.Model(&User{}).Where("id = ?", id).Updates(vals).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DeleteUser(id string) error {
	err := global.DBAccess.Where("id = ?", id).Delete(&User{}).Error
	if err != nil {
		klog.Errorf("err: %v", err)
		return err
	}
	return nil

}
