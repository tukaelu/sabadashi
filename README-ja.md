# sabadashi
![Latest GitHub Release](https://img.shields.io/github/release/tukaelu/sabadashi.svg)
[![Go Report Card](https://goreportcard.com/badge/tukaelu/sabadashi)](https://goreportcard.com/report/tukaelu/sabadashi)
![Github Actions Test](https://github.com/tukaelu/sabadashi/workflows/test/badge.svg?branch=main)

![](./images/sabadashi-logo.png)

## 概要

Mackerelの任意のホストの指定した期間に投稿されたメトリックをCSVファイルに出力する非公式なコマンドラインツールです。  
ホストメトリックとサービスメトリックの取得に対応しています。

なお、このツールは大量のAPIリクエストを行うことがあります。多重実行はサービスに負荷をかける事も考えられるのでお控えください。

## インストール方法

### Homebrewでインストール

```
brew install tukaelu/tap/sabadashi
```

### バイナリを使用する

[リリースページ](https://github.com/tukaelu/sabadashi/releases)から使用する環境にあったZipアーカイブをダウンロードしてご使用ください。

## 使用方法

### ホストメトリック
```
NAME:
   sabadashi host - Retrieves host metrics

USAGE:
   sabadashi host [command options] [arguments...]

OPTIONS:
   --id value                       メトリックを取得するホストIDを指定
   --from value                     YYYYMMDD形式でメトリック取得を開始する日付を指定 （例: 20230101）
   --to value                       YYYYMMDD形式でメトリック取得を終了する日付を指定 （例: 20231231）
   --granularity value, -g value    取得するメトリックの粒度を 1m, 5m, 10m, 1h, 2h, 4h, 1d から指定 （デフォルト: 1m）
   --with-friendly-date-format, -f  フラグを有効にするとCSVの行頭に読みやすい日付のカラムを追加 （デフォルト: 追加しない）
   --help, -h                       ヘルプを表示する
```

### サービスメトリック
```
NAME:
   sabadashi service - Retrieves service metrics

USAGE:
   sabadashi service [command options] [arguments...]

OPTIONS:
   --name value, -n value           メトリックを取得するサービス名を指定
   --from value                     YYYYMMDD形式でメトリック取得を開始する日付を指定 （例: 20230101）
   --to value                       YYYYMMDD形式でメトリック取得を終了する日付を指定 （例: 20231231）
   --granularity value, -g value    取得するメトリックの粒度を 1m, 5m, 10m, 1h, 2h, 4h, 1d から指定 （デフォルト: 1m）
   --with-friendly-date-format, -f  フラグを有効にするとCSVの行頭に読みやすい日付のカラムを追加 （デフォルト: 追加しない）
   --with-external-monitors, -e     フラグを有効にすると外形監視で計測したメトリックも含める （デフォルト: 追加しない）
   --help, -h                       ヘルプを表示する
```

メトリックを取得するホストのID、YYYYMMDD形式の取得の開始日と終了日を指定します。  
コマンドオプションのfromに指定された`YYYY/MM/DD 00:00:00`から、toに指定された`YYYY/MM/DD 23:59:59`までに投稿されたメトリックを取得します。
有効な期間の範囲は460日間です。

```
# APIキーが環境変数 MACKEREL_APIKEY に設定されている場合
sabadashi host -id <your host id> -from <YYYYMMDD> -to <YYYYMMDD>

# APIキーをオプションで指定する場合
sabadashi host -apikey <your api key> -id <your host id> -from <YYYYMMDD> -to <YYYYMMDD>
```

環境変数の`MACKEREL_APIKEY`もしくは`-apikey`オプションに指定するAPIキーには参照権限が必要となります。

なおコマンドを実行すると、作業ディレクトリの配下にホストIDと開始日・終了日によって命名されたディレクトリが作成され、その中にメトリックごとのCSVファイルが出力されます。

## 注意

- 非公式なプラグインのため、ご質問はIssueやSNSなどでお願いします。
- 前述の通り、複数のホストのメトリクスを同時に取得する行為はサービスに負荷をかけることがあるためお控えください。
- ホストメトリックは[ホストのメトリック名の一覧API](https://mackerel.io/ja/api-docs/entry/hosts#metric-names)を元にメトリックを取得しています。
- サービスメトリックは[サービスのメトリック名の一覧API](https://mackerel.io/ja/api-docs/entry/services#metric-names)を元にして、外形監視を対象にするオプションが有効な場合に限り、[監視ルールの一覧](https://mackerel.io/ja/api-docs/entry/monitors#list)から外形監視で計測されたメトリックを追加してメトリックを取得しています。
- 指定した期間中にメトリックが投稿されていない場合、その時間帯のデータはCSVには出力されず、空のファイルだけが作成されることもあります。

## ライセンス

Copyright 2023 tukaelu (Tsukasa NISHIYAMA)

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
