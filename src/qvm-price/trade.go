type TradeClient struct {
	Host string
	*rpc.Client
}


func NewTradeClient(host, authUser, authPassword, authHost string) *TradeClient {
	cli := &http.Client{
		Transport: NewOauthTransport(authUser, authPassword, authHost, nil),
	}
	return &TradeClient{
		Host:   strings.TrimSuffix(host, "/"),
		Client: &rpc.Client{Client: cli},
	}
}

func (tc *TradeClient) CreateSeller(l utils.Log, opt *params.CreateSellerOpt) (seller *params.Seller, err error) {
	val, err := convert.ToWWWFormUrlEncoded(opt, "json")
	if err != nil {
		return
	}
	err = tc.CallWithForm(logger.NewLogger(l), &seller, http.MethodPost, tc.Host+"/seller/new", val)
	return
}