前端:

type ResourceMessageRequest struct {
 +	Resources []*ResourceMessageInfo `json:"resources"`
 +}
 +
 +type ResourceMessageInfo struct {
 +	RegionId     string             `json:"region_id"`
 +	ResourceType enums.ResourceType `json:"resource_type"`
 +	ResourceId   string             `json:"resource_id"`
 +	UpdatedAt    time.Time          `json:"updated_at"`
 +	IsUpdated    bool               `json:"is_updated"`
 +}


request:
	vm 1  updateAt: 3点
	vm 2  updateAt: 2点
	vm 3
	eip 1
	eip 2


db:  vm  1 - msg  updateAt: 1点  --> resource.IsUpdated = true
	 vm   2 - msg 



db:
- ticker:
	+type TicketModel struct {
	 +	Id           bson.ObjectId          `bson:"_id" json:"_"`
	 +	Uid          uint32                 `bson:"uid" json:"uid"`
	 +	RegionId     string                 `bson:"region_id" json:"region_id"`
	 +	ResourceType enums.ResourceType     `bson:"resource_type" json:"resource_type"`
	 +	Status       enums.TicketStatusType `bson:"status" json:"status"`
	 +
	 +	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	 +	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	 +
	 +	isNewRecord bool `bson:"-" json:"-"`
	 +}

用户A	cn-qingdao		EIP		
用户A	cn-qingdao		EIP		
用户A	cn-qingdao		EIP		

- task
	+type TaskModel struct {
	 +	Id           bson.ObjectId      `bson:"_id" json:"_"`
	 +	Uid          uint32             `bson:"uid" json:"uid"`
	 +	RegionId     string             `bson:"region_id" json:"region_id"`
	 +	ResourceType enums.ResourceType `bson:"resource_type" json:"resource_type"`
	 +	ResourceId   string             `bson:"resource_id" json:"resource_id"`
	 +	TargetStatus string             `bson:"target_status" json:"target_status"`
	 +
	 +	Timeout   int       `bson:"timeout" json:"timeout"` // default unit second
	 +	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	 +	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	 +
	 +	isNewRecord bool `bson:"-" json:"-"`
	 +}


用户A	cn-qingdao		EIP		状态A(这里应该是最终状态)
用户A	cn-qingdao		EIP		状态B
用户A	cn-qingdao		EIP		状态C


- ResourceMessageModel
+type ResourceMessageModel struct {
 +	Id           bson.ObjectId      `bson:"_id" json:"_"`
 +	Uid          uint32             `bson:"uid" json:"uid"`
 +	RegionId     string             `bson:"region_id" json:"region_id"`
 +	ResourceType enums.ResourceType `bson:"resource_type" json:"resource_type"`
 +	ResourceId   string             `bson:"resource_id" json:"resource_id"`
 +
 +	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
 +	CreatedAt time.Time `bson:"created_at" json:"created_at"`
 +
 +	isNewRecord bool `bson:"-" json:"-"`
 +}




tick_type.go
idle/busy


checker.go 
type Checker interface {
 ID() string
 Status() string
 Updated() bool
}



task_do.go 
do(uid,regionId,resouceType,resouceId []string)(checkers []Checker, err error)



type TaskManager struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Worker chan bool

	NoLeftLock sync.RWMutex
	IsNoLeft   bool
}


- NewTaskManager(worker int) *TaskManager

//新建没有目标状态的task
- func (t *TaskManager) NewTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId string) (err error) 
	- return t.NewStatusTask(uid, regionId, resourceType, resourceId, "")
		- func (t *TaskManager) NewStatusTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId, targetStatus string) (err error)
			- // 去db里查找是否有这个task对应的条目 如果没有 就添加到db
			- task, err := model.Task.FindOrCreate(uid, regionId, resourceType, resourceId) 
			- //设置目标状态
			- task.SetTargetStatus(targetStatus)
					- func (task *TaskModel) SetTargetStatus(status string)
						- task.TargetStatus = status
			- err = task.Save()
			- // 在 ticket 库 里 寻找 userA、cn-qingdao、EIP
			- _, err = model.Ticket.FindbyRegionIdAndResourcetype(uid, regionId, resourceType)
			- if err == model.ErrNotFound {
				- ticket := model.NewTicketModel(uid, regionId, resourceType)
				- return ticket.Save()

//新建有目标状态的task
- func (t *TaskManager) NewStatusTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId, targetStatus string) (err error) {


- func (t *TaskManager) run()
	- defer func(){ t.run() }()
	- for 循环
		- t.Worker<-true  //这里t.Worker是缓冲的
		- go routine
			- t.NoLeftLock.Rlock()
			- if t.IsNoLeft
					- time.Sleep(time.Second * 3)
					- t.NoLeftLock.Lock()
					- t.IsNoLeft = false
					- t.NoLeftLock.Unlock()
			- t.NoLeftLock.RUnlock()
			- isNoLeft := t.execute()
			- if isNoLeft
				- t.NoLeftLock.Lock()
				- t.IsNoLeft = true
				- t.NoLeftLock.Unlock()
			- <-t.Worker

- func(t *TaskManager) refreshTicke()
	- defer func(){ t.run() }()
	- for 循环
		- time.Sleep(time.Second * 9)
		- //在task表里查找 比现在 小 18 s 的 tasks
		- iter := model.Task.FindUpdatedIter()
		- for task, isError := model.Task.Next(iter); isError; {
				// 在ticket里寻找 userA cn-qingdao EIP
			- _, err := model.Ticket.FindbyRegionIdAndResourcetype(task.Uid, task.RegionId, task.ResourceType)
			- //如果没有就创建这个ticket
			- ... 

- func(t *TaskManager)execute()bool
	- // 在ticket表里寻找 距离现在小于3s 且 状态 为idel的条目 并且更新 busy + update_at
	- ticket, err := model.Ticket.FindAndLock()
	- //如果 在 ticket表里 距离现在 小于3s 的 条目没找到,就 回 true
	- if err == model.ErrNotFound  - return true 
	- // 在 task表里寻找 userA cn-qingdao EIP  -> 有多个条目
	- tasks, err := model.Task.AllByUidRegionAndResourceType(ticket.Uid, ticket.RegionId, ticket.ResourceType)
	- // 如果 找不到 tasks 就 remove对应的 ticket
	- if len(tasks) == 0 - ticket.Remove()
	- // 合并 tasks
	- ... 
	- checkers, err := t.do(ticket.Uid, ticket.RegionId, ticket.ResourceType, resourceIds)
	- taskMap := make(map[string]*model.TaskModel)
	- for _,task:=range tasks -- taskMap[task.ResourceId] = task
	- now:=time.Now()
	- 遍历 checkers
		- task, ok := taskMap[checker.ID()]
		- if checker.Updated()
			- // 在 ResourceMessage 表里 find or create
			- msg, err := model.ResourceMessage.FindOrCreate(task.Uid, task.RegionId, task.ResourceType, task.ResourceId)
			- err = msg.Save()
		- //如果 task.TargetStatus != "" && task.TargetStatus == checker.Status() 
		- if task.TargetStatus != "" && task.TargetStatus == checker.Status() - task.Remove()
		- ... -> task.Remove()
		- timeout --> task.Remove()
		- ticket.Unlock()
				- //重新回到idel状态
				- ticketModel.Status = enums.TicketStatusTypeIdle

