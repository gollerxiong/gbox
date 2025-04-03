package userInfoDao

// helper抽象到单独的目录中，主要解决go不能互相引入的问题
type userInfoHelper struct {
}

func NewUserInfoHelper() *userInfoHelper {
	return &userInfoHelper{}
}
