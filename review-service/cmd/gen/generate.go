package main

// GEN 生成代码配置
// gorm gen configure

import (
	"errors"
	"flag"
	"fmt"
	"review-service/internal/conf"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite" // ✅ 新增：sqlite 驱动
	"gorm.io/gorm"

	"gorm.io/gen"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

// 连接数据库，根据配置里的 driver/source 选择不同驱动
func connectDB(cfg *conf.Data_Database) *gorm.DB {
	if cfg == nil {
		// ✅ 去掉 .New，直接用 fmt.Errorf
		panic(fmt.Errorf("GEN:connectDB fail, need config"))
	}

	switch strings.ToLower(cfg.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect mysql db fail: %w", err))
		}
		return db
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect sqlite db fail: %w", err))
		}
		return db
	default:
		panic(errors.New("GEN:connectDB fail unsupported db driver"))
	}
}

func main() {
	// 从配置文件读取数据库连接信息
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	// 指定生成代码的具体相对目录(相对当前文件)，默认为：./query
	// 默认生成需要使用 WithContext 之后才可以查询的代码，
	// 可以通过设置 gen.WithoutContext 禁用该模式
	g := gen.NewGenerator(gen.Config{
		// 生成的 CRUD 代码会放到这里
		OutPath: "../../internal/data/query",

		// gen.WithoutContext：禁用 WithContext 模式
		// gen.WithDefaultQuery：生成一个全局 Query 对象 Q
		// gen.WithQueryInterface：生成 Query 接口
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true, //deleted_at 允许为 null
	})

	// 使用项目中的数据库连接配置
	g.UseDB(connectDB(bc.Data.Database))

	// 为所有表生成 Model 结构体和 CRUD 代码
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行并生成代码
	g.Execute()
}
