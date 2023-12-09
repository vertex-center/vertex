package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"
	"github.com/vertex-center/vertex/cmd/storage"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

var (
	errReverseProxyNotFound       = errors.New("redirects.json doesn't exists or could not be found")
	errReverseProxyFailedToRead   = errors.New("failed to read redirects.json")
	errReverseProxyFailedToDecode = errors.New("failed to decode redirects.json")
)

type proxyFSAdapter struct {
	redirects      types.ProxyRedirects
	redirectsMutex sync.RWMutex

	proxyPath string
}

type ProxyFSAdapterParams struct {
	proxyPath string
}

func NewProxyFSAdapter(params *ProxyFSAdapterParams) port.ProxyAdapter {
	if params == nil {
		params = &ProxyFSAdapterParams{}
	}
	if params.proxyPath == "" {
		params.proxyPath = path.Join(storage.FSPath, "proxy")
	}

	err := os.MkdirAll(params.proxyPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.proxyPath),
		)
		os.Exit(1)
	}

	adapter := &proxyFSAdapter{
		redirects:      types.ProxyRedirects{},
		redirectsMutex: sync.RWMutex{},

		proxyPath: params.proxyPath,
	}

	err = adapter.read()
	if errors.Is(err, errReverseProxyFailedToDecode) {
		log.Error(err)
	}

	return adapter
}

func (a *proxyFSAdapter) GetRedirects() types.ProxyRedirects {
	a.redirectsMutex.RLock()
	defer a.redirectsMutex.RUnlock()

	return a.redirects
}

func (a *proxyFSAdapter) GetRedirectByHost(host string) *types.ProxyRedirect {
	a.redirectsMutex.RLock()
	defer a.redirectsMutex.RUnlock()

	for _, redirect := range a.redirects {
		if redirect.Source == host {
			return &redirect
		}
	}
	return nil
}

func (a *proxyFSAdapter) AddRedirect(id uuid.UUID, redirect types.ProxyRedirect) error {
	func() {
		a.redirectsMutex.Lock()
		defer a.redirectsMutex.Unlock()
		a.redirects[id] = redirect
	}()
	return a.write()
}

func (a *proxyFSAdapter) RemoveRedirect(id uuid.UUID) error {
	func() {
		a.redirectsMutex.Lock()
		defer a.redirectsMutex.Unlock()
		delete(a.redirects, id)
	}()
	return a.write()
}

func (a *proxyFSAdapter) read() error {
	p := path.Join(a.proxyPath, "redirects.json")
	file, err := os.ReadFile(p)

	if errors.Is(err, os.ErrNotExist) {
		return errReverseProxyNotFound
	} else if err != nil {
		return fmt.Errorf("%w: %w", errReverseProxyFailedToRead, err)
	}

	a.redirectsMutex.Lock()
	defer a.redirectsMutex.Unlock()

	err = json.Unmarshal(file, &a.redirects)
	if err != nil {
		return fmt.Errorf("%w: %w", errReverseProxyFailedToDecode, err)
	}

	return nil
}

func (a *proxyFSAdapter) write() error {
	p := path.Join(a.proxyPath, "redirects.json")

	a.redirectsMutex.RLock()
	defer a.redirectsMutex.RUnlock()

	bytes, err := json.MarshalIndent(a.redirects, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(p, bytes, os.ModePerm)
}
