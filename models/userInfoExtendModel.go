package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// gorm:"table:user_info_extension"
type UserInfoExtendModel struct {
	ID                int64     `json:"id" gorm:"primarykey"`
	UserId            int64     `json:"user_id" gorm:"not null;comment:用户id"`
	FamilyId          int64     `json:"family_id" gorm:"not null;default:0;comment:家族id"`
	InviteCode        string    `json:"invite_code" gorm:"not null;default:'';comment:邀请码"`
	InviteUserId      int64     `json:"invite_user_id" gorm:"not null;default:0;comment:邀请人用户id"`
	RegisterUmengId   string    `json:"register_umeng_id" gorm:"not null;default:'';comment:一个用户只能存在唯一的设备id如果重复了，就保存第一个注册的用户id"`
	Token             string    `json:"token" gorm:"not null;default:0;comment:用户登陆态token"`
	LastLoginIp       string    `json:"last_login_ip" gorm:"not null;default:'';comment:用户最后登录的ip"`
	RegisterIp        string    `json:"register_ip" gorm:"not null;default:'';comment:用户注册ip"`
	ImageSelf         string    `json:"image_self" gorm:"not null;default:'';comment:真人认证过的照片"`
	Realname          string    `json:"realname" gorm:"not null;default:'';comment:用户真实姓名"`
	Identity          string    `json:"identity" gorm:"not null;default:'';comment:用户身份证号"`
	ImageIdentityDown string    `json:"image_identity_down" gorm:"not null;default:'';comment:用户身份证号反面"`
	ImageIdentityUp   string    `json:"image_identity_up" gorm:"not null;default:'';comment:用户身份证号正面"`
	Voice             string    `json:"voice" gorm:"not null;default:'';comment:用户的语音签名"`
	VoiceTime         int64     `json:"voice_time" gorm:"comment:语音签名时常"`
	RechargeLevel     int64     `json:"recharge_level" gorm:"not null;default:0;comment:充值等级"`
	InviteAt          int64     `json:"invite_at" gorm:"comment:邀请时间"`
	IsSuperInvite     int64     `json:"is_super_invite" gorm:"not null;default:0;comment:是否为超推"`
	SuperInviteTime   int64     `json:"super_invite_time"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Connect           *gorm.DB  `json:"-" gorm:"-"`
}

func NewUserInfoExtendModel() *UserInfoExtendModel {
	// 这里的链接从连接池里面获取
	dsn := "root:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	ins := &UserInfoExtendModel{}
	ins.SetConnect(db)
	return ins
}

func (t *UserInfoExtendModel) GetTableName() string {
	return "user_info"
}

func (t *UserInfoExtendModel) GetId() int64 {
	return t.ID
}

func (t *UserInfoExtendModel) GetConnectName() string {
	return "default"
}

func (t *UserInfoExtendModel) GetConnect() *gorm.DB {
	if t.Connect == nil {
		panic("table of " + t.GetTableName() + " connect is nil!")
	}

	// 指定表名
	return t.Connect.Table(t.GetTableName()).Debug()
}

/*
*
设置链接
*/
func (t *UserInfoExtendModel) SetConnect(connect *gorm.DB) {
	t.Connect = connect
}
