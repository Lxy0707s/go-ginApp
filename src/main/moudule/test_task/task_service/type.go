package task_service

type QueryArgs struct {
	TaskName string `json:"task_name,omitempty"`
	TaskId   string `json:"task_id,omitempty"`
}
