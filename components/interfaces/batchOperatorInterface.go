package interfaces

type BatchOperatorInterface interface {
	GetModel() interface{}
	SetModel(*ModelInterface) *BatchOperatorInterface
	SetPage(int65 int64) *BatchOperatorInterface
	SetPageSize(int64) *BatchOperatorInterface
	SetField(string) *BatchOperatorInterface
	SetIds([]int64) *BatchOperatorInterface
	SetFieldValue(interface{}) *BatchOperatorInterface
	BuildParams()
}
