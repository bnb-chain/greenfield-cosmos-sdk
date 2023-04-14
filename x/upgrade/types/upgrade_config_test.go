package types

import (
	"reflect"
	"testing"
)

func TestUpgradeConfig(t *testing.T) {
	type args struct {
		height int64
	}
	tests := []struct {
		name string
		c    *UpgradeConfig
		args args
		want []*Plan
	}{
		{
			name: "TestUpgradeConfig_GetPlan Case 1",
			c: NewUpgradeConfig().SetPlan(&Plan{
				Name:   "Upgrade-1",
				Height: 1,
			}).SetPlan(&Plan{
				Name:   "Upgrade-2",
				Height: 11,
			}).SetPlan(&Plan{
				Name:   "Upgrade-3",
				Height: 20,
			}),
			args: args{
				height: 10,
			},
			want: []*Plan{{
				Name:   "Upgrade-2",
				Height: 11,
			}},
		},
		{
			name: "TestUpgradeConfig_GetPlan Case 2",
			c: NewUpgradeConfig().SetPlan(&Plan{
				Name:   "Upgrade-1",
				Height: 1,
			}).SetPlan(&Plan{
				Name:   "Upgrade-2",
				Height: 11,
			}).SetPlan(&Plan{
				Name:   "Upgrade-3",
				Height: 20,
			}),
			args: args{
				height: 20,
			},
			want: []*Plan{{
				Name:   "Upgrade-3",
				Height: 20,
			}},
		},
		{
			name: "TestUpgradeConfig_SetPlan Override 1",
			c: NewUpgradeConfig().SetPlan(&Plan{
				Name:   "Upgrade-1",
				Height: 1,
			}).SetPlan(&Plan{
				Name:   "Upgrade-2",
				Height: 11,
			}).SetPlan(&Plan{
				Name:   "Upgrade-3",
				Height: 20,
			}).SetPlan(&Plan{
				Name:   "Upgrade-2",
				Height: 19,
			}),
			args: args{
				height: 11,
			},
			want: []*Plan{{
				Name:   "Upgrade-2",
				Height: 19,
			}},
		},
		{
			name: "TestUpgradeConfig_SetPlan Override 2",
			c: NewUpgradeConfig().SetPlan(&Plan{
				Name:   "Upgrade-1",
				Height: 1,
			}).SetPlan(&Plan{
				Name:   "Upgrade-2",
				Height: 11,
			}).SetPlan(&Plan{
				Name:   "Upgrade-3",
				Height: 20,
			}).SetPlan(&Plan{
				Name:   "Upgrade-2",
				Height: 11,
				Info:   "override",
			}),
			args: args{
				height: 10,
			},
			want: []*Plan{{
				Name:   "Upgrade-2",
				Height: 11,
				Info:   "override",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetPlan(tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpgradeConfig.GetPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}
