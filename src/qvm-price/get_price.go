	getpriceOpts := &params.GetPriceOpt{
		Bandwidth:  &bandwidth,
		CostParams: cost,
		Amount:     oc.amount,
	}


pri, err := price.GetIpPrice(logger, oc.orderOpts.Uid, oc.region, getpriceOpts)



type GetPriceOpt struct {
	InstanceType *string     `json:"instance_type,omitempty"`
	DiskInfo     []*DiskInfo `json:"disk_info,omitempty"`
	Bandwidth    *int        `json:"bandwidth,omitempty"`
	Amount       int         `json:"amount,omitempty"`
	*CostParams
}

type GetUserItemPriceOpt struct {
	Uid  uint32  `json:"uid"`
	Item string  `json:"item"`
	Zone string  `json:"zone"`
	When *string `json:"when,omitempty"`
}


func GetIpPrice(l utils.Log, uid uint32, region string, opt *params.GetPriceOpt) (float64, error) 
	payOpt := payParams.GetUserItemPriceOpt{
		Uid:  uid,
		Zone: getPayZoneFromRegion(region),
		Item: opt.ToIpItem(), // 按照opt设置string
	}
	pri, err := getItemPrice(l, &payOpt)
	if err != nil {
		return 0, err
	}
	return getTotalPrice(*opt, pri), nil
