package data

import (
	"kratos/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	db *MemoryDB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	mDB := NewMemoryDB()
	mDB.userDB["1111"] = "John Doe"
	mDB.userDB["2222"] = "Jane Smith"

	return &Data{
		db: mDB,
	}, cleanup, nil
}

// MemoryDB is a mock in-memory database.
type MemoryDB struct {
	userDB map[string]string
}

func NewMemoryDB() *MemoryDB {

	return &MemoryDB{
		userDB: make(map[string]string),
	}
}
