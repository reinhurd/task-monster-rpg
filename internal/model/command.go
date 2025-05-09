package model

const (
	HELP              = "/help"
	CREATE_TASK_GPT   = "/create_task_gpt"
	UPDATE_TASK_DATE  = "/update_task_date"
	UPDATE_TASK_TITLE = "/update_task_title"
	UPDATE_TASK_DESC  = "/update_task_desc"
	CONNECT_USER      = "/connect_user"
	TASK_LIST         = "/task_list"
	CREATE_TASK       = "/create_task"
	UPDATE_TASK       = "/edit_task"
	CREATE_USER       = "/create_user"
	START             = "/start"
	VIEW_TASK         = "/view_task"
	PROFILE           = "/profile"
)

var Commands = map[string]string{
	HELP: CREATE_TASK_GPT + " <request> - create a task by your request to learn\n" +
		TASK_LIST + " - get list of tasks\n" +
		CREATE_TASK + " <task_goal> <task_description> - create a task without GPT\n" + // todo fix
		PROFILE + " - get your profile\n",
	CREATE_TASK_GPT: "Task created: %v",
	CONNECT_USER:    "User connected with userID: %v",
	TASK_LIST:       "Task description: %v\n Task ID: %v \n", //todo sort tasks by date
	CREATE_TASK:     "Task created: %v",
	UPDATE_TASK:     "Task updated: %v",
	CREATE_USER:     "User created with ID: %v",
	VIEW_TASK:       "Task ID: %v\n TaskUniqID: %v\n\n Title: %v\n Description: %v\n Status: %v\n Executor: %v\n Reviewer: %v\n Deadline: %v\n Created: %v\n Updated: %v\n",
}
