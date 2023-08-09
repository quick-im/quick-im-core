package msgdb

import (
	"reflect"
	"testing"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

func TestModel(t *testing.T) {
	type args struct {
		model any
	}
	tests := []struct {
		name string
		args args
		want tableInfo
	}{
		{
			name: "test1",
			args: args{
				model: model.Msg{},
			},
			want: tableInfo{
				table: "msg",
				pk:    "msg_id",
				index: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Model(tt.args.model); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Model() = %v, want %v", got, tt.want)
			}
		})
	}
}
