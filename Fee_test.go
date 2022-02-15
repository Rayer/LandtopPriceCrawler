package TelecamFeeAnalyzer

import "testing"

func TestFeeEntry_RealMonthlyCost(t *testing.T) {
	tests := []struct {
		name   string
		fields FeeEntry
		want   int
	}{
		{
			name: "Case1",
			fields: FeeEntry{
				ContractMonth: 24,
				MonthlyFee:    499,
				Rebate:        3000,
			},
			want: 374,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FeeEntry{
				Provider:      tt.fields.Provider,
				Description:   tt.fields.Description,
				ContractMonth: tt.fields.ContractMonth,
				PrePaid:       tt.fields.PrePaid,
				MonthlyFee:    tt.fields.MonthlyFee,
				Rebate:        tt.fields.Rebate,
			}
			if got := f.RealMonthlyCost(); got != tt.want {
				t.Errorf("RealMonthlyCost() = %v, want %v", got, tt.want)
			}
		})
	}
}
