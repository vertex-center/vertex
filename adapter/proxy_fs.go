package adapter

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type ProxyFSRepository struct {
	redirects types.ProxyRedirects
	proxyPath string
}

type ProxyFSAdapterParams struct {
	proxyPath string
}

func NewProxyFSAdapter(params *ProxyFSAdapterParams) types.ProxyAdapterPort {
	if params == nil {
		params = &ProxyFSAdapterParams{}
	}
	if params.proxyPath == "" {
		params.proxyPath = path.Join(storage.Path, "proxy")
	}

	err := os.MkdirAll(params.proxyPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Default.Error(err,
			vlog.String("message", "failed to create directory"),
			vlog.String("path", params.proxyPath),
		)
		os.Exit(1)
	}

	adapter := &ProxyFSRepository{
		redirects: types.ProxyRedirects{},
		proxyPath: params.proxyPath,
	}
	adapter.read()

	return adapter
}

func (a *ProxyFSRepository) GetRedirects() types.ProxyRedirects {
	return a.redirects
}

func (a *ProxyFSRepository) AddRedirect(id uuid.UUID, redirect types.ProxyRedirect) error {
	a.redirects[id] = redirect
	return a.write()
}

func (a *ProxyFSRepository) RemoveRedirect(id uuid.UUID) error {
	delete(a.redirects, id)
	return a.write()
}

func (a *ProxyFSRepository) read() {
	p := path.Join(a.proxyPath, "redirects.json")
	file, err := os.ReadFile(p)

	if errors.Is(err, os.ErrNotExist) {
		log.Default.Info("redirects.json doesn't exists or could not be found")
	} else if err != nil {
		log.Default.Error(err)
		return
	}

	err = json.Unmarshal(file, &a.redirects)
	if err != nil {
		log.Default.Error(err)
		return
	}
}

func (a *ProxyFSRepository) write() error {
	p := path.Join(a.proxyPath, "redirects.json")

	bytes, err := json.MarshalIndent(a.redirects, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(p, bytes, os.ModePerm)
}
