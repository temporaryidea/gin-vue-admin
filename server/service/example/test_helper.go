package example

import (
	"os"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestEnv(t *testing.T) func() {
	// Create a temporary SQLite database file
	tmpDB := t.TempDir() + "/test.db"
	
	db, err := gorm.Open(sqlite.Open(tmpDB), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	// Auto migrate the required tables
	err = db.AutoMigrate(&example.ExaFile{}, &example.ExaFileChunk{}, &example.ExaCustomer{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Store original globals
	originalDB := global.GVA_DB
	originalLogger := global.GVA_LOG
	originalConfig := global.GVA_CONFIG
	originalVP := global.GVA_VP
	originalDBList := global.GVA_DBList
	originalActiveDBName := global.GVA_ACTIVE_DBNAME

	// Set up test globals
	global.GVA_DB = db
	global.GVA_LOG, _ = zap.NewDevelopment()
	global.GVA_VP = viper.New()
	global.GVA_CONFIG = config.Server{
		Zap: config.Zap{
			Level: "info",
		},
		System: config.System{
			DbType: "sqlite",
		},
	}
	
	// Initialize database list with test database
	global.GVA_DBList = map[string]*gorm.DB{
		"default": db,
	}
	defaultDBName := "default"
	global.GVA_ACTIVE_DBNAME = &defaultDBName

	// Return cleanup function
	return func() {
		global.GVA_DB = originalDB
		global.GVA_LOG = originalLogger
		global.GVA_CONFIG = originalConfig
		global.GVA_VP = originalVP
		global.GVA_DBList = originalDBList
		global.GVA_ACTIVE_DBNAME = originalActiveDBName
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
		os.Remove(tmpDB)
	}
}
