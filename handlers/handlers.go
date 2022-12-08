package handlers

import (
	"blogpost/article_service/config"
	"blogpost/article_service/storage"
)

// Handler...
type handler struct {
	Stg storage.StorageI
	Cfg config.Config
}

// NewHandler ...
func NewHandler(stg storage.StorageI, cfg config.Config) handler {
	return handler{
		Stg: stg,
		Cfg: cfg,
	}
}
