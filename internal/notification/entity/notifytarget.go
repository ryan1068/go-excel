package entity

type NotifyTarget struct {
	ID        int   `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`
	NotifyID  int   `gorm:"index;column:notify_id;type:int(10) unsigned;not null"` // 通知id
	VersionID int8  `gorm:"column:version_id;type:tinyint(1) unsigned;not null"`   // 版本id
	Scene     int8  `gorm:"index;column:scene;type:tinyint(1) unsigned;not null"`  // 发送场景：0集团，1门店
	AdminID   int   `gorm:"index;column:admin_id;type:int(10) unsigned;not null"`  // 管理员id
	IsRead    int8  `gorm:"column:is_read;type:tinyint(1) unsigned;not null"`      // 是否已读：0否，1是
	IsDel     int8  `gorm:"column:is_del;type:tinyint(1) unsigned;not null"`       // 是否已删除，0否，1是
	CreatedAt int64 `gorm:"column:created_at;type:int(11) unsigned;not null"`      // 创建时间
	CreatedBy int64 `gorm:"column:created_by;type:int(10) unsigned;not null"`      // 创建人
	UpdatedAt int64 `gorm:"column:updated_at;type:int(11) unsigned;not null"`      // 更新时间
	UpdatedBy int64 `gorm:"column:updated_by;type:int(10) unsigned;not null"`      // 更新人
}

func (m NotifyTarget) TableName() string {
	return "cu_notify_target"
}
