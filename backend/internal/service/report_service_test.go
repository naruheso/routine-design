package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateYearMonth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		year    int
		month   int
		wantErr bool
	}{
		// 正常系
		{name: "正常系: 2026年4月", year: 2026, month: 4, wantErr: false},
		{name: "正常系: 境界値 2020年1月", year: 2020, month: 1, wantErr: false},
		{name: "正常系: 境界値 2100年12月", year: 2100, month: 12, wantErr: false},
		// 異常系: 年
		{name: "異常系: 年が小さすぎる", year: 2019, month: 4, wantErr: true},
		{name: "異常系: 年が大きすぎる", year: 2101, month: 4, wantErr: true},
		{name: "異常系: 年が0", year: 0, month: 4, wantErr: true},
		{name: "異常系: 年がマイナス", year: -1, month: 4, wantErr: true},
		// 異常系: 月
		{name: "異常系: 月が0", year: 2026, month: 0, wantErr: true},
		{name: "異常系: 月が13", year: 2026, month: 13, wantErr: true},
		{name: "異常系: 月がマイナス", year: 2026, month: -1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateYearMonth(tt.year, tt.month)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrValidation)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
