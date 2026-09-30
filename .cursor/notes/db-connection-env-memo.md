# DB接続エラー（Unix ドメインソケット接続失敗）の原因と対処法

## 発生したエラー

`make up`（`go run ./cmd/main.go`）実行時に以下のエラーで起動に失敗した。

```text
[error] failed to initialize database, got error failed to connect to `user=hirano-ta database=`: /private/tmp/.s.PGSQL.5432 (/private/tmp): dial error: dial unix /private/tmp/.s.PGSQL.5432: connect: no such file or directory
failed to connect to `user=hirano-ta database=`: /private/tmp/.s.PGSQL.5432 (/private/tmp): dial error: dial unix /private/tmp/.s.PGSQL.5432: connect: no such file or directory
exit status 1
```

---

## 結論（対処法）

`Makefile` の `up` コマンドに `GO_ENV=dev` が付いていなかったことが原因。

```makefile
# 修正前
up:
	go run ./cmd/main.go

# 修正後
up:
	GO_ENV=dev go run ./cmd/main.go
```

またはターミナル実行時に手動で環境変数を指定する:

```bash
GO_ENV=dev make up
```

---

## なぜこのエラーが発生したのか（メカニズム）

### 1. `GO_ENV=dev` が無いと `.env` がロードされない

`db/db.go` では、環境変数 `GO_ENV` が `"dev"` の場合のみ `godotenv.Load()` を実行する実装になっている。

```go
func NewDB() *gorm.DB {
	if os.Getenv("GO_ENV") == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}
    // ...
```

- Go 標準の `os.Getenv` は、環境変数が未設定の場合に例外やエラーを投げず、**空文字列 `""` を返す** 仕様。
  - 一次ソース: [Go 公式ドキュメント - os.Getenv](https://pkg.go.dev/os#Getenv)
- そのため `make up` だと `if` 文の条件を満たさず、`.env` ファイルが一切読み込まれない。

### 2. 接続文字列（DSN）が空のパラメータで生成される

`.env` が読まれないため、後続の `os.Getenv("POSTGRES_USER")` などもすべて空文字 `""` になる。

```go
url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PW"),
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_PORT"),
	os.Getenv("POSTGRES_DB"),
)
```

結果として、`url` は `postgres://:@:/` という空の URL 文字列になる。

### 3. PostgreSQL ドライバが Unix ドメインソケット & OS ユーザー名にフォールバックする

PostgreSQL の接続ライブラリ（libpq / pgx）には以下の公式仕様が存在する。

1. **ホスト名が省略（空）の場合**:
   TCP 接続ではなく、ローカルの Unix ドメインソケット（Mac の場合 `/private/tmp/.s.PGSQL.5432` など）を探して接続を試みる。
2. **ユーザー名が省略（空）の場合**:
   現在ログインしている OS のユーザー名（今回は Mac のアカウント名 `hirano-ta`）をデフォルト値として自動補完する。

- 一次ソース: [PostgreSQL 15 公式ドキュメント - 34.1. Database Connection Control Functions](https://www.postgresql.org/docs/15/libpq-connect.html#LIBPQ-PARAMKEYWORDS)
  > _"If the host is not specified, libpq will connect to a local server using a Unix-domain socket..."_
  > _"If user is omitted, it defaults to the user name of the OS user running the application."_

### 4. Docker 側の PostgreSQL と噛み合わずエラーになる

- Docker コンテナ上の PostgreSQL は、TCP ポート `5434:5432` で待ち受けている。
- しかしアプリケーション側は「ホスト名が空」だったため、Mac 本体（ホスト OS 側）のファイル `/private/tmp/.s.PGSQL.5432` を開こうとした。
- Mac 本体の `/private/tmp/` には PostgreSQL のソケットファイルが存在しないため、`no such file or directory` で接続失敗となった。

---

## TS / PHP との比較・学び

- **Node.js (`dotenv`) / PHP (`vlucas/phpdotenv`)**:
  これらの言語と同様、Go でも設定ファイル（`.env`）は言語組み込みではなく外部ライブラリ（`godotenv`）で明示的に読み込む必要がある。
- **Go の `os.Getenv`**:
  存在しないキーに対して `null` や `undefined` ではなく `""`（ゼロ値の空文字列）を返すため、タイポや設定漏れがサイレントに進行しやすい点に注意が必要。
