package models

import (
	"time"

	"github.com/AVANT/felicium/model"
	"github.com/robfig/revel/cache"
)

func tmpCache(m model.IsModel) error {
	return cache.Set(m.GetId(), m, 120*time.Second)
}
