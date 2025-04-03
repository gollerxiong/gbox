package interfaces

type ListFormatterInterface interface {
	SetList([]ModelInterface) *ListFormatterInterface
	SetFields(string) *ListFormatterInterface
	Formate() []map[string]interface{}
	Do(int, ModelInterface)
}
