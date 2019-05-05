	projectList, _ := service.ProjectService.GetAllProject()
	permList := service.SystemService.GetPermList()

	chkmap := make(map[string]string)
	for _, v := range role.PermList {
		chkmap[v.Key] = "checked"
	}
	if role.ProjectIds != "" {
		pids := strings.Split(role.ProjectIds, ",")
		for _, pid := range pids {
			chkmap[pid] = "checked"
		}
	}