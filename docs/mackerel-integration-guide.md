# mackerel-plugin-puma-v2 導入ガイド

このドキュメントでは、mackerel-plugin-puma-v2 を Mackerel エージェントに組み込む手順を説明します。

## 目次

1. [前提条件](#前提条件)
2. [インストール方法](#インストール方法)
3. [Puma の設定](#puma-の設定)
4. [プラグインの設定](#プラグインの設定)
5. [動作確認](#動作確認)
6. [トラブルシューティング](#トラブルシューティング)
7. [メトリクス一覧](#メトリクス一覧)

## 前提条件

### 必要な環境

- **Mackerel エージェント**: v0.72.0 以降
- **Puma**: 4.0 以降（6.6.1 以降推奨）
- **Go**: 1.24 以降（ソースからビルドする場合）
- **OS**: Linux, macOS（Unix ソケットをサポートする OS）

### 権限要件

- Mackerel エージェントの実行ユーザーが Puma の Unix ソケットファイルを読み取れること
- `/opt/mackerel-agent/plugins/bin/` ディレクトリへの書き込み権限（インストール時）

## インストール方法

### 方法 1: mkr を使用したインストール（推奨）

```bash
# プラグインをインストール
sudo mkr plugin install srockstyle/mackerel-plugin-puma-v2

# インストール確認
ls -la /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2
```

### 方法 2: GitHub Releases からダウンロード

```bash
# 最新版をダウンロード（例: Linux amd64）
curl -L https://github.com/srockstyle/mackerel-plugin-puma-v2/releases/latest/download/mackerel-plugin-puma-v2_linux_amd64.tar.gz -o mackerel-plugin-puma-v2.tar.gz

# 解凍
tar xzf mackerel-plugin-puma-v2.tar.gz

# プラグインディレクトリに配置
sudo mv mackerel-plugin-puma-v2 /opt/mackerel-agent/plugins/bin/
sudo chmod +x /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2
```

### 方法 3: ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/srockstyle/mackerel-plugin-puma-v2.git
cd mackerel-plugin-puma-v2

# ビルド
make build

# インストール
sudo cp bin/mackerel-plugin-puma-v2 /opt/mackerel-agent/plugins/bin/
sudo chmod +x /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2
```

## Puma の設定

### 1. Puma の設定ファイルを編集

`config/puma.rb` に以下を追加：

```ruby
# Unix ソケットでコントロールサーバーを有効化
activate_control_app 'unix:///tmp/puma.sock'

# 別のパスを使用する場合
# activate_control_app 'unix:///var/run/puma/pumactl.sock'
```

### 2. Puma を再起動

```bash
# systemd を使用している場合
sudo systemctl restart puma

# または直接再起動
bundle exec pumactl restart
```

### 3. ソケットファイルの確認

```bash
# ソケットファイルが作成されているか確認
ls -la /tmp/puma.sock

# 出力例
# srw-rw-rw- 1 puma puma 0 8月 16 10:00 /tmp/puma.sock
```

## プラグインの設定

### 1. Mackerel エージェントの設定ファイルを編集

`/etc/mackerel-agent/mackerel-agent.conf` を編集：

```toml
# 基本設定（Unix ソケット使用）
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2"

# カスタムソケットパスを使用する場合
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/var/run/puma/pumactl.sock"

# 拡張メトリクスを有効にする場合
[plugin.metrics.puma_extended]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock -extended"

# カスタムプレフィックスを使用する場合
[plugin.metrics.puma_prod]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock -metric-key-prefix=puma_prod"
```

### 2. 環境変数を使用する場合

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2"
env = { PUMA_SOCKET = "/var/run/puma/pumactl.sock" }
```

### 3. Mackerel エージェントを再起動

```bash
sudo systemctl restart mackerel-agent
```

## 動作確認

### 1. プラグインの手動実行

```bash
# 基本的な動作確認
/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock

# 出力例
puma.workers	2	1692237600
puma.booted_workers	2	1692237600
puma.running	4	1692237600
puma.pool_capacity	16	1692237600
...
```

### 2. デバッグモードで実行

```bash
# 詳細なログを出力
MACKEREL_PLUGIN_DEBUG=1 /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma.sock
```

### 3. Mackerel エージェントのログ確認

```bash
# エージェントのログを確認
sudo journalctl -u mackerel-agent -f

# または
sudo tail -f /var/log/mackerel-agent.log
```

### 4. Mackerel Web UI での確認

1. [Mackerel](https://mackerel.io) にログイン
2. 対象のホストページを開く
3. 「カスタムメトリック」タブを確認
4. `puma.*` メトリクスが表示されていることを確認

## トラブルシューティング

### よくある問題と解決方法

#### 1. Permission denied エラー

```bash
# エラー例
Error: dial unix /tmp/puma.sock: connect: permission denied
```

**解決方法**:
```bash
# ソケットファイルの権限を確認
ls -la /tmp/puma.sock

# Mackerel エージェントのユーザーを確認
ps aux | grep mackerel-agent

# 権限を修正（例）
sudo chmod 666 /tmp/puma.sock
# または、エージェントユーザーをpumaグループに追加
sudo usermod -a -G puma mackerel-agent
```

#### 2. No such file or directory エラー

```bash
# エラー例
Error: dial unix /tmp/puma.sock: connect: no such file or directory
```

**解決方法**:
```bash
# Puma が起動しているか確認
ps aux | grep puma

# コントロールサーバーが有効か確認
grep activate_control_app config/puma.rb

# Puma を再起動
sudo systemctl restart puma
```

#### 3. メトリクスが表示されない

**確認事項**:
1. プラグインが正しくインストールされているか
2. 設定ファイルの記述が正しいか
3. Mackerel エージェントが再起動されているか
4. ファイアウォールやセキュリティグループの設定

```bash
# 設定ファイルの構文チェック
mackerel-agent configtest

# dry-run モードで実行
sudo mackerel-agent -conf=/etc/mackerel-agent/mackerel-agent.conf -dry-run -once
```

#### 4. 拡張メトリクスが取得できない

```bash
# GC stats エンドポイントの確認
curl --unix-socket /tmp/puma.sock http://localhost/gc-stats

# エラーが出る場合、Puma のバージョンが古い可能性があります
```

## メトリクス一覧

### 基本メトリクス

| メトリクス名 | 説明 | 単位 | 種別 |
|------------|------|-----|------|
| `puma.workers` | ワーカープロセス数 | 個 | gauge |
| `puma.booted_workers` | 起動済みワーカー数 | 個 | gauge |
| `puma.old_workers` | 古いワーカー数（phased restart 時） | 個 | gauge |
| `puma.phase` | 現在のフェーズ番号 | - | gauge |
| `puma.backlog` | リクエストバックログ | 個 | gauge |
| `puma.running` | 実行中のスレッド数 | 個 | gauge |
| `puma.pool_capacity` | スレッドプール容量 | 個 | gauge |
| `puma.max_threads` | 最大スレッド数 | 個 | gauge |

### Puma 6.x 専用メトリクス

| メトリクス名 | 説明 | 単位 | 種別 |
|------------|------|-----|------|
| `puma.requests_count` | 処理したリクエスト総数 | 個 | counter |
| `puma.uptime` | サーバー稼働時間 | 秒 | gauge |
| `puma.thread_utilization` | スレッド使用率 | % | gauge |

### 拡張メトリクス（-extended フラグ使用時）

#### メモリメトリクス

| メトリクス名 | 説明 | 単位 | 種別 |
|------------|------|-----|------|
| `puma.memory.alloc` | 割り当て済みメモリ | MB | gauge |
| `puma.memory.sys` | システムメモリ | MB | gauge |
| `puma.memory.heap_inuse` | 使用中のヒープメモリ | MB | gauge |

#### Ruby GC メトリクス

| メトリクス名 | 説明 | 単位 | 種別 |
|------------|------|-----|------|
| `puma.ruby.gc.count` | GC 実行回数 | 回 | counter |
| `puma.ruby.gc.minor_count` | Minor GC 実行回数 | 回 | counter |
| `puma.ruby.gc.major_count` | Major GC 実行回数 | 回 | counter |
| `puma.ruby.gc.heap_available_slots` | 利用可能なヒープスロット | 個 | gauge |
| `puma.ruby.gc.heap_live_slots` | 使用中のヒープスロット | 個 | gauge |
| `puma.ruby.gc.heap_free_slots` | 空きヒープスロット | 個 | gauge |
| `puma.ruby.gc.old_objects` | Old 世代オブジェクト数 | 個 | gauge |
| `puma.ruby.gc.oldmalloc_bytes` | Old malloc バイト数 | bytes | gauge |

#### Go ランタイムメトリクス

| メトリクス名 | 説明 | 単位 | 種別 |
|------------|------|-----|------|
| `puma.gc.num_gc` | Go の GC 実行回数 | 回 | counter |
| `puma.go.goroutines` | Goroutine 数 | 個 | gauge |

## 高度な設定

### 複数の Puma インスタンスを監視

```toml
# アプリケーション 1
[plugin.metrics.puma_app1]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/var/run/puma/app1.sock -metric-key-prefix=puma_app1"

# アプリケーション 2
[plugin.metrics.puma_app2]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/var/run/puma/app2.sock -metric-key-prefix=puma_app2"
```

### アラート設定例

Mackerel の Web UI でアラートルールを設定：

1. **高負荷アラート**
   - メトリクス: `puma.thread_utilization`
   - 条件: 80% 以上が 5 分間継続
   - 重要度: Warning

2. **ワーカー異常アラート**
   - メトリクス: `puma.booted_workers`
   - 条件: 期待値を下回る
   - 重要度: Critical

3. **リクエストバックログアラート**
   - メトリクス: `puma.backlog`
   - 条件: 100 以上
   - 重要度: Warning

## まとめ

mackerel-plugin-puma-v2 を使用することで、Puma の詳細なメトリクスを Mackerel で監視できます。Unix ソケットを使用することで、セキュアで効率的な監視が可能です。

問題が発生した場合は、[GitHub Issues](https://github.com/srockstyle/mackerel-plugin-puma-v2/issues) でお問い合わせください。