package accrualsdk

type OrderAccrualInfoStatus string

const (
	OrderAccrualInfoStatusRegistered OrderAccrualInfoStatus = "REGISTERED"
	OrderAccrualInfoStatusInvalid    OrderAccrualInfoStatus = "INVALID"
	OrderAccrualInfoStatusProcessing OrderAccrualInfoStatus = "PROCESSING"
	OrderAccrualInfoStatusProcessed  OrderAccrualInfoStatus = "PROCESSED"
)

type OrderAccrualInfo struct {
	OrderNumber string                 `json:"order"`
	Status      OrderAccrualInfoStatus `json:"status"`
	Accrual     *float64               `json:"accrual,omitempty"`
}
