package model

const (
	HELP            = "/help"
	CREATE_TASK_GPT = "/create_task_gpt"
	CONNECT_USER    = "/connect_user"
)

var Commands = map[string]string{
	HELP: "/create_task_gpt <request> - create a task by your request to learn\n" +
		"/connect_user <user_ID> - connect user to telegram\n" +
		"/set_csv <filename> - set default csv\n",
	CREATE_TASK_GPT: "Task created: %v",
	CONNECT_USER:    "User connected: %v",
}
