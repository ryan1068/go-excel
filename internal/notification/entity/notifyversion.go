package entity

import "github.com/jinzhu/gorm"

type NotifyVersion struct {
	ID        int   `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`
	NotifyID  int   `gorm:"index;column:notify_id;type:int(10) unsigned;not null"`  // 通知id
	VersionID int64 `gorm:"index;column:version_id;type:int(10) unsigned;not null"` // 版本id
	IsDel     int8  `gorm:"column:is_del;type:tinyint(1) unsigned;not null"`        // 是否已删除，0否，1是
	CreatedAt int   `gorm:"column:created_at;type:int(11) unsigned;not null"`       // 创建时间
	CreatedBy int   `gorm:"column:created_by;type:int(10) unsigned;not null"`       // 创建人
	UpdatedAt int   `gorm:"column:updated_at;type:int(11) unsigned;not null"`       // 更新时间
	UpdatedBy int   `gorm:"column:updated_by;type:int(10) unsigned;not null"`       // 更新人

	Groups []Group `gorm:"ForeignKey:VersionID;AssociationForeignKey:VersionID"`
}

func (m NotifyVersion) TableName() string {
	return "cu_notify_version"
}

// 获取关联的集团
func (m NotifyVersion) GetGroups(db *gorm.DB) ([]Group, error) {
	var models []Group
	err := db.Model(&m).Where("is_del = 0").Association("Groups").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (m NotifyVersion) GetGroupAdmins(db *gorm.DB, vid int64) ([]Admin, error) {
	var models []Admin
	err := db.Model(&m).Joins("left join cu_group on cu_group.id = cu_group_admin.group_id").
		Where("cu_group_admin.is_del = 0 and cu_group.version_id = ?", vid).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (m NotifyVersion) GetAreaAdmins(db *gorm.DB, vid int64) ([]AreaAdmin, error) {
	var models []AreaAdmin
	err := db.Model(&m).Joins("left join cu_group on cu_group.id = cu_area_admin.group_id").
		Where("cu_area_admin.is_del = 0 and cu_group.is_del = 0 and cu_group.version_id = ?", vid).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
