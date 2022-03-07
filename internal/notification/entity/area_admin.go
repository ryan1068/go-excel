package entity

type AreaAdmin struct {
	ID            int    `gorm:"primary_key;column:id;type:int(10) unsigned;not null"`
	GroupID       int    `gorm:"column:group_id;type:int(10) unsigned;not null"`          // 集团id
	AreaID        int    `gorm:"index;column:area_id;type:int(10) unsigned;not null"`     // 门店id
	Account       string `gorm:"index;column:account;type:varchar(15);not null"`          // 账号
	Nickname      string `gorm:"column:nickname;type:varchar(15);not null"`               // 昵称
	Number        string `gorm:"column:number;type:varchar(20);not null"`                 // 员工编号
	Sex           string `gorm:"column:sex;type:varchar(1);not null"`                     // 性别
	Avatar        string `gorm:"column:avatar;type:varchar(100);not null"`                // 头像
	Tel           string `gorm:"column:tel;type:varchar(11);not null"`                    // 手机
	Password      string `gorm:"column:password;type:varchar(100);not null"`              // 密码
	PostStatus    int8   `gorm:"column:post_status;type:tinyint(1) unsigned;not null"`    // 岗位状态：0在职，1离职
	AccountStatus int8   `gorm:"column:account_status;type:tinyint(1) unsigned;not null"` // 账号状态：0启用，1禁用
	DepartmentID  int    `gorm:"column:department_id;type:int(10) unsigned;not null"`     // 部门id
	Type          int8   `gorm:"column:type;type:tinyint(1) unsigned;not null"`           // 用户类型:0自定义增加，1业财增加
	CurrentRoleID int    `gorm:"column:current_role_id;type:int(10) unsigned;not null"`   // 当前角色
	AuthKey       string `gorm:"column:auth_key;type:varchar(100);not null"`              // 认证key
	WechatID      string `gorm:"column:wechat_id;type:varchar(50);not null"`              // 微信号
	Email         string `gorm:"column:email;type:varchar(50);not null"`                  // 邮箱
	Intro         string `gorm:"column:intro;type:varchar(255);not null"`                 // 简介
	Source        int8   `gorm:"column:source;type:tinyint(1) unsigned;not null"`         // 用户来源:0门店添加，1集团添加
	IsDel         uint8  `gorm:"index;column:is_del;type:tinyint(4) unsigned;not null"`   // 是否已删除，0否，1是
	CreatedAt     int    `gorm:"index;column:created_at;type:int(10) unsigned;not null"`  // 创建时间
	CreatedBy     int    `gorm:"column:created_by;type:int(10) unsigned;not null"`        // 创建人
	UpdatedAt     int    `gorm:"column:updated_at;type:int(10) unsigned;not null"`        // 更新日期
	UpdatedBy     int    `gorm:"column:updated_by;type:int(10) unsigned;not null"`        // 更新人
}

func (a AreaAdmin) TableName() string {
	return "cu_area_admin"
}
