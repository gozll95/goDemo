//types
type RunMode string

func (mode RunMode) IsValid() bool {
	switch mode {
	case Development, Test, Production:
		return true
	}

	return false
}


// EvmOpType
package enums


// 资源操作类型
type EvmOpType int

const (
	_evmOpType EvmOpType = iota
	EvmOpCreate
	EvmOpUpdate
	EvmOpShutdown
	EvmOpDelete
	EvmOpStop
	EvmOpReboot
	EvmOpStart
	EvmOpUpdateCostMode
	evmOpType_
)

// 计费操作类型
type EvmCostOpType int

const (
	_evmCostOpType   EvmCostOpType = iota
	EvmCostOpDegrade               // 配置更新为降低配置时，对应的对计费的操作类型
	EvmCostOpRenew                 // 续费操作
	EvmCostOpUpgrade               // 配置更新为升高配置时，对应的对计费的操作类型
	evmCostOpType_
)

// 资源类型
type EvmResourceType int

const (
	_evmResourceType EvmResourceType = iota
	EvmResourceComputer
	EvmResourceVolume
	EvmResourceFloatingip
	EvmResourceSnapshot
	EvmResourceListener
	EvmResourceImage
	EvmResourceLoadbalancer
	EvmResourcePool
	EvmResourceSecurityGroup
	EvmResourceSecurityRule
	EvmResourceKey
	EvmResourceLoadbalancerV2
	evmResourceType_
)


//类方法:
Valid()
IsValid()
IsCreate()
IsUpdate()
... 
Humanize()string 
String()string


// 提供方

type Provider string


const (
	OpenStack Provider = "openstack"
	Ustack    Provider = "ustack"
	T2cloud   Provider = "t2cloud"
)

//类方法:
IsValid()
func (p Provider) Name() string {
	return string(p)
}


//计费类型

type EvmCostModeType int

const (
	_evmCostModeType EvmCostModeType = iota
	EvmCostModeHourly
	EvmCostModeMonthly
	EvmCostModeYearly
	evmCostModeType_
)


main.go
- 输入 runMode
- 输入 srcPath
- // verify run mode
- if mode := gogo.RunMode(runMode); !mode.IsValid()
- // adjust src path
- srcPath=... 
- evmapp := app.New(runMode, srcPath) 
- evmlogger := evmapp.Logger()
- evmlogger := evmapp.Logger()
- // oplog
	id := ""
	limit := 100
	for {
		evmlogger.Infof("Starting adjusting cost mode from %s with %d ...", id, limit)

		oplogs, err := models.Oplog.AllByID(id, limit)
		if err != nil {
			evmlogger.Errorf("Oplog.AllByID(%s, %d): %v", id, limit, err)
			return
		}

		for _, oplog := range oplogs {
			if strings.ContainsAny(oplog.Args, "%") {
				evmlogger.Printf("%#v", oplog)

				matches := ri.FindAllStringSubmatch(oplog.Args, -1)
				if len(matches) > 0 {
					oplog.Args = ri.ReplaceAllStringFunc(oplog.Args, func(matched string) string {
						for i := 0; i < len(matches); i++ {
							if matches[i][0] == matched {
								return matches[i][1]
							}
						}

						return matched
					})
				}

				matches = rs.FindAllStringSubmatch(oplog.Args, -1)
				if len(matches) > 0 {
					oplog.Args = rs.ReplaceAllStringFunc(oplog.Args, func(matched string) string {
						for i := 0; i < len(matches); i++ {
							if matches[i][0] == matched {
								return matches[i][1]
							}
						}

						return matched
					})
				}

				models.Oplog.Query(func(c *mgo.Collection) {
					err := c.UpdateId(oplog.Id, bson.M{
						"$set": bson.M{
							"args": oplog.Args,
						},
					})
					if err != nil {
						evmlogger.Errorf(">>> %#v : %v", oplog, err)
					}
				})
			}
		}

		if len(oplogs) < limit {
			break
		}

		id = oplogs[len(oplogs)-1].Id.Hex()
	}


应用到controller:
AuditlogHelper.GenerateAuditlog(requester, evmEnums.AuditCreate, vmId, "", enums.EvmResourceComputer, enums.EvmResourceComputer, request, ctx.Logger)
OplogHelper.GenerateCostOplog(requester, evmEnums.EvmOpCreate, vmId, enums.EvmResourceComputer, request, ctx.Logger)
		- func (_ *_OplogHelper) GenerateCostOplog(requester RequesterInfoer,op evmEnums.EvmOpType, resourceId string, resourceType enums.EvmResourceType, param CostParamer, logger gogo.Logger) error {
					- // encoding resource attrs
					- json.Marshal(param)-->string-->args 
					- //get some info
					- uid := requester.UserID()
					- platformRegion := requester.CurrentRegion()
					- platfromProvider := requester.CurrentPlatformProvider()
					- // new oplog model
					- oplog := models.NewOplogModel(uid, platformRegion, platfromProvider)
					- // set resource infp
					- oplog.WithResourceInfo(op, resourceId, resourceType, args)
					- // set cost info
					- oplog.WithCostInfo(param.GetCostMode(), param.GetAmount())
					- // save oplog
					- err=oplog.Save()
					- // set cost info
					- err = CostHelper.SetCostInfo(uid, resourceId, resourceType, param.GetCostMode(), param.GetAmount(), param.GetAutoRenew())



oplog_helper.go
//types:
type _OplogHelper struct{}

var (
	OplogHelper *_OplogHelper
)

type RequesterInfoer interface {
	UserID() uint32
	CurrentRegion() string
	CurrentPlatformProvider() platform.Provider
}

type CostParamer interface {
	GetCostMode() enums.EvmCostModeType
	GetAmount() int
	GetAutoRenew() *bool
}
//类方法:

func (_ *_OplogHelper) GenerateCostOplog(requester RequesterInfoer,
	op evmEnums.EvmOpType, resourceId string, resourceType enums.EvmResourceType, param CostParamer, logger gogo.Logger) error

func (_ *_OplogHelper) GenerateOplog(requester RequesterInfoer,
	op evmEnums.EvmOpType, resourceId string, resourceType enums.EvmResourceType, param interface{}, logger gogo.Logger)

// 貌似是清理没有的oplog?
func (_ *_OplogHelper) FilterExistOplogs(oplogs []*models.OplogModel) (resOplogs []*models.OplogModel,
	err error)




