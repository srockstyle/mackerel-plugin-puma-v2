# mackerel-plugin-puma-v2 クイックスタートガイド

5分で mackerel-plugin-puma-v2 を導入する手順です。

## 🚀 最速セットアップ

### 1. プラグインのインストール（30秒）

```bash
# mkr がインストールされている場合
sudo mkr plugin install srockstyle/mackerel-plugin-puma-v2
```

または

```bash
# 直接ダウンロード（Linux amd64の例）
curl -L https://github.com/srockstyle/mackerel-plugin-puma-v2/releases/latest/download/mackerel-plugin-puma-v2_linux_amd64.tar.gz | tar xz
sudo mv mackerel-plugin-puma-v2 /opt/mackerel-agent/plugins/bin/
sudo chmod +x /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2
```

### 2. Puma の設定（1分）

`config/puma.rb` に1行追加：

```ruby
# この行を追加
activate_control_app 'unix:///tmp/puma.sock'
```

Puma を再起動：

```bash
sudo systemctl restart puma
# または
bundle exec pumactl restart
```

### 3. Mackerel エージェントの設定（1分）

`/etc/mackerel-agent/mackerel-agent.conf` に追加：

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2"
```

### 4. 動作開始（30秒）

```bash
# エージェントを再起動
sudo systemctl restart mackerel-agent
```

完了！ 🎉

## 📊 確認方法

### コマンドラインで確認

```bash
# メトリクスが取得できているか確認
/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2
```

正常な出力例：
```
puma.workers	2	1692237600
puma.running	4	1692237600
puma.backlog	0	1692237600
```

### Mackerel Web UI で確認

1. [Mackerel](https://mackerel.io) にログイン
2. ホスト → カスタムメトリック
3. `puma.*` グラフが表示されていれば成功

## 🔧 よくあるエラーと対処法

### ❌ "permission denied" エラー

```bash
# ソケットファイルの権限を緩める
sudo chmod 666 /tmp/puma.sock
```

### ❌ "no such file or directory" エラー

```bash
# Puma が起動しているか確認
ps aux | grep puma

# ソケットファイルが存在するか確認
ls -la /tmp/puma.sock
```

### ❌ メトリクスが Mackerel に表示されない

```bash
# デバッグモードで実行
MACKEREL_PLUGIN_DEBUG=1 /opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2

# エージェントのログを確認
sudo journalctl -u mackerel-agent -n 50
```

## 🎯 次のステップ

### 拡張メトリクスを有効化

メモリ使用量や GC 統計も監視したい場合：

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -extended"
```

### カスタムソケットパス

デフォルト以外の場所を使用する場合：

```toml
[plugin.metrics.puma]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/var/run/puma/pumactl.sock"
```

### 複数インスタンスの監視

```toml
[plugin.metrics.puma_app1]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma_app1.sock -metric-key-prefix=app1"

[plugin.metrics.puma_app2]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-puma-v2 -socket=/tmp/puma_app2.sock -metric-key-prefix=app2"
```

## 📚 詳細情報

- [完全な導入ガイド](./mackerel-integration-guide.md)
- [メトリクス一覧](./mackerel-integration-guide.md#メトリクス一覧)
- [トラブルシューティング](./mackerel-integration-guide.md#トラブルシューティング)

## 💬 サポート

問題が解決しない場合：

- [GitHub Issues](https://github.com/srockstyle/mackerel-plugin-puma-v2/issues) で質問
- [Mackerel Slack](https://mackerel-users-jp.slack.com/) で相談

---

**所要時間: 約5分** ⏱️