package TelecamFeeAnalyzer

import "fmt"

type FeeEntry struct {
	Provider      string
	Plan          string
	SignalType    string //4G or 5G
	Limit         int    //0 for unlimited
	HotspotLimit  int
	ContractMonth int
	PrePaid       int
	MonthlyFee    int
	Rebate        int
	Description   string
}

func (f *FeeEntry) RealMonthlyCost() int {
	return f.MonthlyFee - (f.Rebate / f.ContractMonth)
}

func GetLimitString(limit int) string {
	if limit == 0 {
		return "吃到飽"
	}
	return fmt.Sprintf("%d GB", limit)
}

func (f FeeEntry) String() string {
	return fmt.Sprintf("%v[%v](%v)\t合約期 %v\t每月單價 %v\t實際月繳 %v\t用量限制 %v\t%v", f.Provider, f.Plan, f.SignalType, f.ContractMonth, f.RealMonthlyCost(), f.MonthlyFee, GetLimitString(f.Limit), f.Description)
}
