package gopay

type WeChatNotifyRequest struct {
	ReturnCode         string `xml:"return_code"`
	ReturnMsg          string `xml:"return_msg"`
	ResultCode         string `xml:"result_code"`
	ErrCode            string `xml:"err_code"`
	ErrCodeDes         string `xml:"err_code_des"`
	Appid              string `xml:"appid"`
	MchId              string `xml:"mch_id"`
	DeviceInfo         string `xml:"device_info"`
	NonceStr           string `xml:"nonce_str"`
	Sign               string `xml:"sign"`
	SignType           string `xml:"sign_type"`
	Openid             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           int    `xml:"total_fee"`
	SettlementTotalFee int    `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            int    `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	CouponFee          int    `xml:"coupon_fee"`
	CouponCount        int    `xml:"coupon_count"`
	CouponType0        string `xml:"coupon_type_0"`
	CouponId0          string `xml:"coupon_id_0"`
	CouponFee0         int    `xml:"coupon_fee_0"`
	TransactionId      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	Attach             string `xml:"attach"`
	TimeEnd            string `xml:"time_end"`
}

//企业付款到零钱--请求参数
type WeChatMchTransferReuqest struct {
	MchAppid       string `xml:"mch_appid"`
	Mchid          string `xml:"mchid"`
	DeviceInfo     string `xml:"device_info"`
	NonceStr       string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	PartnerTradeNo string `xml:"partner_trade_no"`
	Openid         string `xml:"openid"`
	CheckName      string `xml:"check_name"`
	ReUserName     string `xml:"re_user_name"`
	Amount         int    `xml:"amount"`
	Desc           string `xml:"desc"`
	SpbillCreateIp string `xml:"spbill_create_ip"`
}
