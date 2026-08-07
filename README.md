# Go Feature Template

[laravel-feature-template](../laravel-feature-template) の Go 版です。Go を体系的に学ぶため、自分が慣れているLaravelとの比較対象として作りました。特定プロジェクトのビジネスロジックは含まれておらず、`Task` という汎用的なサンプル機能(CRUD)だけを同梱しています。新規プロジェクトを始めるときはこのリポジトリをベースに `internal/features/task` を消して自分たちの Feature を追加してください。

## 思想

Laravel 版と同じく、技術的な層(Controller / Service / Repository 等)ではなく機能(Feature)を第一の分割単位とし、1機能 = 1ディレクトリ(1パッケージ)で、Feature 内のコードは原則そのディレクトリの中に閉じ込める構造にしています。Go の `internal/` はこの構成と相性がよく、`internal/features/*` の外からは import できません。

## 技術スタック

- Go 1.25
- MySQL 8.4(ローカル開発用、Docker Compose で起動)
- [chi](https://github.com/go-chi/chi) — 軽量ルーター
- [GORM](https://gorm.io/) — Eloquent に最も近い ActiveRecord 風 ORM
- [go-playground/validator](https://github.com/go-playground/validator) — 構造体タグによるバリデーション
- [golang-migrate](https://github.com/golang-migrate/migrate) — SQL マイグレーション管理
- [swaggo/swag](https://github.com/swaggo/swag) — コードのコメントアノテーションから OpenAPI を生成
- [testify](https://github.com/stretchr/testify) + `net/http/httptest` — テスト
- golangci-lint / gofmt / goimports — 静的解析・フォーマット
- Docker / Docker Compose — ローカル開発環境

## ディレクトリ構成

```
.
├── compose.yaml              # ローカル開発環境の定義
├── docker/app/                # Dockerfile(local/builder/runtime のマルチステージ)
├── Makefile                   # よく使う操作のショートカット
├── .env.example                # compose.yaml が参照する環境変数
├── migrations/                 # golang-migrate 用の SQL マイグレーション(一元管理)
├── docs/                       # swaggo が生成する OpenAPI ドキュメント
├── cmd/
│   ├── api/main.go             # HTTP サーバーのエントリポイント
│   └── seed/main.go            # サンプルデータ投入コマンド
└── internal/
    ├── features/                # ここが本体。1機能 = 1パッケージ(ディレクトリ)
    │   └── task/                 # サンプル Feature(CRUD)
    ├── models/                   # GORM モデル(2つ以上の Feature から参照される)
    ├── apperror/                 # ドメインエラー型
    ├── httpx/                    # JSON エンコード/デコード・バリデーション・エラーレンダリング
    ├── api/
    │   ├── router.go             # ルーター組み立て(ミドルウェア・ヘルスチェック・Swagger UI)
    │   └── routes.go             # 全 Feature のルート一元登録(routes/api.php 相当)
    ├── config/                   # 環境変数の読み込み
    ├── database/                 # *gorm.DB の初期化
    └── testutil/                 # テスト用の in-memory DB・ルーター構築
```

## Package by Feature の規約

1機能につき `internal/features/{feature}/` を1パッケージ作り、その機能に関するものは原則すべてそのディレクトリの中に閉じ込めます。ただし **ルーティングとマイグレーションは例外** で、`internal/api/routes.go` と `migrations/` に一元管理します(Laravel 版が Controller 等は `app/Features/{Feature}/` に閉じ込めつつ、`routes/api.php` と `database/migrations/` だけは一元管理しているのと同じ方針)。

```
internal/features/{feature}/
├── handler.go     # 薄い Handler。1メソッド = 1アクション(Controller 相当)
├── request.go     # リクエストの型 + バリデーション(FormRequest + Input DTO 相当)
├── response.go     # レスポンスの型と変換関数(Resource 相当)
└── usecase.go      # ビジネスロジック本体
```

**リクエストの流れ**: `Router → Middleware → Handler → Request(検証) → UseCase → Model → Response`

**ルール**

- Repository パターンは使わない。UseCase から `*gorm.DB` を直接操作する(Laravel 版が Eloquent を直接叩くのと同じ方針)。
- ドメイン固有のエラーは `internal/apperror.AppError` を使う。`internal/httpx.Error()` が code/message/status から一貫した JSON エラーレスポンスに変換する。
- 「god Service」は作らない。複数 Feature にまたがる共通処理が本当に必要になったときだけ `internal/shared/` のようなパッケージを切り出す。
- DI コンテナは使わない。`main.go` / `internal/api/routes.go` で `task.NewHandler(task.NewUseCase(db))` のように手で組み立てる(下記「Laravel との対応」参照)。

新しい Feature を追加するときは `internal/features/task/` をコピーしてリネームするのが一番早い方法です。

## セットアップ

```bash
git clone <this-repo>
cd go-feature-template

cp .env.example .env

make build
make up
make migrate
```

`http://localhost:8000` で API が起動します(ポートは `.env` の `APP_PORT` で変更可能)。

Swagger UI は `http://localhost:8000/swagger/index.html` で確認できます。

## よく使うコマンド(Makefile)

| コマンド | 内容 |
| --- | --- |
| `make up` / `make down` | ローカル環境の起動・停止 |
| `make bash` | app コンテナに入る |
| `make migrate` / `make migrate-down` / `make fresh` | マイグレーション関連 |
| `make seed` | サンプルデータ投入 |
| `make test` | `go test ./...` 実行 |
| `make format` / `make format-test` | gofmt + goimports によるフォーマット・チェックのみ |
| `make lint` | golangci-lint による静的解析 |
| `make openapi` | swaggo で `docs/` を再生成 |
| `make ci` | `format-test` → `lint` → `test` を一括実行(CI と同じチェック) |

`make help` で一覧を表示できます。

## テストの書き方

- テストは `internal/features/{feature}/` 配下に、エンドポイント単位のファイルとして置きます(`index_test.go`, `store_test.go` など)。Laravel 版の `tests/Feature/Api/{Feature}/` と同じ「1ファイル1エンドポイント」規約です。
- `internal/testutil.NewDB` が in-memory SQLite(`glebarez/sqlite`、CGO 不要)+ `AutoMigrate` で毎テスト独立した DB を用意するので、Docker なしで高速に実行できます。本番のスキーマは `migrations/*.sql`(MySQL 用)が一元管理しており、テストは GORM モデルからの AutoMigrate という別経路を使う点に注意してください。
- `Handler → Request → UseCase → DB` を `httptest` 経由で通しで検証する統合テストとし、UseCase 単体のユニットテストは基本的に書きません。

## CI

`.github/workflows/ci.yml` で `gofmt` チェック → `golangci-lint` → `go test -race` → OpenAPI ドキュメントの鮮度チェックを PR ごとに実行します。デプロイ関連のワークフローは含まれていないので、デプロイ先に合わせて別途追加してください。

## Laravel 版との対応

Go を体系的に学ぶための比較用に、Laravel 版の各要素が Go 版でどう表現されているかをまとめています。

| Laravel (`laravel-feature-template`) | Go (`go-feature-template`) | 補足 |
| --- | --- | --- |
| `bootstrap/app.php`(ミドルウェア設定・ルーティングの配線) | `internal/api/router.go` | どちらも「何を挟むか・どこにマウントするか」の配線担当で、具体的なエンドポイント一覧は持たない |
| `routes/api.php`(エンドポイント一覧そのもの) | `internal/api/routes.go` | Laravel 版が Controller 等は Feature ディレクトリに閉じ込めつつルーティングだけ一元管理しているのと同じ方針を踏襲。ちなみにこれは言語の制約ではなく設計判断で、Laravel 側も `app/Features/Task/routes.php` を作って読み込む分散スタイルにできる(`nwidart/laravel-modules` 等が採用) |
| `TaskController`(Controller) | `Handler`(`handler.go`) | 薄いレイヤという役割は同じ |
| `FormRequest` + `Inputs/`(spatie/laravel-data の `Optional`) | `Request` 構造体(`json` + `validate` タグ) | FormRequest は `Illuminate\Http\Request` に紐づくため Input DTO と分かれているが、Go の構造体はもともと transport 非依存なので1つの型に統合できる |
| PATCH の「未指定」と「明示的な null」を区別する `Optional` | ジェネリクスの `Field[T]`(`field.go`) | `encoding/json` は JSON に存在しないキーの `UnmarshalJSON` を呼ばないことを利用して同じ問題を解いている |
| `UseCases/`(1 execute() ずつ) | `usecase.go` の `UseCase` 構造体 | Repository パターンを使わず ORM を直接叩くという方針も同じ |
| `Resources/`(`JsonResource`) | `response.go`(`toTaskResponse`) | `{"data": ...}` エンベロープも踏襲 |
| Eloquent Model | GORM Model(`internal/models/task.go`) | ActiveRecord 風という点は共通。GORM はマジックが少ない分やや冗長 |
| `AppDomainException` + `ApiErrorResponse` | `apperror.AppError` + `httpx.Error()` | Go に例外機構は無いので、UseCase は panic ではなく `error` を戻り値として返す |
| `ForceJsonRequest` / `CamelCaseJsonResponse` ミドルウェア | なし | Handler は最初から JSON しか書かないし、`json:"dueDate"` のような構造体タグがそのまま wire 上のキーになるので、後から snake_case→camelCase に変換する層が要らない |
| `database/migrations`(Laravel Migration、PHP で記述) | `migrations/*.sql` + golang-migrate | 「マイグレーションを一元管理する」という思想は同じ。Go 版は生 SQL |
| `RefreshDatabase` + SQLite in-memory | `internal/testutil.NewDB`(SQLite in-memory + `AutoMigrate`) | 高速なテストのために本番と違う DB を使う、という狙いは同じ |
| サービスコンテナによる自動 DI | `main.go`/`router.go` での手動ワイヤリング | Go には IoC コンテナが無いのが一般的。`task.NewHandler(task.NewUseCase(db))` のように呼び出し側が明示的に依存を組み立てる |
| Laravel Pint | `gofmt` + `goimports` | どちらも「議論の余地がない」フォーマッタ |
| Larastan(PHPStan) | golangci-lint | Go は静的型付けなので、PHP ほど静的解析に頼らなくても取れるバグの割合は多い |
| Scramble(コードから自動で OpenAPI 生成) | swaggo/swag(コメントアノテーションから生成) | Go 版は `@Summary` 等のコメントを書く必要があり、完全自動ではない |
| `php artisan serve` / 本番の PHP-FPM + nginx | 単一の Go バイナリ(`net/http`) | コンパイル言語であり、Go の HTTP サーバーは標準ライブラリだけで本番運用できるレベルにあるため、PHP-FPM と nginx のようなプロセス分割が要らない。`compose.yaml` の `app` サービスが1つで済むのもこのため |
| `.env` + `config:cache` | 起動時に一度だけ環境変数を読む `config.Load()` | Go はリクエストごとにプロセスを起動し直さないので、そもそも「設定をキャッシュする」という概念がない |

## フレームワークを使わない理由

Laravel は「これ一択」という強いフレームワーク文化がありますが、Go の Web API 開発は標準ライブラリ (`net/http`) を土台に、ルーティングだけ軽量ライブラリ(この template では chi)を足すスタイルが主流です。`gin` や `echo` のようなフルスタックフレームワークもありますが、DI・ORM・認証などは個別のライブラリを組み合わせるのが一般的で、Laravel ほど「フレームワークが全部やってくれる」わけではありません。この template も chi 以外は素の Go の標準ライブラリ(`net/http`, `encoding/json`, `log/slog` など)を中心に構成しています。
