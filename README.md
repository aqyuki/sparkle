# README

Sparkleは，Goで実装されたBotで便利な機能から，ネタよりな機能まで様々な機能を提供しています．

## Features

> [!WARNING]
> 現在，**Sparkle**は非安定版です．以下の機能は未実装，もしくは仕様が破壊的に変更される可能性があります．

- メッセージリンクを展開する．
- チャンネルリンクを展開する．
- VC Diff
- Alias 機能

## How to Run
> [!NOTE]
> Dockerでの実行のみサポートしています．

1. このリポジトリをローカルにクローンする
2. Discord Developer PortalでDiscord Botのトークンを作成してください．
3. `.env.sample`を`.env`に改名し，`DISCORD_TOKEN`にDiscord Botのトークンを指定してください．
4. `docker compose up -d`で起動できます．
