package orm

import (
	"database/sql"
	"fmt"
	"github.com/goiiot/libmqtt/edge_gateway/initialize/server"
	"github.com/goiiot/libmqtt/edge_gateway/utils"
	_ "github.com/mattn/go-sqlite3" // SQLite 驱动
	"log"
)

var DB *DBKit

func Init() {
	// 初始化数据库工具包
	var err error
	DB, err = NewDBKit()
	if err != nil {
		log.Fatalf("Failed to initialize DBKit: %v", err)
	}
	server.RegisterOnShutdown(func() {
		_ = DB.Close()
	})
}

// DBKit 数据库工具包
type DBKit struct {
	db *sql.DB
}

// NewDBKit 初始化数据库连接
func NewDBKit() (*DBKit, error) {
	// 从环境变量读取数据库路径
	dbPath := utils.GetStrEnv("DB_PATH", "./data.DB")

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// 初始化表结构
	if err := initTable(db); err != nil {
		return nil, fmt.Errorf("failed to init table: %v", err)
	}

	return &DBKit{db: db}, nil
}

// initTable 初始化表结构
func initTable(db *sql.DB) error {
	log.Printf("初始化表结构:%s", "channels")
	// 创建 channels 表
	query := `
	CREATE TABLE IF NOT EXISTS channels (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		parent_id INTEGER,
		config TEXT,
		desc TEXT
	);
	CREATE TABLE IF NOT EXISTS devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL default 'default',
		channel_id INTEGER,
		config TEXT,
		desc TEXT
	);
	CREATE TABLE IF NOT EXISTS device_things_model (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		device_id INTEGER,
		config TEXT,
		desc TEXT
	);
`
	_, err := db.Exec(query)
	return err
}

// Close 关闭数据库连接
func (d *DBKit) Close() error {
	return d.db.Close()
}
