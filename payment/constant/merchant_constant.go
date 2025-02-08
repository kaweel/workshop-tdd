package constant

type MerchantStatus string

const (
	MerchantStatusActive   MerchantStatus = "active"
	MerchantStatusSuspend  MerchantStatus = "suspend"
	MerchantStatusInActive MerchantStatus = "inactive"
)

func IsActiveMerchant(status MerchantStatus) bool {
	return MerchantStatusActive == status
}
