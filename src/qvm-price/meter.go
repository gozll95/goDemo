’meter 是 跟商业运营的计费项相关的


以eip为例:

      {
        "DataItems": {
          "DataItem": [
            {
              "Name": "NetworkOut",
              "Value": "0"
            },
            {
              "Name": "ProviderId",
              "Value": "1103909446200972"
            },
            {
              "Name": "EndTime",
              "Value": "2018-04-10T01:00:00Z"
            },
            {
              "Name": "NetworkIn",
              "Value": "0"
            },
            {
              "Name": "Bandwidth",
              "Value": "1048576"
            },
            {
              "Name": "StartTime",
              "Value": "2018-04-10T00:00:00Z"
            },
            {
              "Name": "UsedTime",
              "Value": "3600"
            },
            {
              "Name": "Region",
              "Value": "cn-beijing-btc-a01"
            },
            {
              "Name": "EipId",
              "Value": "eip-2ze3ncglp9z5rvtake4hx"
            },
            {
              "Name": "EIP",
              "Value": "47.95.229.115"
            }
          ]
        }


// -------------------------------------------------------
// -------------------------------------------------------

// model - 学会留存源数据
type AliEcsEipMeterModel struct {
	Id        bson.ObjectId `bson:"_id"`
	Uid       uint32        `bson:"uid"`
	RegionId  string        `bson:"region_id"`
	StartTime time.Time     `bson:"start_time"`
	EndTime   time.Time     `bson:"end_time"`
	Spec      string        `bson:"spec"` // 对于eip来说是 bandwidth

	meter.EipMetric `bson:"data"` 
}


// OMS接口返回的eip的response  func (m *MeterResponse) EipMetric() (metrics []*EipMetric, err error) 
type EipMetric struct {
	ProviderId string `ason:"ProviderId" json:"provider_id" bson:"provider_id"` // 合作运营商，在阿里云 AppStore 中的统一编号， 阿里云的编号为 26842
	EipId      string `ason:"EipId" json:"eip_id" bson:"eip_id"`                // EIP 实例 ID
	StartTime  string `ason:"StartTime" json:"start_time" bson:"-"`             // 计费数据的发生开始时间，日期格式按照 ISO8601 标准表示，并需要使用 UTC 时间。格式 为:yyyy-MM-ddTHH:mm:ssZ
	EndTime    string `ason:"EndTime" json:"end_time" bson:"-"`                 // 计费数据的发生结束时间，日期格式按照 ISO8601 标准表示，并需要使用 UTC 时间。格式 为:yyyy-MM-ddTHH:mm:ssZ
	EIP        string `ason:"EIP" json:"eip" bson:"eip"`                        // EIP 的 IP 地址
	Bandwidth  string `ason:"Bandwidth" json:"bandwidth" bson:"bandwidth"`      // 购买带宽，单位为 bps
	NetworkIn  string `ason:"NetworkIn" json:"network_in" bson:"network_in"`    // 在 StartTime 和 EndTime 时间段内发生的从外部 网络流入 ECS 的数据量，单位 Bytes
	NetworkOut string `ason:"NetworkOut" json:"network_out" bson:"network_out"` // 在 StartTime 和 EndTime 时间段内发生的从 ECS 流出到外部网络的数据量，单位 Bytes
	UsedTime   string `ason:"UsedTime" json:"used_time" bson:"used_time"`       // 公网 IP 绑定服务器的使用时长， 单位为秒
	Region     string `ason:"Region" json:"region" bson:"region"`               // EIP 所在区域
}

methods:
// 通过oplog确认是否是prepaid资源
func (a *AliEcsEipMeterModel) IsInvalid() (invalid bool, err error) {
	invalid, err = Oplog.IsPrePaidResource(a.Uid, a.RegionId, a.EipId, a.StartTime)
	return
}

func (a *AliEcsEipMeterModel) Save() (err error) 

// 这里主要是设置spec
func (_ *_AliEcsEipMeter) NewModel(uid uint32, m *meter.EipMetric) (eipModel *AliEcsEipMeterModel, err error)

func (_ *_AliEcsEipMeter) Find(m MeterQuery) (aliModels []*AliEcsEipMeterModel, err error) {
	AliEcsEipMeter.Query(func(c *mgo.Collection) {
		err = c.Find(m.Query).All(&aliModels)
	})
	return
}

// 根据query 求出 count,并且返回 MeterResult 这个interface(ToMeterResult() (res *params.Meter) )
func (_ *_AliEcsEipMeter) QueryMeter(m MeterQuery) (eipMeter MeterResult, err error)
	countMeter := &CountMeter{}
	count := 0
	AliEcsEipMeter.Query(func(c *mgo.Collection) {
		trafficSpec := enums.GenerateIpItem(-1, "")
		if m.GetSpec() == trafficSpec {
			pipeline := []bson.M{
				{
					"$match": m.Query,
				},
				{
					"$group": bson.M{
						"_id": nil,
						"count": bson.M{
							"$sum": "$data.network_out",
						},
					},
				},
			}
			err = c.Pipe(pipeline).One(&countMeter)
			return
		}
		count, err = c.Find(m.Query).Count()
		countMeter.Count = int64(count)
	})

	if err == mgo.ErrNotFound {
		return &CountMeter{
			StartTime: m.StartTime,
			Count:     0,
		}, nil
	}

	if err != nil {
		return nil, err
	}
	countMeter.StartTime = m.StartTime

	return countMeter, nil

// group by spec
func (_ *_AliEcsEipMeter) QueryMeterItems(m MeterQuery) (items []*MeterItems, err error) {
	AliEcsEipMeter.Query(func(c *mgo.Collection) {
		pipeline := []bson.M{
			{
				"$match": m.Query,
			},
			{
				"$group": bson.M{
					"_id": "$spec",
				},
			},
		}

		err = c.Pipe(pipeline).All(&items)
	})
	return
}

// group by uid
func (_ *_AliEcsEipMeter) QueryMeterUsers(m MeterQuery) (users []*MeterUser, err error) {
	AliEcsEipMeter.Query(func(c *mgo.Collection) {
		pipeline := []bson.M{
			{
				"$match": m.Query,
			},
			{
				"$group": bson.M{
					"_id": "$uid",
				},
			},
		}

		err = c.Pipe(pipeline).All(&users)
	})
	return
}

// 
func (_ *_AliEcsEipMeter) SupportSpec(spec string) bool {
	return strings.HasPrefix(spec, enums.MeterProduct+":"+enums.MeterGroupEip.String())
}
// -------------------------------------------------------
// -------------------------------------------------------




// -------------------------------------------------------
// -------------------------------------------------------
type CountMeter struct {
	Id        interface{} `bson:"_id"`
	StartTime time.Time   `bson:"-"`
	Count     int64       `bson:"count"`
}

// method:
func (cm *CountMeter) ToMeterResult() (res *params.Meter) {
	res = &params.Meter{
		Time: cm.StartTime.Format(params.MeterTimeFormat),
		Vals: map[string]int64{
			"count": cm.Count,
		},
	}
	return
}
// -------------------------------------------------------
// -------------------------------------------------------


type Meter struct {
	Time   string            `json:"time"`
	Vals   map[string]int64  `json:"values"`
	Extras map[string]string `json:"extras"`
}



// -------------------------------------------------------
// -------------------------------------------------------
type MeterQuery struct {
	StartTime time.Time
	EndTime   time.Time
	Query     bson.M
}

// methods: pipeline
func NewMeterQuery() MeterQuery {
	return MeterQuery{
		Query: make(bson.M),
	}
}

func (m MeterQuery) Start(t time.Time) MeterQuery {
	m.StartTime = t
	if m.Query == nil {
		m.Query = make(bson.M)
	}
	m.Query["start_time"] = bson.M{
		"$gte": t,
	}
	return m
}
// -------------------------------------------------------
// -------------------------------------------------------

interface{}

// eip_meter/disk_meter/instance_meter 都满足IMeter interface
type IMeter interface {
	QueryMeter(m MeterQuery) (MeterResult, error)
	QueryMeterItems(m MeterQuery) ([]*MeterItems, error)
	SupportSpec(spec string) bool
}

type MeterResult interface {
	ToMeterResult() *params.Meter
}





// -------------------------------------------------------
// -------------------------------------------------------
总的 model/meter.go 

var (
	Meter = &_Meter{
		meters: []IMeter{AliEcsInstanceMeter, AliEcsDiskMeter, AliEcsEipMeter},
	}
)

type _Meter struct {
	meters []IMeter
}

method:


// q.Query 增加 spec query

// qvm:ip:2m 有多少条
func (m *_Meter) QueryMeter(spec string, q MeterQuery) (MeterResult, error) 




// 找到 
// 所有的类型的group by spec
// qvm:ip:2m
// qvm:ip:3m
// qvm:ecs:aa
// qvm:ecs:bb
func (m *_Meter) QueryMeterItems(q MeterQuery) ([]*MeterItems, error)



// 找到users
func (m *_Meter) QueryMeterUsers() ([]uint32, error) {
	pageSize := 50
	users := []uint32{}
	for i := 1; ; i++ {
		userModels, err := User.All(i, pageSize)
		if err == mgo.ErrNotFound || len(userModels) == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		for _, v := range userModels {
			users = append(users, v.Uid)
		}
	}

	return users, nil
}

// -------------------------------------------------------
// -------------------------------------------------------

// 先找users
// 再找每个user下面用了多少种类的资源
// 再找每种资源用了多少条



// -------------------------------------------------------
// -------------------------------------------------------
API meter



func MeterRoutes(group *gin.RouterGroup) {
	group.Use(
		middleware.Logger(),
		middleware.bbbbAccountVerify(),
	)
	group.GET("/:option", Meter.GetMeters)
}

query:
GET
/items?start=xxx&end=xxx
/uids?start=xxx&end=xxx


context set key value
start-time
end-time


/items
	- Meter.ListMeterItems(c)
/uids
	- Meter.ListMeterUsers(c)
/default
	- Meter.ListMeters(c)






// -------------------------------------------------------
// -------------------------------------------------------