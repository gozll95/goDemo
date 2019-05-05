package main

import (
	"fmt"
)

type Phone struct {
	PayMap map[string]Pay
}

func (p *Phone) OpenWeChatPay() {
	weChatPay := &WeChatPay{}
	p.PayMap["wechat_pay"] = weChatPay
}

func (p *Phone) OpenAliPay() {
	aliPay := &AliPay{}
	p.PayMap["ali_pay"] = aliPay
}

func (p *Phone) OpenPay(name string, pay Pay) {
	p.PayMap[name] = pay
}

func (p *Phone) PayMoney(name string, money float32) (err error) {
	pay, ok := p.PayMap[name]
	if !ok {
		err = fmt.Errorf("不支持[%s]支付方式", name)
		return
	}

	err = pay.pay(1023, money)
	return
}



////////////////////////////////////////////////
	 var tmp interface{} = weChat
	 _, ok := tmp.(Pay)
	 if ok {
		 fmt.Println("weChat is implement Pay interface")
	 	//phone.OpenPay("wechat_pay", weChat)
	 }
////////////////////////////////////////////////

随机:

type RandBalance struct {

}

func (r *RandBalance) DoBalance(addrList []string) string {
	 l := len(addrList)
	 index := rand.Intn(l)
	 return addrList[index]
}


////////////////////////////////////////////////
轮询
type RoundBalance struct {
	curIndex int
}

func (r *RoundBalance) DoBalance(addrList []string) string {
	l := len(addrList)
	r.curIndex = r.curIndex % l
	addr := addrList[r.curIndex]
	r.curIndex++
	return addr
}
////////////////////////////////////////////////