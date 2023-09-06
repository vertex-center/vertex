package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/types"
)

type ProxyFSRepository struct {
	redirects types.ProxyRedirects
	proxyPath string
}

type ProxyRepositoryParams struct {
	proxyPath string
}

func NewProxyFSRepository(params *ProxyRepositoryParams) ProxyFSRepository {
	if params == nil {
		params = &ProxyRepositoryParams{}
	}
	if params.proxyPath == "" {
		params.proxyPath = path.Join(storage.Path, "proxy")
	}

	err := os.MkdirAll(params.proxyPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		logger.Error(err).
			AddKeyValue("message", "failed to create directory").
			AddKeyValue("path", params.proxyPath).
			Print()
		os.Exit(1)
	}

	repo := ProxyFSRepository{
		redirects: types.ProxyRedirects{},
		proxyPath: params.proxyPath,
	}
	repo.read()

	return repo
}

func (r *ProxyFSRepository) GetRedirects() types.ProxyRedirects {
	return r.redirects
}

func (r *ProxyFSRepository) AddRedirect(id uuid.UUID, redirect types.ProxyRedirect) error {
	r.redirects[id] = redirect
	return r.write()
}

func (r *ProxyFSRepository) RemoveRedirect(id uuid.UUID) error {
	delete(r.redirects, id)
	return r.write()
}

func (r *ProxyFSRepository) read() {
	p := path.Join(r.proxyPath, "redirects.json")
	file, err := os.ReadFile(p)

	if errors.Is(err, os.ErrNotExist) {
		logger.Log("redirects.json doesn't exists or could not be found").Print()
	} else if err != nil {
		logger.Error(err).Print()
		return
	}

	err = json.Unmarshal(file, &r.redirects)
	if err != nil {
		logger.Error(err).Print()
		return
	}
}

func (r *ProxyFSRepository) write() error {
	p := path.Join(r.proxyPath, "redirects.json")

	bytes, err := json.MarshalIndent(r.redirects, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(p, bytes, os.ModePerm)
}
