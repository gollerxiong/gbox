package interfaces

type FormatterInterface interface {
	SetData(map[string]interface{}) interface{}
	GetData() map[string]interface{}
	ColumnFormate(key string, value interface{}) interface{}
	SetModel(*ModelInterface) interface{}
	Formate() map[string]interface{}
}
