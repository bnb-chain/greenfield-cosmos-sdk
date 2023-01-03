package types

import (
	"reflect"
	"testing"
)

func TestUpgradeConfig_GetPlan(t *testing.T) {
	type args struct {
		height int64
	}
	tests := []struct {
		name string
		c    UpgradeConfig
		args args
		want []*Plan
	}{
		{
			name: "TestUpgradeConfig_GetPlan Case 1",
			c: UpgradeConfig(map[int64][]*Plan{
				1: {{
					Name:   "Upgrade-1",
					Height: 1,
				}},
				11: {{
					Name:   "Upgrade-2",
					Height: 11,
				}},
				20: {{
					Name:   "Upgrade-3",
					Height: 20,
				}},
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
			c: UpgradeConfig(map[int64][]*Plan{
				1: {{
					Name:   "Upgrade-1",
					Height: 1,
				}},
				11: {{
					Name:   "Upgrade-2",
					Height: 11,
				}},
				20: {{
					Name:   "Upgrade-3",
					Height: 20,
				}},
			}),
			args: args{
				height: 20,
			},
			want: []*Plan{{
				Name:   "Upgrade-3",
				Height: 20,
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
