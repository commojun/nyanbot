package cache

import (
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
	globalCache *Cache
)

type Cache struct {
	Tables  *table.Tables
	KeyVals *key_value.KVs
}

// Initialize: Pod起動時に一度だけ呼ばれる（リトライあり）
func Initialize(cfg config.Config) error {
	if globalCache != nil {
		return nil // 既に初期化済み
	}

	var cache *Cache
	err := retry.Retry(3, 2*time.Second, func() error {
		log.Println("Attempting to load data from Google Sheets...")
		c, err := loadFromSheet(cfg)
		if err != nil {
			log.Printf("Failed to load from Google Sheets: %v", err)
			return err
		}
		cache = c
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to initialize cache after 3 retries: %w", err)
	}

	globalCache = cache
	log.Println("Cache initialized successfully")
	return nil
}

func loadFromSheet(cfg config.Config) (*Cache, error) {
	s, err := sheet.New(sheet.Config{
		Email:        cfg.GoogleClientEmail,
		PrivateKey:   strings.ReplaceAll(cfg.GooglePrivateKey, `\n`, "\n"),
		PrivateKeyID: cfg.GooglePrivateKeyID,
		TokenURL:     cfg.GoogleTokenURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet service: %w", err)
	}

	tables, err := table.LoadTablesFromSheet(s, cfg.SheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load tables: %w", err)
	}

	kvs, err := key_value.LoadKVsFromSheet(s, cfg.SheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to load key-values: %w", err)
	}

	return &Cache{
		Tables:  tables,
		KeyVals: kvs,
	}, nil
}

// Get: グローバルキャッシュを取得
func Get() *Cache {
	return globalCache
}

// GetAlarms: アラームデータを取得
func GetAlarms() []table.Alarm {
	if globalCache == nil {
		return []table.Alarm{}
	}
	return globalCache.Tables.Alarms
}

// GetAnniversaries: 記念日データを取得
func GetAnniversaries() []table.Anniversary {
	if globalCache == nil {
		return []table.Anniversary{}
	}
	return globalCache.Tables.Anniversaries
}

// GetRoomID: ルームキーからルームIDを取得
func GetRoomID(roomKey string) (string, error) {
	if globalCache == nil {
		return "", fmt.Errorf("cache not initialized")
	}
	roomID, ok := globalCache.KeyVals.Rooms[roomKey]
	if !ok {
		return "", fmt.Errorf("room key not found: %s", roomKey)
	}
	return roomID, nil
}

// GetNickname: ユーザーIDからニックネームを取得
func GetNickname(userID string) (string, error) {
	if globalCache == nil {
		return "", fmt.Errorf("cache not initialized")
	}
	nickname, ok := globalCache.KeyVals.Nickname[userID]
	if !ok {
		return "", fmt.Errorf("nickname not found for user id: %s", userID)
	}
	return nickname, nil
}

// SetTestCache: テスト用にキャッシュをセット
func SetTestCache(c *Cache) {
	globalCache = c
}
