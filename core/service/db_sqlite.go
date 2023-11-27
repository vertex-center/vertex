package service

import (
	"os"
	"path"

	"github.com/vertex-center/vertex/pkg/storage"
)

func (s *DbService) deleteSqliteDB() error {
	p := path.Join(storage.Path, "database", "vertex.db")
	return os.Remove(p)
}
