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

## 提案アーキテクチャ: シンプルメモリキャッシュ

### 設計方針（改訂版）

1. **Pod起動時の一括ロード**: 全コマンド（server, alarm, anniversary）起動時に、Google Sheetsから読み込む
2. **各Pod独立**: Pod間のデータ共有は考慮せず、各Podが独自のメモリキャッシュを持つ
3. **リトライ機能**: Google Sheets API呼び出しは3回リトライ。失敗時はPod起動を失敗させる
4. **並行アクセス制御不要**: 各Podのメモリ内で完結するため、ロック機構は不要
5. **データ更新方法**: Google Sheetsのデータを更新した際は、Podを再起動して反映（Kubernetesのローリングアップデート機能を活用）

### 実装案

#### 新しいパッケージ構造: `cache/`

```go
// cache/cache.go
package cache

import (
    "fmt"
    "log"
    "time"

    "github.com/Songmu/retry"
    "github.com/commojun/nyanbot/masterdata/table"
    "github.com/commojun/nyanbot/masterdata/key_value"
)

var (
    // グローバルキャッシュインスタンス（Pod内でシングルトン）
    globalCache *Cache
)

type Cache struct {
    Tables  *table.Tables
    KeyVals *key_value.KVs
}

// Initialize: Pod起動時に一度だけ呼ばれる（リトライあり）
func Initialize() error {
    if globalCache != nil {
        return nil // 既に初期化済み
    }

    var cache *Cache
    err := retry.Retry(3, 2*time.Second, func() error {
        log.Println("Attempting to load data from Google Sheets...")
        c, err := loadFromSheet()
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

// loadFromSheet: Google Sheetsから全データを読み込む
func loadFromSheet() (*Cache, error) {
    // テーブルデータを読み込み
    tables := &table.Tables{}
    if err := tables.LoadTablesFromSheet(); err != nil {
        return nil, fmt.Errorf("failed to load tables: %w", err)
    }

    // Key-Valueデータを読み込み
    kvs := &key_value.KVs{}
    if err := kvs.LoadKVsFromSheet(); err != nil {
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

// SetTestCache: テスト用にキャッシュをセット
func SetTestCache(c *Cache) {
    globalCache = c
}
```

#### 既存コードの修正

1. **`cmd/nyan/nyan.go`**
   - 全コマンド（`Server()`, `Alarm()`, `Anniversary()`, `Hello()`）の開始時に`cache.Initialize()`を呼び出す
   - `Export()`コマンドは削除（不要になる）

   ```go
   func Server() {
       if err := cache.Initialize(); err != nil {
           log.Fatal(err)
       }

       server, err := nyanbot.NewServer()
       if err != nil {
           log.Fatal(err)
       }
       // ...
   }

   func Alarm() {
       if err := cache.Initialize(); err != nil {
           log.Fatal(err)
       }

       alm, err := alarm.New()
       // ...
   }

   func Anniversary() {
       if err := cache.Initialize(); err != nil {
           log.Fatal(err)
       }

       anniv, err := anniversary.New()
       // ...
   }
   ```

2. **`masterdata/table/table.go`**
   - `RestoreFromRedis()`削除
   - `SaveToRedis()`削除
   - `Initialize()`削除（Redisへの保存が不要になる）

3. **`masterdata/table/alarm.go`**
   ```go
   func Alarms() ([]Alarm, error) {
       return cache.GetAlarms(), nil
   }
   ```

4. **`masterdata/table/anniversary.go`**
   ```go
   func Anniversaries() ([]Anniversary, error) {
       return cache.GetAnniversaries(), nil
   }
   ```

5. **`masterdata/key_value/key_value.go`**
   - `SaveToRedis()`削除
   - `Initialize()`修正（Redisへの保存を削除）

6. **`app/linebot/linebot.go`**
   ```go
   func (bot *LineBot) TextMessageWithRoomKey(msg string, roomKey string) error {
       roomID, err := cache.GetRoomID(roomKey)
       if err != nil {
           return err
       }
       return bot.TextMessageWithRoomID(msg, roomID)
   }
   ```

7. **`app/alarm/alarm.go`**
   - `table.Alarms()`の呼び出しは変更不要（内部実装が変わるだけ）

8. **`app/anniversary/anniversary.go`**
   - `table.Anniversaries()`の呼び出しは変更不要（内部実装が変わるだけ）

### データ更新時の運用方法

Google Sheetsのデータを更新した際に、Podにデータを反映させる方法:

#### 方法1: Podの再起動（推奨）

```bash
# Deployment全体をローリング再起動
kubectl rollout restart deployment/nyanbot-server
kubectl rollout restart deployment/nyanbot-alarm
kubectl rollout restart deployment/nyanbot-anniversary

# または、Makefile経由で
make restart/server
make restart/alarm
make restart/anniversary
```

#### 方法2: 新しいデプロイメント

```bash
# イメージタグを変更せずに、環境変数やアノテーションを更新して再デプロイ
kubectl patch deployment nyanbot-server -p \
  "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"reloaded-at\":\"$(date +%s)\"}}}}}"
```

#### CronJob（alarm, anniversary）の場合

CronJobは次回の実行時に自動的に新しいPodが起動するため、データ更新後は次回実行時に自動反映されます。即座に反映したい場合は、手動でJobを実行:

```bash
kubectl create job --from=cronjob/nyanbot-alarm manual-alarm-$(date +%s)
```

## 移行ステップ（改訂版）

### フェーズ1: メモリキャッシュの実装と基本機能の置き換え

**目的**: Redisを使わない新しいキャッシュ機構を実装

1. `cache/cache.go`パッケージを作成
   - `Initialize()`: リトライ機能付きで3回試行
   - `loadFromSheet()`: Google Sheetsから読み込むプライベート関数
   - `GetAlarms()`, `GetAnniversaries()`, `GetRoomID()`: データ取得関数
   - `SetTestCache()`: テスト用モック注入関数

2. `cmd/nyan/nyan.go`を修正
   - 全コマンドの先頭に`cache.Initialize()`を追加
   - `Export()`コマンドを削除

3. `masterdata/table/`を修正
   - `alarm.go`の`Alarms()`を`cache.GetAlarms()`を使うように変更
   - `anniversary.go`の`Anniversaries()`を`cache.GetAnniversaries()`を使うように変更
   - `table.go`から`SaveToRedis()`, `RestoreFromRedis()`を削除

4. `app/linebot/linebot.go`を修正
   - `TextMessageWithRoomKey()`を`cache.GetRoomID()`を使うように変更

5. ローカル環境でテスト実行
   ```bash
   make testall
   go run cmd/nyan/nyan.go server  # 起動確認
   go run cmd/nyan/nyan.go alarm   # アラーム実行確認
   ```

### フェーズ2: Redis関連コードの完全削除

**目的**: Redisパッケージと依存関係を全て削除

1. `app/redis/redis.go`の削除（または使用箇所がゼロになったことを確認）

2. `masterdata/table/table.go`の`Initialize()`を削除
   - これはRedisへの保存のみを行っていたため不要

3. `masterdata/key_value/key_value.go`の`SaveToRedis()`を削除

4. `go.mod`からRedis関連の依存を削除（必要に応じて）
   ```bash
   go mod tidy
   ```

5. 全テストを再実行
   ```bash
   make testall
   ```

### フェーズ3: Kubernetes構成の更新

**目的**: Kubernetesリソースからredisを削除

1. `k8s/`配下のRedis関連リソースを削除
   - Redis Deployment
   - Redis Service
   - Export CronJob

2. `Makefile`からRedis関連コマンドを削除
   - `redis-cli`コマンド
   - Redisに関連する`init`処理

3. アプリケーションのDeployment/CronJob定義から環境変数`NYAN_REDIS_HOST`を削除

4. `constant/constant.go`から`NYAN_REDIS_HOST`定義を削除

5. デプロイして動作確認
   ```bash
   make deploy
   make logs/server
   ```

### フェーズ4: ドキュメント・テスト・運用手順の整備

**目的**: ドキュメント更新と運用手順の整備

1. `README.md`と`CLAUDE.md`のRedis関連記述を更新
   - アーキテクチャ図の更新
   - データフローの説明を「Google Sheets → メモリキャッシュ」に変更
   - `export`コマンドの説明を削除
   - データ更新時の運用手順（Pod再起動）を追加

2. `Makefile`にPod再起動コマンドを追加
   ```makefile
   restart/server:
       kubectl rollout restart deployment/nyanbot-server

   restart/alarm:
       kubectl rollout restart deployment/nyanbot-alarm

   restart/anniversary:
       kubectl rollout restart deployment/nyanbot-anniversary

   restart/all:
       make restart/server
       make restart/alarm
       make restart/anniversary
   ```

3. テストの整備
   - `cache`パッケージのテストを追加
   - `SetTestCache()`を使ってモックデータを注入するテストを追加

4. 最終確認
   - 全テストがパスすること
   - ローカル環境で全コマンドが動作すること
   - Kubernetes環境で本番動作確認
   - データ更新→Pod再起動のフローを確認

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
   - **対策**: Google Sheetsが真のデータソース（Single Source of Truth）なので問題なし
   - Podが再起動しても、起動時に最新データを読み込む

### 2. **複数Pod間でデータ不整合の可能性**
   - **対策**: 各Podが独立してGoogle Sheetsから読み込むため、起動タイミングによって若干のデータ差が生じる可能性がある
   - しかし、データ更新頻度は低いため、実運用上の問題は小さい
   - データ更新時は`kubectl rollout restart`でローリング再起動することで、全Podに最新データを反映

### 3. **メモリ使用量の増加**
   - **対策**: 現状のデータ量は少ない（アラーム・記念日で数十件程度）ため、影響は微小（数KB～数MB程度）
   - Redisのネットワークオーバーヘッドがなくなるトレードオフとして十分価値がある

### 4. **Google Sheets APIのレート制限**
   - **対策**: Pod起動時のみアクセスするため、通常運用では問題なし
   - Pod再起動は手動実行のため、APIレート制限に達する心配はない

## リスク評価

| リスク | 発生確率 | 影響度 | 対策 |
|--------|----------|--------|------|
| メモリ不足 | 低 | 低 | データ量は少なく（数十件）、問題なし |
| 複数Pod間の不整合 | 中 | 低 | データ更新頻度が低いため影響小。必要時はPod再起動 |
| Google Sheets APIレート制限 | 低 | 中 | 起動時のみアクセス（1回のみ）。再読み込みは手動実行 |
| Google Sheets API障害時の起動失敗 | 低 | 高 | 3回リトライ実装。失敗時はPod起動失敗でKubernetesが再試行 |

## 成功基準

- ✅ Redisコンテナの削除が完了
- ✅ すべてのテストがパス
- ✅ アラーム・記念日機能が正常動作
- ✅ ドキュメントが更新されている
- ✅ Pod再起動によるデータ更新フローが確立している
- ✅ Makefileに`restart/*`コマンドが追加されている

## まとめ

この計画により、nyanbotは以下を実現します:

1. **シンプルなアーキテクチャ**: Redisを削除し、Pod起動時にGoogle Sheetsから読み込むメモリキャッシュのみでデータ管理
2. **低運用コスト**: Redisコンテナが不要になり、インフラリソースとメンテナンスコストが削減
3. **高速アクセス**: ネットワークI/O（Redis接続）なしのメモリ直接アクセス
4. **開発体験の向上**: ローカル開発時にRedis起動不要。`export`コマンド実行も不要
5. **各Pod独立動作**: Pod間のデータ同期を考慮しない、シンプルな設計

### 技術的な改善点

- **並行制御の削除**: `sync.RWMutex`などの複雑なロック機構が不要
- **リトライ機能**: Google Sheets API呼び出しに3回リトライを実装し、一時的な障害に対応
- **テスト容易性**: `SetTestCache()`でモックデータを簡単に注入可能
- **Kubernetesネイティブな運用**: Pod再起動によるデータ更新は、Kubernetesの標準的な運用パターン

### 適用条件

この設計は以下の条件下で最適です:

- データ量が少ない（数十～数百件程度）
- データ更新頻度が低い
- 読み取り専用運用
- Google SheetsがSingle Source of Truth

nyanbotはこれらの条件を全て満たしており、Redisというオーバースペックなミドルウェアを排除することで、保守性と運用性を大幅に向上させます。
