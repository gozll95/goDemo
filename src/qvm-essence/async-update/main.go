	// set up task manager
	controller.Setup(app.Config.Task)

	var TaskManager *task.TaskManager

func Setup(taskConfig *conf.TaskConfig) {
	TaskManager = task.NewTaskManager(taskConfig.Worker)
	TaskManager.Start()
	utils.StdLog.Info("started task")
}