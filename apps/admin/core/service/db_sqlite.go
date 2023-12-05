package service

import (
	"os"
	"path"

	"github.com/vertex-center/vertex/core/types/storage"
)

func (s *DbService) deleteSqliteDB() error {
	p := path.Join(storage.FSPath, "database", "vertex.db")
	return os.Remove(p)
}
