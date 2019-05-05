
后付费-----> 直接 create ---->oplog 
预付费-----> create order --->将条目arg填进去---> 创建订单--->更新订单价格----->提交订单
续费 -----> 调bss


商业运营 callback 
 
----解析 订单里 memo信息(mem里存入的是创建的args的string化) ---> 直接create



// ---------------------------------
// 预备创建订单
type OrderCreator struct {
	region    string
	trade     float64  // 这里是qvm计算的价格
	period    int
	amount    int
	orderOpts payParams.CreateOrderOpt
}

// methods:
-  NewOrderCreator(uid uint32, region string, period, amount int) *OrderCreator--->  // eg: orderCreator := NewOrderCreator(uid, args.RegionId, buyMonth, 1) // 这里 amount是什么)

// orderType enums: create | renew | to_prepaid  
// logger + 什么类型的订单 + 创建的是什么资源eip/disk + 创建参数 -->memo --> encode --> set order memo
- AddMemo(logger utils.Log, orderType params.OrderType, resourceType enums.ResourceType, data interface{}) 
		- dataStr = json.Marshal(data) // 序列化data
		- orderMemo:=&params.OrderMemo{ 	// 创建OrderMemo(里面存创建资源的信息)
			OrderType:    orderType,
			ResourceType: resourceType,
			Data:         string(dataStr),
		}
		- memo=orderMemo.Encode()
		- oc.orderOpts.Memo = memo // 放到预创建订单的memo里


// 用于 设置 价格  从商业运营那边 获取 产品ID
- AddEip(logger utils.Log, bandwidth int, cost *params.CostParams) error
	- getpriceOpts := &params.GetPriceOpt{ // 形成GetPriceOpt
		Bandwidth:  &bandwidth,
		CostParams: cost,
		Amount:     oc.amount,
	}
	- pri, err := price.GetIpPrice(logger, oc.orderOpts.Uid, oc.region, getpriceOpts) // 获取价格
	- oc.trade += pri // 设置价格
		getProductOpts := &payParams.GetProductOpt{ // 用于获取产品信息
			Model: common.String(enums.GenerateIpItem(bandwidth, "")),
		}
	- prod, err := client.GetOrderClient().GetProduct(logger, getProductOpts) // 获取产品
	- oc.addProductID(prod.ID) // 添加 产品ID

//
func (oc *OrderCreator) Submit(logger utils.Log) (*payParams.Order, error) {
	orderHash, err := client.GetOrderClient().CreateOrder(logger, &oc.orderOpts) // 创建订单
	if err != nil {
		logger.Errorf("client.GetOrderClient().CreateOrder(%v, %v):%v", logger, oc.orderOpts, err)
		return nil, err
	}
	order, err := client.GetOrderClient().UpdateOrder(logger, &payParams.UpdateOrderOpt{ // 设置价格
		OrderHash:   &orderHash.OrderHash,
		ActuallyFee: &oc.trade,
	})
	return order, err
}

// ------------------------------------------------------------------
prepaid eip
case enums.PrePaid:
	// create order
	order, err := EipPayHelper.CreateEipOrder(request.Logger(), request.User().Uid, &allocateEipAddressArgs)
			buyMonth, err := extractMonthFromCostParams(args.CostParams)
			if err != nil {
				logger.Errorf("extractMonthFromCostParams(%v): %v", args, err)
				return nil, errors.ErrorParameters
			}
			// 创建【预备创建订单参数】
			orderCreator := NewOrderCreator(uid, args.RegionId, buyMonth, 1)
			// add memo - 将此次创建参数以及其他信息添加进memo中
			if err = orderCreator.AddMemo(logger, params.OrderTypeCreate, enums.ResourceTypeIp, args); err != nil {
				return nil, err
			}

			if args.Bandwidth == nil {
				logger.Errorf("args.Bandwidth can not be nil: %v", args)
				return nil, errors.InvalidArgument
			}

			bandwidth, err := strconv.Atoi(*args.Bandwidth)
			if err != nil {
				logger.Errorf("strconv.Atoi(%v) %v", *args.Bandwidth, err)
				return nil, errors.InvalidArgument
			}
			// 添加价格-产品ID
			if err = orderCreator.AddEip(logger, bandwidth, args.CostParams); err != nil {
				return nil, err
			}
			// 创建订单 - 设置价格
			return orderCreator.Submit(logger)
	if err != nil {
		request.Logger().Errorf("EipPayHelper.CreateEipOrder(logger,%d,%v) %v", request.User().Uid, allocateEipAddressArgs, err)
		errorResponse(c, err)
		return
	}

	response(c, order)

// ------------------------------------------------------------------
callback

middle Callback 
set context key value 

order key
memo key
region key
和 其他 资源所需要的key 比如 ecs /vpc / ... 

order := c.MustGet(middleware.OrderInfoKey).(*payParams.Order)
orderMemo := c.MustGet(middleware.OrderMemoKey).(*params.OrderMemo)

	switch orderMemo.ResourceType {
	case enums.ResourceTypeInstance:
		//...
	case enums.ResourceTypeIp:
		// eip order callback
		// 通过orderMemo里存储的orderMemo.OrderType来执行后续操作
		err := EipPayHelper.EipPayCallbackHandle(c, order.BuyerID, orderMemo)
		if err != nil {
			logger.Errorf("EipPayHelper.EipPayCallbackHandle(%v,%d,%v): %v", c, order.BuyerID, orderMemo, err)
			errorResponse(c, err)
			return
		}
	default:
		logger.Errorf("orderMemo.ResourceType(%s) is not supported now!", orderMemo.ResourceType)
		errorResponse(c, errors.ErrorParameters)
		return
	}

	// 更新 产品
	for _, prod := range order.ProductOrders {
		opts := &payParams.UpdateProductOrderOpt{
			ID: prod.ID,
		}
		_, err := client.GetOrderClient().UpdateProductOrder(logger, opts)
		if err != nil {
			logger.Errorf("client.GetOrderClient().UpdateProductOrder(%v):%v", *opts, err)
			errorResponse(c, err)
			return
		}
	}
	response(c, nil)


// ------------------------------------------------------------------
renew

var args bss.RenewInstance
allocationId := c.Param("id")

args.ProductCode = bss.EipProductCode
args.RegionId = request.region
args.InstanceId = allocationId
args.InitCost()

// find eip in model
// create order
order, err := EipPayHelper.CreateRenewEipOrder(request.Logger(), request.User().Uid, &args, eip)
		- func (_ *_EipPayHelper) CreateRenewEipOrder(logger *utils.Logger, uid uint32, args *bss.RenewInstance, eip *model.EipModel) (order *payParams.Order, err error)
				- orderCreator := NewOrderCreator(uid, args.RegionId, args.RenewPeriod, 1) // 创建【预创建订单参数】
				- err = orderCreator.AddMemo(logger, params.OrderTypeRenew, enums.ResourceTypeIp, args) // 添加memo
				- err = orderCreator.AddEip(logger, bandwidth, args.CostParams)
				- return orderCreator.Submit(logger)
// ------------------------------------------------------------------
call back ---> resource = eip的时候

EipPayCallbackHandle
	switch orderMemo.OrderType
		case params.OrderTypeCreate:
		case params.OrderTypeRenew:
			args := &bss.RenewInstance{}
			err := orderMemo.GetData(args) //从memo中将bss renew参数解析出来
			err = DoRenewEip(args, request, uid)
		default:

// ------------------------------------------------------------------







// 计价参数
type CostParams struct {
	CostChargeType enums.ChargeType `ason:"-" json:"cost_charge_type" bson:"cost_charge_type"` // 付费方式 PostPaid | PrePaid
	CostChargeMode enums.ChargeMode `ason:"-" json:"cost_charge_mode" bson:"cost_charge_mode"` // 按什么计费 
	CostPeriodUnit enums.PeriodUnit `ason:"-" json:"cost_period_unit" bson:"cost_period_unit"` // 计费单元 Hour | Day | Week | Month | Year 
	CostPeriod     int              `ason:"-" json:"cost_period" bson:"cost_period"`           // 计费周期
}

type CostParamer interface {
	GetChargeType() enums.ChargeType
	GetChargeMode() enums.ChargeMode // enums: traffic | bandwith ...
	GetPeriodUnit() enums.PeriodUnit // enums: hour | week | month | year
	GetPeriod() int
	IsCostValid() bool
}

// 订单

type OrderStatus int

const (
	OrderStatusUnpaid OrderStatus = iota + 1 // 未支付
	OrderStatusPaid 	// 支付
	OrderStatusExpired // 过期
	OrderStatusPaidByOthers // 被其他支付
)

// 一个订单下有很多子订单
type Order struct {
	ID            int64          `json:"id"`
	OrderHash     string         `json:"order"`
	SellerID      int64          `json:"seller_id"`
	BuyerID       uint32         `json:"buyer_id"`
	Fee           float64        `json:"fee"`
	ActuallyFee   float64        `json:"actually_fee"`
	Memo          string         `json:"memo"`
	UpdateTime    time.Time      `json:"update_time"`
	CreateTime    time.Time      `json:"create_time"`
	PayTime       time.Time      `json:"pay_time"`
	Status        int            `json:"status"`
	Products      []Product      `json:"products"`
	ProductOrders []ProductOrder `json:"product_orders"`
}

type Product struct {
	ID          int64         `json:"id"`
	SellerID    int64         `json:"seller_id"`
	Name        string        `json:"name"`
	Model       string        `json:"model"`
	Spu         string        `json:"spu"`
	Unit        ProductUnit   `json:"unit"`
	Price       float64       `json:"price"`
	ExpiresIn   int           `json:"expires_in"`
	Property    string        `json:"property"`
	Description string        `json:"description"`
	UpdateTime  time.Time     `json:"update_time"`
	CreateTime  time.Time     `json:"create_time"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Status      ProductStatus `json:"status"`
	Version     int           `json:"version"`
}

type ProductOrder struct {
	ID              int64              `json:"id"`
	ProductID       int64              `json:"product_id"`
	SellerID        int64              `json:"seller_id"`
	BuyerID         uint32             `json:"buyer_id"`
	OrderID         int64              `json:"order_id"`
	OrderHash       string             `json:"order_hash"`
	OrderType       OrderType          `json:"order_type"`
	ProductOrderID  int64              `json:"product_order_id"`
	ProductName     string             `json:"product_name"`
	ProductProperty string             `json:"product_property"`
	Property        string             `json:"property"`
	Duration        int                `json:"duration"`
	TimeDuration    time.Duration      `json:"time_duration"`
	Quantity        int                `json:"quantity"`
	ItemFee         float64            `json:"item_fee"`
	Fee             float64            `json:"fee"`
	UpdateTime      time.Time          `json:"update_time"`
	CreateTime      time.Time          `json:"create_time"`
	StartTime       time.Time          `json:"start_time"`
	EndTime         time.Time          `json:"end_time"`
	Status          ProductOrderStatus `json:"status"`
}

order 对象:

一个订单可以包涵同一商家， 同一订单类型的多个商品订单

{
	"id":17, //订单id
	"order_hash":"47ed733b8d10be225eceba344d533586", //订单号
	"seller_id":19, //商家id
	"buyer_id":3774353, //买家uid
	"fee":339, //订单费用
	"actually_fee":339, //用户需要支付的费用
	"memo":"computer update", //订单说明
	"update_time":"2015-10-21T16:43:25+08:00", //更新时间
	"create_time":"2015-10-21T16:43:25+08:00", //创建时间
	"pay_time":"0001-01-01T08:00:00+08:00", //支付时间
    "expired_time": "2015-10-22T16:43:25+08:00", // 订单过期时间
	"status":1, //订单状态 1: 未支付 2: 已支付 3:作废, 4: 垫付
    “products":[
        product1,
		
        product2,
        ...
    ], // 可能为 null
	"product_orders":
	[
		product_order1, // 单个商品订单
		product_order2,
		...
	] // 可能为 null
}

product 对象:

{
	"id":19, //商品id
	"seller_id":18, //商家id 外键 **seller**对象的id
	"name":"compute_3", //商品名称
	"model":"evm:compute:c:3", //产品型号，sku, 支持根据产品型号获取产品信息， 推荐构成方式:seller_name:product_property
    "spu": "evm:compute", 
	"unit":1, //单位, 默认是按月 1: 年，2: 月 3: 周 4: 天 99: 一次性购买
	"price":0, //价格， 支持价格为0， 单位：元，精度5位，长度15位
    "expires_in": 3600, // 订单未支付过期时间，0表示不过期， 单位 second
	"property":"{\"cpu\":\"16 core\", \"memory\":\"64\"}", //产品属性
	"description":"4核1G内存", // 产品描述
	"update_time":"2015-10-20T11:25:55+08:00", // 更新时间
	"create_time":"2015-10-20T11:24:33+08:00", // 创建时间
	"start_time":"2015-12-20T00:00:00+08:00", // 产品上线时间
	"end_time":"2018-12-19T23:59:59+08:00", // 产品下线时间
	"status":1 // 产品状态  1:新建 2:在线 3:已失效 4:已删除
}

roduct_order 对象:

{
	"id":14, // id
	"product_id":21, // 产品id
	"seller_id":19, // 商家id
	"buyer_id":3774353, // 买家uid
	"order_id":17, // 总订单id
	"order_hash":"47ed733b8d10be225eceba344d533586", // 订单号
    "order_type":1, //订单类型  1: 新建 2:续费 3: 升级 4: 补偿 5:退款
    "product_order_id": 0, // 关联订单id，默认为0
    "product_name": "入门型点播云",
    "product_property": "{\"space\":20, \"tranfer\":60}",
	"property":"", // 订单属性
	"duration":2, // unit对应的倍数，如服务器按月购买，这边可以一次购买2个月
    "time_duration":0, // 自定义时间长度, 单位： second
	"quantity":2, // 数量，同样配置
	"item_fee":339, // 费用
    "fee": 339, // 实际支付费用
	"update_time":"2015-10-21T16:43:25+08:00", //更新时间
	"create_time":"2015-10-21T16:43:25+08:00", //创建时间
	"start_time":"0001-01-01T08:00:00+08:00", //服务开始时间
	"end_time":"0001-01-01T08:00:00+08:00", //服务结束时间
	"status": 1 // 订单状态： 1: 新建 2:完成
}









type CreateOrderOpt struct {
	Uid    uint32       `json:"uid"`
	Memo   *string      `json:"memo,omitempty"`
	Orders []SmartOrder `json:"orders,omitempty"`
}

type SmartOrder struct {
	ProductID    int64   `json:"product_id"`  //产品ID
	Duration     *int    `json:"duration,omitempty"`
	TimeDuration *int    `json:"time_duration,omitempty"` // 自定义事件长度
	Quantity     int     `json:"quantity"`  // 数量
	Property     *string `json:"property,omitempty"` // 订单属性
}




const (
	OrderTypeCreate    OrderType = "create"
	OrderTypeRenew     OrderType = "renew"
	OrderTypeToPrePaid OrderType = "to_prepaid"
)

//
type OrderMemo struct {
	OrderType    OrderType          `json:"order_type"`
	ResourceType enums.ResourceType `json:"resource_type"`
	Data         string             `json:"data"`
}


type GetPriceOpt struct {
	InstanceType *string     `json:"instance_type,omitempty"`
	DiskInfo     []*DiskInfo `json:"disk_info,omitempty"`
	Bandwidth    *int        `json:"bandwidth,omitempty"`
	Amount       int         `json:"amount,omitempty"`
	*CostParams
}


type GetProductOpt struct {
	ID    *int64  `json:"id,omitempty"`
	Model *string `json:"model,omitempty"`
}

//----------------------
商业运营
produce 产品
seller 卖家
order 订单
//----------------------