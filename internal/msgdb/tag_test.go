package msgdb

import (
	"testing"

	"github.com/quick-im/quick-im-core/internal/msgdb/model"
)

func TestParseStruct(t *testing.T) {
	d := model.Msg{
		MsgId: "",
	}
	name, info, err := parseStruct(d)
	if err != nil {
		t.Error(err)
	}
	t.Log(name, info[0].Field, info[0].Tags)
}
