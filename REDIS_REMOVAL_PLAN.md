# Redis削減・メモリキャッシュ移行計画

## 概要

現在のnyanbotアーキテクチャでは、Google SheetsからRedisにデータをエクスポートし、各機能がRedisから読み取る構成となっています。しかし、データ量が少なく、Redisはオーバースペックであると判断されました。

この計画では、**Pod起動時にGoogle Sheetsから一度だけデータを読み込み、メモリ上にキャッシュする方式**に移行することで、Redisを削除し、システムを簡素化します。

## 現状分析

### Redisの現在の使用状況

1. **テーブルデータ** (`masterdata/table/`)
   - `Alarm`: アラーム設定（JSON文字列として保存）
   - `Anniversary`: 記念日設定（JSON文字列として保存）

2. **Key-Valueデータ** (`masterdata/key_value/`)
   - `room`: ルームキー → LINE Room IDのマッピング（Redisハッシュ）
   - `nickname`: ニックネーム（Redisハッシュ）
   - `testkv`: テスト用データ（Redisハッシュ）

### データフロー

**現状:**
```
起動時: Google Sheets → `export`コマンド → Redis
実行時: Redis → アプリケーション
```

**移行後:**
```
起動時: Google Sheets → メモリキャッシュ
実行時: メモリキャッシュ → アプリケーション
```

### Redisに依存しているコンポーネント

1. **`cmd/nyan/nyan.go`**: `Export()`コマンドでRedis接続確認とデータ保存
2. **`masterdata/table/table.go`**: `RestoreFromRedis()`, `SaveToRedis()`
3. **`masterdata/key_value/key_value.go`**: `SaveToRedis()`、Redis HSetでkey-value保存
4. **`app/linebot/linebot.go`**: `TextMessageWithRoomKey()`でRedisからルームID取得（42-54行目）
5. **`app/alarm/alarm.go`**: `table.Alarms()`でRedisから読み込み
6. **`app/anniversary/anniversary.go`**: `table.Anniversaries()`でRedisから読み込み

## 提案アーキテクチャ: グローバルメモリキャッシュ

### 設計方針

1. **Pod起動時の一括ロード**: サーバー起動時に全データをGoogle Sheetsから読み込む
2. **グローバルシングルトンキャッシュ**: パッケージレベルの変数でデータを保持
3. **オンデマンド再読み込み**: 必要に応じて手動でデータを再取得できる関数/エンドポイントを提供
4. **並行アクセス安全性**: 読み取り専用運用を前提とし、更新時のみロックを使用

### 実装案

#### 新しいパッケージ構造: `cache/`

```go
// cache/cache.go
package cache

import (
    "sync"
    "github.com/commojun/nyanbot/masterdata/table"
    "github.com/commojun/nyanbot/masterdata/key_value"
)

var (
    // グローバルキャッシュインスタンス
    globalCache *Cache
    once        sync.Once
    mu          sync.RWMutex
)

type Cache struct {
    Tables   *table.Tables
    KeyVals  *key_value.KVs
}

// Initialize: Pod起動時に一度だけ呼ばれる
func Initialize() error {
    var err error
    once.Do(func() {
        globalCache = &Cache{}
        err = globalCache.Reload()
    })
    return err
}

// Reload: Google Sheetsから再読み込み
func (c *Cache) Reload() error {
    mu.Lock()
    defer mu.Unlock()

    // テーブルデータを読み込み
    tables := &table.Tables{}
    if err := tables.LoadTablesFromSheet(); err != nil {
        return err
    }

    // Key-Valueデータを読み込み
    kvs := &key_value.KVs{}
    if err := kvs.LoadKVsFromSheet(); err != nil {
        return err
    }

    c.Tables = tables
    c.KeyVals = kvs
    return nil
}

// Get: グローバルキャッシュを取得
func Get() *Cache {
    mu.RLock()
    defer mu.RUnlock()
    return globalCache
}

// GetAlarms: アラームデータを取得
func GetAlarms() []table.Alarm {
    return Get().Tables.Alarms
}

// GetAnniversaries: 記念日データを取得
func GetAnniversaries() []table.Anniversary {
    return Get().Tables.Anniversaries
}

// GetRoomID: ルームキーからルームIDを取得
func GetRoomID(roomKey string) (string, error) {
    roomID, ok := Get().KeyVals.Rooms[roomKey]
    if !ok {
        return "", fmt.Errorf("room key not found: %s", roomKey)
    }
    return roomID, nil
}
```

#### 既存コードの修正

1. **`cmd/nyan/nyan.go`**
   - `Server()`関数の開始時に`cache.Initialize()`を呼び出す
   - `Export()`コマンドは削除、または`cache.Get().Reload()`に変更

2. **`masterdata/table/table.go`**
   - `RestoreFromRedis()`削除
   - `SaveToRedis()`削除
   - `Alarms()`, `Anniversaries()`は`cache.GetAlarms()`などに置き換え

3. **`masterdata/key_value/key_value.go`**
   - `SaveToRedis()`削除

4. **`app/linebot/linebot.go`**
   ```go
   func (bot *LineBot) TextMessageWithRoomKey(msg string, roomKey string) error {
       roomID, err := cache.GetRoomID(roomKey)
       if err != nil {
           return err
       }
       return bot.TextMessageWithRoomID(msg, roomID)
   }
   ```

5. **`app/alarm/alarm.go`**
   ```go
   func New() (*AlarmManager, error) {
       alms := cache.GetAlarms()
       // ...
   }
   ```

6. **`app/anniversary/anniversary.go`**
   ```go
   func New() (*AnniversaryManager, error) {
       annivs := cache.GetAnniversaries()
       // ...
   }
   ```

### 再読み込みエンドポイントの追加

データを更新した際に、Podを再起動せずに反映するためのエンドポイント:

```go
// api/reload/reload.go
package reload

import (
    "net/http"
    "github.com/commojun/nyanbot/cache"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    if err := cache.Get().Reload(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write([]byte("Cache reloaded successfully"))
}
```

## 移行ステップ

### フェーズ1: メモリキャッシュの実装（Redisと並行稼働）

1. `cache/cache.go`パッケージを作成
2. `Initialize()`と`Get()`関数を実装
3. サーバー起動時に`cache.Initialize()`を呼び出す
4. 既存のRedis読み取り処理をキャッシュ読み取りに段階的に置き換え
5. テストを実行して動作確認

### フェーズ2: Redis依存の削除

1. `app/redis/redis.go`の使用箇所をすべて削除
2. `masterdata/table/table.go`から`SaveToRedis()`, `RestoreFromRedis()`を削除
3. `masterdata/key_value/key_value.go`から`SaveToRedis()`を削除
4. `cmd/nyan/nyan.go`から`Export()`コマンドを削除（または`Reload()`に変更）
5. 環境変数`NYAN_REDIS_HOST`を削除

### フェーズ3: Kubernetesデプロイメントの更新

1. `k8s/`配下のRedis関連リソース（Deployment, Service）を削除
2. `Makefile`からRedis関連コマンド（`redis-cli`など）を削除
3. Export用のCronJobを削除
4. アプリケーションPodのRedis環境変数を削除

### フェーズ4: ドキュメント更新

1. `README.md`と`CLAUDE.md`のRedis関連記述を削除
2. 新しいキャッシュアーキテクチャの説明を追加
3. 再読み込みエンドポイントの使い方を記載

## メリット

### 1. **運用コストの削減**
   - Redisコンテナが不要になり、インフラコストが削減
   - Kubernetesリソース（Pod、Service）が削減

### 2. **システムの簡素化**
   - データフローが単純化（Sheets → Memory）
   - `export`コマンドの実行が不要
   - Redis接続エラーのハンドリングが不要

### 3. **パフォーマンス向上**
   - ネットワークI/Oが発生しない（メモリアクセスのみ）
   - Redis接続オーバーヘッドがなくなる

### 4. **開発体験の向上**
   - ローカル開発時にRedisの起動が不要
   - データ更新後、再読み込みエンドポイントを叩くだけで反映

## デメリットと対策

### 1. **データの永続性がない**
   - **対策**: Google Sheetsが真のデータソースなので問題なし

### 2. **複数Pod間でデータ同期が必要**
   - **対策**: 再読み込みエンドポイントを全Podに対して実行するか、Podを再起動

### 3. **メモリ使用量の増加**
   - **対策**: 現状のデータ量は少ないため、影響は微小（数KB～数MB程度）

### 4. **並行書き込みがない**
   - **対策**: 読み取り専用運用が前提なので問題なし。更新時のみ`RWMutex`で保護

## リスク評価

| リスク | 発生確率 | 影響度 | 対策 |
|--------|----------|--------|------|
| メモリ不足 | 低 | 中 | データ量は少なく、モニタリングで監視 |
| 複数Pod間の不整合 | 中 | 低 | 再読み込みエンドポイントで全Pod更新 |
| Google Sheets APIレート制限 | 低 | 中 | 起動時のみアクセス、頻繁な再読み込みは避ける |

## 成功基準

- ✅ Redisコンテナの削除が完了
- ✅ すべてのテストがパス
- ✅ アラーム・記念日機能が正常動作
- ✅ ドキュメントが更新されている
- ✅ 再読み込みエンドポイントが動作

## まとめ

この計画により、nyanbotは以下を実現します:

1. **シンプルなアーキテクチャ**: Redisを削除し、メモリキャッシュのみでデータ管理
2. **低コスト**: インフラリソースの削減
3. **高速アクセス**: ネットワークI/Oなしのメモリアクセス
4. **柔軟な更新**: 再読み込みエンドポイントによるオンデマンド更新

データ量が少ないLINEボットという特性を活かし、過剰な複雑性を排除した設計に移行します。
