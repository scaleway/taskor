package handler

// Metric contains all metrics available
type Metric struct {
	TaskSent            uint32
	TaskDoneWithSuccess uint32
	TaskDoneWithError   uint32
}
