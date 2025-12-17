# nyanbot
linebotを作りたい

- スプレッドシート上でcrontab的な使用感のデータを設定し定期的にメッセージを送る
- 相手の発言に応じていろいろな返答をしてくれる

## Architecture

This bot loads data from a Google Sheet into memory on startup. Data is cached in memory for the lifetime of the pod.

## Data Updates

To update the data from the Google Sheet, the corresponding Kubernetes pods must be restarted. The Makefile contains helper commands for this:

- `make restart/server`
- `make restart/alarm`
- `make restart/anniversary`
- `make restart/all`
