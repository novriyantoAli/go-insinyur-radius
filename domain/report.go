package domain

import "context"

// ReportFinance ...
type ReportFinance struct {
	Name  *string `json:"name"`
	Value *int64  `json:"value"`
}

// ReportRepository ...
type ReportRepository interface {
	ReportFinanceCurrentDay(ctx context.Context)(res []ReportFinance, err error)
	ReportFinanceCurrentMonth(ctx context.Context)(res []ReportFinance, err error)
	ReportFinanceCurrentYear(ctx context.Context)(res []ReportFinance, err error)
}

// ReportUsecase ...
type ReportUsecase interface {
	ReportFinanceCurrent(ctx context.Context, menu int)(res []ReportFinance, err error)
	ReportExpirationCurrent(c context.Context)(res []Users, err error)
}
