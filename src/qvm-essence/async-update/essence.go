
1."思路"
type TaskManager struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Worker chan bool

	Closed chan bool

	NoLeftLock sync.RWMutex
	IsNoLeft   bool
}


2."TaskManager"
type TaskManager struct {
	ID     string `json:"id"`
	Action string `json:"action"`
	Worker chan bool

	Closed chan bool

	NoLeftLock sync.RWMutex
	IsNoLeft   bool
}


func NewTaskManager(worker int) *TaskManager {
	return &TaskManager{
		Worker: make(chan bool, worker),
		Closed: make(chan bool),
	}
}


func (t *TaskManager) Start() {
	go t.run()
	go t.refreshTicket()

}

- "TaskManager的run()方法"
	- defer func(){
		...  
		t.run()
	}()
	- for{
		select{
		case <-t.Closed:
		case t.Worker <-true: //相当于取得票 然后开启goroutine
			- go func(){
				...  
				// check 如果 is no left 就 sleep 3 s
				isNoLeft := t.execute()
				... 
			}()
		}
	}

- "TaskManager的execute()方法"
	// 在ticket里寻找最近的一个空闲的ticket并lock它,这里采用了mgo的find and modify 原子操作
	- ticket, err := model.Ticket.FindAndLock()
	// 在task里找ticket对应的那一类资源所有的task // 所以这里ticket是task的一个集合
	- tasks, err := model.Task.AllByUidRegionAndResourceType(ticket.Uid, ticket.RegionId, ticket.ResourceType)
	- // 当tasks长度为0 就 ticket.Remove()
	- if len(tasks) == 0 -> err = ticket.Remove()
	- // 执行每个task对应async func
	- checkers, err := t.do(ticket.Uid, ticket.RegionId, ticket.ResourceType, resourceIds)
	- // 一些判断
	- if task.OnlyOnce -> err = task.Remove()
	- if task.TargetStatus != "" && task.TargetStatus == checker.ResourceStatus() -> err = task.Remove()
	- ... 


- "TaskManager的refreshTicket()方法"
	- defer func(){
		... 
		t.refreshTicket()
	}()
	- for{
		select{
			case <-t.Closed: 
			default: 
				- time.Sleep(time.Second * 9)
				- // 在task里寻找更新时间大于18s的游标 iter
				- iter := model.Task.FindUpdatedIter()
				- // iter next
				- task, isError := model.Task.Next(iter)
				-  for isError {
					_, err := model.Ticket.FindbyRegionIdAndResourcetype(task.Uid, task.RegionId, task.ResourceType)
					if err == model.ErrNotFound {
						// Note: only no ticket save the ticket
						// or update ticket updatedAt time will effect ticket wait time
							ticket := model.NewTicketModel(task.Uid, task.RegionId, task.ResourceType)
						err = ticket.Save()
					}

					if err != nil {
						utils.StdLog.Errorf("refresh ticket find or save error %v", err)
					}
					
					task, isError = model.Task.Next(iter)
			}
		}
	}

// 创建task通过tasker
- "TaskManager的NewTaskByTasker(tasker Tasker)方法"
	- func (t *TaskManager) NewTaskByTasker(tasker Tasker) (err error)
		- ... 
		- return t.NewStatusTask(uid, regionId, resourceType, resourceId, "")

- "TaskManager的NewTaskByTaskers(taskers ...Tasker)方法"
	func (t *TaskManager) NewTaskByTaskers(taskers ...Tasker) (err error) 

- "TaskManager的NewTaskByTaskers(taskers ...Tasker)方法"

// 创建 状态的 task
func (t *TaskManager) NewStatusTask(uid uint32, regionId string, resourceType enums.ResourceType, resourceId, targetStatus string) (err error) {
	task, err := model.Task.FindOrCreate(uid, regionId, resourceType, resourceId)
	if err != nil {
		return
	}

	task.SetTargetStatus(targetStatus)
	err = task.Save()

	if err != nil {
		return
	}

	_, err = model.Ticket.FindbyRegionIdAndResourcetype(uid, regionId, resourceType)
	if err == model.ErrNotFound {
		// Note: only no ticket save the ticket
		// or update ticket updatedAt time will effect ticket wait time
		ticket := model.NewTicketModel(uid, regionId, resourceType)
		return ticket.Save()
	}

	return err

}


3."task"
type Tasker interface {
	TaskUid() uint32
	TaskRegion() string
	TaskResourceType() enums.ResourceType
	TaskResourceId() string
}



4."TaskHelper":

- methods:
	- "NewEipTasks"
	- "NewDiskTasks"
	- "NewSnapshotTasks"
	- ... 

func (_ *_TaskHelper) NewEipTasks(logger utils.Log, tasks ...*model.EipModel) {
	for _, tasker := range tasks {
		err := TaskManager.NewTaskByTasker(tasker)

		if err != nil {
			logger.Errorf("New Task  %#v error %v", tasker, err)
		}
	}
}


5."checker"


func eipUpdate(uid uint32, regionId string, resourceType enums.ResourceType, resourceIds []string) (checkers []Checker, err error) {
	var errMsg bool

	utils.StdLog.Debug("start do enums.ResourceTypeIp")

	//make ecs client
	ecsClient, err := newEcs(uid, regionId)
	if err != nil {
		utils.StdLog.Errorf("generateEcs(%v, %v):%v", uid, regionId, err)
		return nil, err
	}

	// call ali describe sdk
	eipEcsClient := ecsClient.NewEip()

	eipAllocationIds := strings.Join(resourceIds, ",")
	describeEipAddressesArgs := params.DescribeEipAddresses{
		RegionId:     regionId,
		AllocationId: &eipAllocationIds,
	}

	desRes, err := eipEcsClient.DescribeEipAddresses(&describeEipAddressesArgs)
	if err != nil {
		utils.StdLog.Errorf("Ecs(?).DescribeEipAddresses(%#v):%v", describeEipAddressesArgs, err)
		return nil, err
	}

	switch desRes.TotalCount {
	case 0:
		return nil, nil
	default:
		for _, item := range desRes.EipAddresses.EipAddress {
			//find from db
			dbEip, err := model.Eip.FindByUidRegionIdAndAllocationId(uid, regionId, item.AllocationId)
			if err != nil {
				errMsg = true
				utils.StdLog.Errorf("model.Eip.FindByUidRegionIdAndAllocationId(%v,%v,%v):%v", uid, regionId, item.AllocationId, err)
				continue
			}
			// compare db with ali
			isUpdated, err := compareAndUpdateEip(dbEip, item)
			if err != nil {
				errMsg = true
				utils.StdLog.Errorf("compareAndUpdateEip(%v,%v):%v", dbEip, item, err)
				continue
			}
			if isUpdated {
				dbEip.SetUpdate(true)
				utils.StdLog.Debugf("update %v successfully", dbEip)
			}
			checkers = append(checkers, dbEip)
		}
	}
	if errMsg {
		return checkers, errors.New(fmt.Sprintf("something wrong in eipUpdate(%v,%v,%v,%v):%v", uid, regionId, resourceType, resourceIds, err))
	}
	return
}

func compareAndUpdateEip(dbEip *model.EipModel, aliEip params.EipAddressSetType) (isUpdated bool, err error) {
	isSame, err := compareEip(dbEip, aliEip)
	if err != nil {
		utils.StdLog.Errorf("compareEip(%v, %v):%v", dbEip, aliEip, err)
		return
	}
	switch isSame {
	case false:
		err = convertEip(dbEip, aliEip)
		if err != nil {
			utils.StdLog.Errorf("convertEip(%v,%v):%v", dbEip, aliEip, err)
			return
		}
		err = dbEip.Save()
		if err != nil {
			utils.StdLog.Errorf("dbEip.Save():%v", err)
			return
		}
		isUpdated = true
		return
	}
	return
}

func compareEip(dbEip *model.EipModel, aliEip params.EipAddressSetType) (isSame bool, err error) {
	var aliExpiredTime, aliAllocationTime time.Time

	isSame, err = utils.QDeepEqual(dbEip.OperationLocks, aliEip.OperationLocks)
	if err != nil || !isSame {
		return
	}

	if aliEip.ExpiredTime != "" {
		aliExpiredTime, err = utils.ConvertTime("2006-01-02T15:04:05Z", aliEip.ExpiredTime)
		if err != nil {
			utils.StdLog.Errorf("utils.ConvertTime(%v):%v", aliEip.ExpiredTime)
			return isSame, err
		}
	}

	if aliEip.AllocationTime != "" {
		aliAllocationTime, err = utils.ConvertTime("2006-01-02T15:04:05Z", aliEip.AllocationTime)
		if err != nil {
			utils.StdLog.Errorf("utils.ConvertTime(%v):%v", aliEip.AllocationTime)
			return isSame, err
		}
	}

	if dbEip.InternetChargeType == aliEip.InternetChargeType &&
		dbEip.ChargeType == aliEip.ChargeType &&
		dbEip.Status == aliEip.Status &&
		dbEip.InstanceId == aliEip.InstanceId &&
		dbEip.InstanceType == aliEip.InstanceType &&
		dbEip.Bandwidth == aliEip.Bandwidth &&
		dbEip.AliExpiredTime.Equal(aliExpiredTime) &&
		dbEip.AliAllocationTime.Equal(aliAllocationTime) {

		isSame = true
		return
	}

	isSame = false

	return
}

func convertEip(dbEip *model.EipModel, aliEip params.EipAddressSetType) (err error) {

	err = utils.CopyStruct(dbEip, aliEip)
	if err != nil {
		utils.StdLog.Errorf("utils.CopyStruct(%v,%v):%v", dbEip, aliEip)
		return
	}

	if aliEip.ExpiredTime != "" {
		dbEip.AliExpiredTime, err = utils.ConvertTime("2006-01-02T15:04:05Z", aliEip.ExpiredTime)
		if err != nil {
			utils.StdLog.Errorf("utils.ConvertTime(%v):%v", aliEip.ExpiredTime, err)
			return
		}
	}

	if aliEip.AllocationTime != "" {
		dbEip.AliAllocationTime, err = utils.ConvertTime("2006-01-02T15:04:05Z", aliEip.AllocationTime)
		if err != nil {
			utils.StdLog.Errorf("utils.ConvertTime(%v):%v", aliEip.AllocationTime, err)
			return
		}
	}

	return
}



2."流程"
修改eip:
	- TaskHelper.NewEipTasks(indexReq.Logger(), eips...)






5."可以修正":
	for {
		select {
		case <-t.Closed:
			panic(closed)
		case t.Worker <- true:
			go func() {
				// check No Left true
				t.NoLeftLock.RLock()

				// if no left task sleep 3 second and release no left
				if t.IsNoLeft {
					t.NoLeftLock.RUnlock()

					time.Sleep(time.Second * 3)

					t.NoLeftLock.Lock()
					t.IsNoLeft = false
					t.NoLeftLock.Unlock()
				} else {
					t.NoLeftLock.RUnlock()
				}

				isNoLeft := t.execute()

				if isNoLeft {
					// set no left marker
					t.NoLeftLock.Lock()
					t.IsNoLeft = true
					t.NoLeftLock.Unlock()
				}

				<-t.Worker
			}()

		}

	}