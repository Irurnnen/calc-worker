package models

type Task struct {
	ID             int     `json:"id"`
	FirstArgument  float64 `json:"arg1"`
	SecondArgument float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	OperationTime  int     `json:"operation_time"`
}

type Answer struct {
	ID     int     `json:"id"`
	Result float64 `json:"result"`
}
