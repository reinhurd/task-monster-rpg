package model

const (
	HELP            = "/help"
	CREATE_TASK_GPT = "/create_task_gpt"
)

var Commands = map[string]string{
	HELP: "/create_task_gpt <request> - create a task by your request to learn\n" +
		"/set_csv <filename> - set default csv\n",
	CREATE_TASK_GPT: "Task created: %v",
}
