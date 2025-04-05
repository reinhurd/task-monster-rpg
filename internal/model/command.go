package model

const (
	HELP            = "/help"
	CREATE_TASK_GPT = "/create_task_gpt"
	CONNECT_USER    = "/connect_user"
	TASK_LIST       = "/task_list"
	CREATE_TASK     = "/create_task"
	UPDATE_TASK     = "/edit_task"
	CREATE_USER     = "/create_user"
	START           = "/start"
	VIEW_TASK       = "/view_task"
)

var Commands = map[string]string{
	HELP: "/create_task_gpt <request> - create a task by your request to learn\n" +
		"/task_list - get list of tasks\n" +
		"/create_task <task_goal> <task_description> - create a task without GPT\n" + // todo fix
		"/profile - get your profile\n",
	CREATE_TASK_GPT: "Task created: %v",
	CONNECT_USER:    "User connected with userID: %v",
	TASK_LIST:       "Task: %v\n /view_task %v - view a task\n", //todo find how make a link in tg // and also sort tasks
	CREATE_TASK:     "Task created: %v",
	UPDATE_TASK:     "Task updated: %v",
	CREATE_USER:     "User created with ID: %v",
	VIEW_TASK:       "Task: %v\n /edit_task %v some params - edit a task\n", //todo find how make a link in tg
}
