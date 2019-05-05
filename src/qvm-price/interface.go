type OrderClient interface {
	// 商家
	CreateSeller(l utils.Log, param *params.CreateSellerOpt) (*params.Seller, error)
	GetSeller(l utils.Log, opt *params.GetSellerOpt) (*params.Seller, error)
	UpdateSeller(l utils.Log, param *params.UpdateSellerOpt) (*params.Seller, error)

	// 产品
	CreateProduct(l utils.Log, opt *params.CreateProductOpt) (*params.Product, error)
	UpdateProduct(l utils.Log, opt *params.UpdateProductOpt) (*params.Product, error)
	GetProduct(l utils.Log, opt *params.GetProductOpt) (*params.Product, error)
	ListProductsBySellerID(l utils.Log, opt *params.ListProductsBySellerIDOpt) ([]*params.Product, error)
	ListProductsByIds(l utils.Log, opt *params.ListProductsByIdsOpt) ([]*params.Product, error)
	ListProducts(l utils.Log, opt *params.ListProductsOpt) ([]*params.Product, error)
	ReleaseProduct(l utils.Log, opt *params.ReleaseProductOpt) (product *params.Product, err error)

	// 订单
	CreateOrder(l utils.Log, opt *params.CreateOrderOpt) (order *params.OrderHash, err error)
	GetOrder(l utils.Log, opt *params.GetOrderOpt) (order *params.Order, err error)
	PayOrder(l utils.Log, opt *params.OrderHash) (err error)
	UpdateOrder(l utils.Log, opt *params.UpdateOrderOpt) (order *params.Order, err error)
	ListOrder(l utils.Log, opt *params.ListOrderOpt) (orders []*params.Order, err error)
	OverdrawOrder(l utils.Log, opt *params.OrderHash) (err error)
	RefundOrder(l utils.Log, opt *params.RefundOrderOpt) (err error)

	// 产品订单
	UpdateProductOrder(l utils.Log, opt *params.UpdateProductOrderOpt) (productOrder *params.ProductOrder, err error)
	UpgradeProductOrder(l utils.Log, opt *params.UpgradeProductOrderOpt) (orderHash *params.OrderHash, err error)
	RefundProductOrder(l utils.Log, opt *params.RefundProductOrderOpt) (err error)
	ListProductOrder(l utils.Log, opt *params.ListProductOrderOpt) (productOrders []*params.ProductOrder, err error)
}

// 获取价格的
type PriceClient interface {
	GetItemPrice(l utils.Log, opt *params.GetUserItemPriceOpt) (price *params.UserItemPrice, err error)
	CreateBasePrice(l utils.Log, opt *params.CreateBasePriceOpt) (err error)
	UpdateBasePrice(l utils.Log, opt *params.UpdateBasePriceOpt) (err error)
	SetItem(l utils.Log, opt *params.SetItemOpt) (err error)
}


type BillClient interface {
}


