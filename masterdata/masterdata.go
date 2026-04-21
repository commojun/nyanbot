package masterdata

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Songmu/retry"
	"github.com/commojun/nyanbot/app/sheet"
	"github.com/commojun/nyanbot/config"
	"github.com/commojun/nyanbot/masterdata/key_value"
	"github.com/commojun/nyanbot/masterdata/table"
)

var (
	global *MasterData
)

type MasterData struct {
	Tables  *table.Tables
	KeyVals *key_value.KVs
}

// Initialize: Pod起動時に一度だけ呼ばれる（リトライあり）
func Initialize(ctx context.Context, cfg config.Config) error {
	if global != nil {
		return nil // 既に初期化済み
	}

	var md *MasterData
	err := retry.WithContext(ctx, 3, 2*time.Second, func() error {
		log.Println("Attempting to load data from Google Sheets...")
		m, err := loadFromSheet(ctx, cfg)
		if err != nil {
			log.Printf("Failed to load from Google Sheets: %v", err)
			return err
		}
		md = m
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to initialize masterdata: %w", err)
	}

	global = md
	log.Println("MasterData initialized successfully")
	return nil
}

func loadFromSheet(ctx context.Context, cfg config.Config) (*MasterData, error) {
	s, err := sheet.New(ctx, sheet.Config{
		Email:        cfg.GoogleClientEmail,
		PrivateKey:   strings.ReplaceAll(cfg.GooglePrivateKey, `\n`, "\n"),
		PrivateKeyID: cfg.GooglePrivateKeyID,
		TokenURL:     cfg.GoogleTokenURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet service: %w", err)
	}

	tables, err := table.LoadTablesFromSheet(ctx, s, cfg.SheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tables: %w", err)
	}

	kvs, err := key_value.LoadKVsFromSheet(ctx, s, cfg.SheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load key-values: %w", err)
	}

	return &MasterData{
		Tables:  tables,
		KeyVals: kvs,
	}, nil
}

// GetTables: テーブルデータを取得（nil安全）
func GetTables() *table.Tables {
	if global == nil {
		return &table.Tables{}
	}
	return global.Tables
}

// GetKeyVals: Key-Valueデータを取得（nil安全）
func GetKeyVals() *key_value.KVs {
	if global == nil {
		return &key_value.KVs{}
	}
	return global.KeyVals
}

// SetTestData: テスト用にマスターデータをセット
func SetTestData(md *MasterData) {
	global = md
}
