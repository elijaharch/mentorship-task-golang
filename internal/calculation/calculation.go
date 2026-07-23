package calculation

import "time"

type Operation string

const (
	OperationAdd      Operation = "+"
	OperationSubtract Operation = "-"
	OperationMultiply Operation = "*"
	OperationDivide   Operation = "/"
)

type Input struct {
	A         float64
	B         float64
	Operation Operation
}

type Calculation struct {
	ID        int64
	A         float64
	B         float64
	Operation Operation
	Result    float64
	CreatedAt time.Time
}

type ListOptions struct {
	Limit  int
	Offset int
}

func (o Operation) Valid() bool {
	switch o {
	case OperationAdd,
		OperationSubtract,
		OperationMultiply,
		OperationDivide:
		return true
	default:
		return false
	}
}
