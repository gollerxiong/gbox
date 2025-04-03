package components

import (
	"github.com/gollerxiong/gbox/components/interfaces"
	"strings"
	"sync"
)

type BaseListFormatter struct {
	List      []interfaces.ModelInterface
	Fields    string
	WG        sync.WaitGroup
	Result    []map[string]interface{}
	Formatter interfaces.FormatterInterface
}

func (b *BaseListFormatter) SetList(list []interfaces.ModelInterface) *BaseListFormatter {
	b.List = list
	return b
}

func (b *BaseListFormatter) SetFields(fields string) *BaseListFormatter {
	b.Fields = fields
	return b
}

func (b *BaseListFormatter) Formate() []map[string]interface{} {
	var length = len(b.List)

	b.WG.Add(length)
	for i := 0; i < length; i++ {
		go b.Do(i, b.List[i])
	}

	b.WG.Wait()
	return b.Result
}

func (b *BaseListFormatter) Do(index int, item interfaces.ModelInterface) {
	defer b.WG.Done()
	result := make(map[string]interface{})
	arr := StructToMap(item)
	if !strings.Contains(b.Fields, "*") {
		fieldArr := strings.Split(b.Fields, ",")
		for _, key := range fieldArr {

			val, ok := arr[key]
			if ok {
				result[key] = b.Formatter.ColumnFormate(key, val)
			}
		}
		b.Result = append(b.Result, result)
	} else {
		b.Result = append(b.Result, arr)
	}
}
