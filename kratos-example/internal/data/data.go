package data

import (
	"kratos-example/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
// 加入 NewDB 函数，这是 wire 需要访问的东西
var ProviderSet = wire.NewSet(NewData, NewDB, NewExampleRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db *gorm.DB) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) *gorm.DB {

	db, err := gorm.Open(mysql.Open(c.Database.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 不建议在生产环境上使用
	if err = db.AutoMigrate(); err != nil {
		panic(err)
	}

	return db
}
