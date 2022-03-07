package entity

import "github.com/jinzhu/gorm"

type Group struct {
	ID         int    `gorm:"primary_key;column:id;type:int(11) unsigned;not null"`
	Name       string `gorm:"column:name;type:varchar(50);not null"`               // 集团名称
	ShortName  string `gorm:"column:short_name;type:varchar(30)"`                  // 企业简称
	Address    string `gorm:"column:address;type:varchar(255);not null"`           // 地址
	Intro      string `gorm:"column:intro;type:text"`                              // 简介
	Logo       string `gorm:"column:logo;type:varchar(255);not null"`              // 集团logo
	AppID      string `gorm:"column:app_id;type:varchar(100);not null"`            // 微信开发者凭据 AppId
	AppSecret  string `gorm:"column:app_secret;type:varchar(2000);not null"`       // 开发者凭据 AppSecret
	MchID      string `gorm:"column:mch_id;type:varchar(200);not null"`            // 微信支付商户号
	PartnerKey string `gorm:"column:partner_key;type:varchar(200);not null"`       // 支付密钥
	QqMapKey   string `gorm:"column:qq_map_key;type:varchar(200);not null"`        // 腾讯位置服务key
	MaxAreaNum int    `gorm:"column:max_area_num;type:int(10) unsigned;not null"`  // 最大门店数量
	IsDel      int8   `gorm:"column:is_del;type:tinyint(1) unsigned;not null"`     // 是否删除：0正常，1删除
	CreatedAt  int    `gorm:"column:created_at;type:int(11) unsigned;not null"`    // 创建时间
	UpdatedAt  int    `gorm:"column:updated_at;type:int(11) unsigned;not null"`    // 更新时间
	CreatedBy  int    `gorm:"column:created_by;type:int(10) unsigned;not null"`    // 创建人员
	UpdatedBy  int    `gorm:"column:updated_by;type:int(10) unsigned;not null"`    // 更新人员
	AdminID    int    `gorm:"column:admin_id;type:int(10) unsigned;not null"`      // 超级管理员id
	ProvinceID int    `gorm:"column:province_id;type:int(10) unsigned;not null"`   // 省份id
	CityID     int    `gorm:"column:city_id;type:int(10) unsigned;not null"`       // 城市id
	RegionID   int    `gorm:"column:region_id;type:int(10) unsigned;not null"`     // 区域id
	VersionID  int8   `gorm:"column:version_id;type:tinyint(3) unsigned;not null"` // 集团版本id
	Status     int8   `gorm:"column:status;type:tinyint(1) unsigned;not null"`     // 系统状态：0开启，1关闭
	DueTime    int    `gorm:"column:due_time;type:int(10) unsigned;not null"`      // 集团到期时间

	Admin  Admin   `gorm:"ForeignKey:ID;AssociationForeignKey:AdminID"`
	Admins []Admin `gorm:"ForeignKey:GroupID;AssociationForeignKey:ID"`
}

func (m Group) TableName() string {
	return "cu_group"
}

// 获取关联的员工
func (m Group) GetAdmins(db *gorm.DB) ([]Admin, error) {
	var models []Admin
	err := db.Debug().Model(&m).Where("is_del = 0").Association("Admins").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
