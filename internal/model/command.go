package model

const (
	HELP            = "/help"
	CREATE_TASK_GPT = "/create_task_gpt"
	CONNECT_USER    = "/connect_user"
	TASK_LIST       = "/task_list"
	CREATE_TASK     = "/create_non_gpt_task"
	UPDATE_TASK     = "/update_task"
	CREATE_USER     = "/create_user"
	START           = "/start"
)

var Commands = map[string]string{
	HELP: "/create_task_gpt <request> - create a task by your request to learn\n" +
		"/connect_user <login> <password> - connect user to telegram\n" +
		"/task_list - get list of tasks\n" +
		"/create_non_gpt_task <task_goal> <task_description> - create a task without GPT\n" +
		"/update_task <task_id> <task_goal> <task_description> - update a task\n" + //todo add additional field for task
		"/create_user <login> <password> - create a user\n",
	CREATE_TASK_GPT: "Task created: %v",
	CONNECT_USER:    "User connected with userID: %v",
	TASK_LIST:       "Task: %v\n",
	CREATE_TASK:     "Task created: %v",
	UPDATE_TASK:     "Task updated: %v",
	CREATE_USER:     "User created with ID: %v",
}
