package models

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type UserInfoModel struct {
	ID               int64               `json:"id" gorm:"primarykey"`
	UserCode         string              `json:"user_code" gorm:"not null;default:'';comment:用户码"`
	Nickname         string              `json:"nickname" gorm:"not null;default:'';comment:用户昵称"`
	Birthday         sql.NullTime        `json:"birthday" gorm:"type:date;not null;default:'';comment:用户生日"`
	Sex              int64               `json:"sex" gorm:"not null;default:0;comment:性别 0.未知 1.男 2.女"`
	Icon             string              `json:"icon" gorm:"not null;default:'';comment:用户头像"`
	RealIcon         string              `json:"real_icon" gorm:"not null;default:'';comment:用户真人头像"`
	Job              string              `json:"job" gorm:"not null;default:'';comment:用户工作"`
	Weight           string              `json:"weight" gorm:"not null;default:'';comment:体重"`
	Height           string              `json:"height" gorm:"not null;default:'';comment:身高"`
	Income           string              `json:"income" gorm:"not null;default:'';comment:收入"`
	Hometown         string              `json:"hometown" gorm:"not null;default:'';comment:家乡"`
	Desc             string              `json:"desc" gorm:"not null;default:'';comment:用户个人简介"`
	IsShow           int64               `json:"is_show" gorm:"not null;default:0;comment:是否展示开关 0.不展示 1.展示"`
	Status           int64               `json:"status" gorm:"not null;default:1;comment:用户状态：1正常 2禁用 3禁言 4注销"`
	City             string              `json:"city" gorm:"not null;default:'';comment:城市"`
	Province         string              `json:"province" gorm:"not null;default:'';comment:省份"`
	IsComplete       int64               `json:"is_complete" gorm:"not null;default:0;comment:是否展示开关 0.未完善 1.已完善"`
	IsReal           int64               `json:"is_real" gorm:"not null;default:0;comment:是否展示开关 0.未真人认证 1.已经真人认证"`
	IsSelf           int64               `json:"is_self" gorm:"not null;default:0;comment:是否展示开关 0.未实名认证 1.已经实名认证"`
	IsHighlyEducated int64               `json:"is_highly_educated" gorm:"not null;default:0;comment:是否展示开关 0.未做学历认证 1.已经做了学历认证"`
	IsStudent        int64               `json:"is_student" gorm:"not null;default:0;comment:是否展示开关 0.已经做了学生认证 1.没有做学生认证"`
	IsNymph          int64               `json:"is_nymph" gorm:"not null;default:0;comment:是否展示开关 0.没有做女神认证 1.做了女神认证"`
	FromChannel      string              `json:"from_channel" gorm:"not null;default:'';comment:来源渠道"`
	FromVersion      string              `json:"from_version" gorm:"not null;default:'';comment:版本号"`
	Recommend        int64               `json:"recommend" gorm:"not null;default:0;comment:推荐值"`
	OnlineStatus     int64               `json:"online_status" gorm:"not null;default:0;comment:在线状态 0.不在线 1.在线"`
	IsVip            int64               `json:"is_vip" gorm:"not null;default:0;comment:在线状态 0.非vip 1.vip"`
	LocationSource   int64               `json:"location_source" gorm:"not null;default:0;comment:在线状态 0.客户端 1.ip"`
	IsShowFromType   int64               `json:"is_show_from_type" gorm:"not null;default:0;comment:is_show字段修改来源"`
	EmotionalState   string              `json:"emotional_state" gorm:"not null;default:'';comment:用户情感状态"`
	UserTag          string              `json:"user_tag" gorm:"not null;default:'';comment:用户标签"`
	CreatedAt        time.Time           `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt        time.Time           `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	Connect          *gorm.DB            `json:"-" gorm:"-"`
	UserInfoExtend   UserInfoExtendModel `json:"user_info_extend" gorm:"foreignKey:user_id"`
}

func NewUserInfoModel() *UserInfoModel {
	// 这里的链接从连接池里面获取
	dsn := "root:123456@tcp(127.0.0.1:3306)/app_api?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	ins := &UserInfoModel{}
	ins.SetConnect(db)
	return ins
}

func (t *UserInfoModel) GetTableName() string {
	return "user_info"
}

func (t *UserInfoModel) GetId() int64 {
	return t.ID
}

func (t *UserInfoModel) GetConnectName() string {
	return "default"
}

func (t *UserInfoModel) GetConnect() *gorm.DB {
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
func (t *UserInfoModel) SetConnect(connect *gorm.DB) {
	t.Connect = connect
}
