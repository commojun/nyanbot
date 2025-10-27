# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

nyanbotはGo言語で書かれたLINEボットアプリケーションです。定期的なメッセージ送信、インタラクティブな応答、記念日リマインダー機能を提供します。データ管理にはGoogle Sheetsを、キャッシュ/ストレージにはRedisを使用し、Kubernetes上にデプロイされます。

## ビルドとテストコマンド

```bash
# 全テストを実行
make testall
# または
go test ./...

# マルチアーキテクチャ対応のDockerイメージをビルド
make dockerbuild VERSION=x.x.x

# 新しいリリースタグを作成
make release VERSION=x.x.x
```

## アプリケーションの実行

メインエントリポイントは`cmd/nyan/nyan.go`で、複数のサブコマンドをサポートします:

```bash
# メインHTTPサーバーを起動（LINE webhookハンドラー）
go run cmd/nyan/nyan.go server

# テスト用のhelloメッセージを送信
go run cmd/nyan/nyan.go hello

# Google SheetsからRedisにデータをエクスポート
go run cmd/nyan/nyan.go export

# アラームチェッカーを実行（cron的なスケジュール送信）
go run cmd/nyan/nyan.go alarm

# 記念日チェッカーを実行
go run cmd/nyan/nyan.go anniversary
```

## Kubernetesデプロイ

```bash
# 全サービスを初期化
make init

# サーバーのみデプロイ
make deploy

# envfileからシークレットを作成
make secret

# 特定のアプリのログを表示
make logs/server
make logs/alarm

# Redis CLIにアクセス
make redis-cli

# Podにシェルでアクセス
make shell/pod-name
```

## アーキテクチャ

### コアコンポーネント

1. **Server (`app/server/`)**: APIエンドポイントを登録し、リクエストを処理するHTTPサーバー。webhookシステムのエントリポイント。

2. **API Layer (`api/`)**: GET/POSTハンドラーを持つHTTPエンドポイントを定義:
   - `line_hook`: LINEボットメッセージ用webhookエンドポイント
   - `message`: 手動メッセージ送信エンドポイント
   - `test`: テストエンドポイント

3. **LineBot (`app/linebot/`)**: LINE Bot SDKのラッパーでメッセージ送信機能を提供:
   - `TextMessage()`: デフォルトルームに送信
   - `TextMessageWithRoomKey()`: キーでルームを指定して送信（RedisからルームIDを取得）
   - `TextMessageWithRoomID()`: 特定のルームIDに送信
   - `TextReply()`: メッセージイベントに返信

4. **Text Message Actions (`app/linebot/text_message_action/`)**: 異なるテキストメッセージのプレフィックスを処理するためのコマンドパターン実装。LINEメッセージを受信すると、プレフィックスをチェックして対応するアクション（占い、おじさんチャット、記念日検索など）を実行。プレフィックスにマッチしない場合はechoにフォールバック。

5. **Alarm System (`app/alarm/`)**: cron的なスケジュールメッセージシステム。Redis（元はGoogle Sheets）からアラーム設定を読み込み、`checker.go`を使って現在時刻がアラームスケジュールと一致するかチェック。別のKubernetes CronJobとして実行。

6. **Anniversary System (`app/anniversary/`)**: 重要な日付を追跡してリマインダーを送信。イベントからの日数/年数を計算し、定期的にリマインダーを送信。これもKubernetes CronJobとして実行。

### データフロー

1. **初期化（Exportコマンド）**:
   - Google Sheets → `masterdata/table`および`masterdata/key_value` → Redis
   - テーブルはJSON文字列として、key-valueはRedisハッシュとして保存
   - 2種類のデータ:
     - **Tables** (`masterdata/table/`): 構造化データ（Alarm、Anniversary）で複数行を持つ
     - **Key-Values** (`masterdata/key_value/`): シンプルなキーバリューペア（ルームID、ニックネーム）

2. **メッセージ処理**:
   - LINE Webhook → `api/line_hook` → LineBotがイベントをパース → Text Message Actions → レスポンス

3. **スケジュールタスク**:
   - Kubernetes CronJobs → alarm/anniversaryコマンド → Redisからデータ取得 → 条件チェック → LINEメッセージ送信

### 設定

環境変数は`constant/constant.go`で定義（必要な値は`envfile.sample`を参照）:
- LINE Bot認証情報: `NYAN_CHANNEL_SECRET`, `NYAN_ACCESS_TOKEN`, `NYAN_DEFAULT_ROOM_ID`
- Google Sheets API: `NYAN_GOOGLE_CLIENT_EMAIL`, `NYAN_GOOGLE_PRIVATE_KEY`, `NYAN_SHEET_ID`
- Redis: `NYAN_REDIS_HOST`
- Server: `NYAN_SERVER_PORT`, `NYAN_MESSAGE_TOKEN`

### 時刻処理

すべての時刻操作は`app/time_util/`を使用し、JST（日本標準時）変換を処理します。ボットは日本のタイムゾーンで動作するため、アラームと記念日のスケジューリングにはこれが重要です。

## 開発メモ

- プロジェクトは`masterdata/table/`と`masterdata/key_value/`でリフレクションを多用し、構造体タグに基づいてGoogle Sheetsからデータを動的にロードしてRedisに保存します。
- ルームIDは「ルームキー」を使って抽象化されています - Redisに保存された人間が読める形式のキーが、実際のLINEルームIDにマッピングされます。
- 新しいテキストメッセージアクションを追加する場合は、`app/linebot/text_message_action/text_message_action.go`の`actions()`スライスに追加してください。
- プロジェクトは、各主要機能（alarm、anniversaryなど）がスタンドアロンコマンドとしても、サーバーの一部としても実行できるパターンに従っています。
