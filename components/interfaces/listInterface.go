package interfaces

type ListInterface interface {
	SetField(string) *ListInterface
	SetPage(int64) *ListInterface
	SetPagenate(bool) *ListInterface
	SetPageSize(int64) *ListInterface
	SetOrder(string) *ListInterface
	SetIds([]int64) *ListInterface
	BuildParams()
	Filter()
	Find()
	AddField(string) *ListInterface
	GetFormatter()
}
