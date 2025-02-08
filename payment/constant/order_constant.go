package constant

type OrderStatus string

const (
	OrderStatusOpen           OrderStatus = "open"
	OrderStatusRequestPayment OrderStatus = "request_payment"
	OrderStatusConfirm        OrderStatus = "confirm"
	OrderStatusReject         OrderStatus = "reject"
)

func IsOrderRequestPayment(status OrderStatus) bool {
	return OrderStatusRequestPayment == status
}
