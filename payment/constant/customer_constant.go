package constant

type CustomerStatus string

const (
	CustomerStatusActive   CustomerStatus = "active"
	CustomerStatusInActive CustomerStatus = "inactive"
)

func IsActiveCustomer(status CustomerStatus) bool {
	return CustomerStatusActive == status
}
