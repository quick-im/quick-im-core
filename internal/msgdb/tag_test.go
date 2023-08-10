package msgdb

import (
	"reflect"
	"testing"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

func TestParseStruct(t *testing.T) {
	type args struct {
		data any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []structInfo
		wantErr bool
	}{
		{
			name: "tast1",
			args: args{
				data: model.Msg{
					MsgId: "",
				},
			},
			want: "msg",
			want1: []structInfo{
				{
					Field: "msg_id",
					Tags:  []string{"pk"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := parseStruct(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseStruct() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("parseStruct() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
