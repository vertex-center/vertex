package service

import (
	"errors"
	"reflect"

	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gorm.io/gorm"
)

func (s *DbService) copyDb(from *gorm.DB, to *gorm.DB) error {
	// Adding a table here must also be added in the Connect() method in adapter/db_config_fs.go
	tables := []interface{}{
		types.AdminSettings{},
	}

	toTx := to.Begin()

	for _, t := range tables {
		err := s.copyDbTable(t, from, toTx)
		if err != nil {
			toTx.Rollback()
			return err
		}
	}

	return toTx.Commit().Error
}

func (s *DbService) copyDbTable(tp interface{}, from *gorm.DB, to *gorm.DB) error {
	t := reflect.TypeOf(tp)
	t = reflect.SliceOf(t)
	items := reflect.New(t).Interface()

	log.Info("copying table",
		vlog.String("table", reflect.TypeOf(items).String()),
		vlog.String("from", from.Name()),
		vlog.String("to", to.Name()),
	)

	err := from.Find(items).Error
	if errors.Is(err, gorm.ErrEmptySlice) {
		return nil
	} else if err != nil {
		return err
	}
	return to.Create(items).Error
}
