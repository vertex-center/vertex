package service

import (
	"reflect"

	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
	"gorm.io/gorm"
)

func (s *DbService) copyDb(from *gorm.DB, to *gorm.DB) error {
	dbCopy := types.NewEventDbCopy()
	dbCopy.AddTable(types.AdminSettings{})
	s.ctx.DispatchEvent(dbCopy)

	toTx := to.Begin()

	for _, t := range dbCopy.All() {
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
		vlog.String("table", reflect.TypeOf(tp).String()),
		vlog.String("from", from.Name()),
		vlog.String("to", to.Name()),
	)

	log.Debug("getting items from database", vlog.String("table", reflect.TypeOf(tp).String()))
	res := from.Find(items)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Debug("no items found", vlog.String("table", reflect.TypeOf(tp).String()))
		return nil
	}

	log.Debug("found items", vlog.String("table", reflect.TypeOf(tp).String()), vlog.Int64("count", res.RowsAffected))
	return to.Create(items).Error
}
