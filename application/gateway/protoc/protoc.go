package protoc

import (
	"context"
	"net/http"
	"sync"

	"github.com/quick-im/quick-im-core/internal/quickerr"
)

type ProtocHandler interface {
	Handler(ctx context.Context) http.HandlerFunc
}

var (
	lock        sync.RWMutex
	protocDrive = make(map[string]ProtocHandler)
)

func RegisterDrive(name string, drive ProtocHandler) {
	lock.Lock()
	defer lock.Unlock()
	protocDrive[name] = drive
}

func Handler(name string) (ProtocHandler, error) {
	lock.RLock()
	defer lock.RUnlock()
	if drive, ok := protocDrive[name]; ok {
		return drive, nil
	}
	return nil, quickerr.ErrDriveNotSupport
}
