package constant

type PaymentChannel string

const (
	PaymentChannelDebit     PaymentChannel = "debit"
	PaymentChannelCredit    PaymentChannel = "credit"
	PaymentChannelPromptPay PaymentChannel = "promptpay"
	PaymentChannelQRPayment PaymentChannel = "qrpayment"
)

func IsValidPaymentChannel(channel PaymentChannel) bool {
	switch channel {
	case PaymentChannelDebit, PaymentChannelCredit, PaymentChannelPromptPay, PaymentChannelQRPayment:
		return true
	default:
		return false
	}
}

type PaymentTranasctionStatus string

const (
	PaymentTranasctionStatusConfirm PaymentTranasctionStatus = "comfirm"
	PaymentTranasctionStatusReject  PaymentTranasctionStatus = "reject"
)

var KafkaTopicPaymentTransaction = "payment-transaction"
