package components

import (
	"github.com/gollerxiong/gbox/components/interfaces"
	"strings"
)

type BaseFormatter struct {
	data          map[string]interface{}
	Model         *interfaces.ModelInterface
	ColumnFuncMap map[string]func(interface{}) interface{}
}

func (b *BaseFormatter) SetData(data map[string]interface{}) interface{} {
	b.data = data
	return b
}

func (b *BaseFormatter) GetData() map[string]interface{} {
	return b.data
}

func (b *BaseFormatter) SetModel(model *interfaces.ModelInterface) interface{} {
	b.Model = model
	return b
}

func (b *BaseFormatter) Formate() map[string]interface{} {
	var result = make(map[string]interface{})
	for key, val := range b.GetData() {
		tmp := b.ColumnFormate(key, val)
		result[key] = tmp
	}

	b.SetData(result)
	return b.GetData()
}

func (b *BaseFormatter) ColumnFormate(key string, value interface{}) interface{} {
	key = strings.ToLower(key)
	callback, ok := b.ColumnFuncMap[key]

	if ok {
		return callback(value)
	} else {
		return value
	}
}
