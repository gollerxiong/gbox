package interfaces

type ObjectInterface interface {
	// 获取方法集合
	GetModel() interface{}
	GetFormatter() interface{}
	GetAttributes() map[string]interface{}
	GetOldAttributes() map[string]interface{}

	// 操作方法集合
	IsNew() bool
	Save() bool
	Delete() bool

	// 加载方法集合
	LoadById(int64) interface{}
	DeleteById(int64) bool

	// 设置方法集合
	SetField(string) *ObjectInterface
	SetAttributes(map[string]interface{}) *ObjectInterface
	SetParams(map[string]interface{}) *ObjectInterface
	SetFormatter(*ObjectInterface) *ObjectInterface
	//SetValues(map[string]interface{}) *BaseObjectInterface
	SetModel(*ObjectInterface) *ObjectInterface
	SetOldAttributes(map[string]interface{}) *ObjectInterface
}
