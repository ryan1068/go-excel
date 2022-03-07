package entity

import "github.com/jinzhu/gorm"

type Notify struct {
	ID             int    `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`
	Type           uint8  `gorm:"index;column:type;type:tinyint(2) unsigned;not null"`          // 消息类型：0系统公告，1版本通知
	Title          string `gorm:"column:title;type:varchar(50);not null"`                       // 消息标题
	ContentType    int8   `gorm:"column:content_type;type:tinyint(1) unsigned;not null"`        // 内容类型：0消息详情，1跳转URL
	Content        string `gorm:"column:content;type:text;not null"`                            // 消息的内容
	SendStatus     int8   `gorm:"column:send_status;type:tinyint(1) unsigned;not null"`         // 发送状态：0待发送，1已发送
	SendType       int8   `gorm:"index;column:send_type;type:tinyint(1) unsigned;not null"`     // 发送模式：0定时发送，1及时发送
	SendTime       int    `gorm:"column:send_time;type:int(10) unsigned;not null"`              // 发送时间
	CompletionTime int64  `gorm:"column:completion_time;type:int(10) unsigned;not null"`        // 任务完成时间
	ReceiveType    uint8  `gorm:"index;column:receive_type;type:tinyint(10) unsigned;not null"` // 接收群体：0所有版本，1指定版本，2指定用户
	Channel        uint8  `gorm:"column:channel;type:tinyint(3) unsigned;not null"`             // 通知渠道：0 WEB后台
	IsDel          int8   `gorm:"column:is_del;type:tinyint(1) unsigned;not null"`              // 是否已删除，0否，1是
	CreatedAt      int    `gorm:"column:created_at;type:int(11) unsigned;not null"`             // 创建时间
	CreatedBy      int    `gorm:"column:created_by;type:int(10) unsigned;not null"`             // 创建人
	UpdatedAt      int    `gorm:"column:updated_at;type:int(11) unsigned;not null"`             // 更新时间
	UpdatedBy      int    `gorm:"column:updated_by;type:int(10) unsigned;not null"`             // 更新人

	Versions []NotifyVersion `gorm:"ForeignKey:NotifyID;AssociationForeignKey:ID"`
}

func (m Notify) TableName() string {
	return "cu_notify"
}

// 获取关联的版本
func (m Notify) GetVersions(db *gorm.DB) ([]NotifyVersion, error) {
	var models []NotifyVersion
	err := db.Debug().Model(&m).Where("is_del = 0").Association("Versions").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return models, nil
}
