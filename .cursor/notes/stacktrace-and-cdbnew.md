## `errors.New` とは

`github.com/cockroachdb/errors` の **`errors.New("email is required")`** は、**メッセージ文字列だけから新しい `error` を1つ作る**関数です。

- 戻り値の型は `error`（インターフェース）
- **元になる別の `error` は無い**（25行目は `mail.ParseAddress` の前なので、引き継ぐ `err` も無い）
- 呼び出し側は `if err != nil { ... }` で受け取る

公式（一次ソース）では `New` は **「creates an error with a simple error message」** で、**スタックトレースも保持する**と書かれています。  
[pkg.go.dev - errors.New](https://pkg.go.dev/github.com/cockroachdb/errors#New)

Go 標準の `import "errors"` の `errors.New` も「文字列から error を作る」点は同じですが、**cockroachdb/errors の `New` はスタック付き・詳細表示向け**などが足された版、と考えるとよいです（README でも `github.com/pkg/errors` や標準 `errors` の代替として使う、と説明されています）。

25行目の意味:

```go
return nil, errors.New("email is required")
```

- 第1戻り値 `nil` … ユーザー entity を作れなかったのでポインタは無い
- 第2戻り値 … 「メールが空（必須）」という **ドメイン側の理由** を error として返す

PHP でいう **「新しく例外を投げる（前の例外なし）」** に近いです。

---

## スタックトレースとは

**その error が「どの関数の、どの行」で作られた（返された）か** を、呼び出しの連鎖として記録したものです。

例（イメージ）:

```
NewUserEntity (user_entity.go:25)
  ← SignUp (user_usecase.go:50)
    ← SignUp (user_controller.go:45)
      ← ...
```

「どこからこの `email is required` が返ってきたか」を後から辿れるログ、というイメージです。

---

## 「残る」と何がいいの？

### 1. 原因の特定が速い

本番やログで `invalid email address` だけ見ても、**コントローラなのか entity なのか usecase なのか** が分かりにくいことがあります。  
スタックがあれば **最初に `errors.New` / `Wrap` した行** にすぐたどり着けます。

cockroachdb/errors では、普通の `err.Error()` や `%v` ではメッセージ中心で、**`%+v` でフォーマットするとスタックなど詳細が出る**、と API コメントにあります（`New` / `Wrap` の "Detail output"）。

### 2. `Wrap` と組み合わせたとき

- **メッセージ** … 人間向け（「メール不正」）
- **Unwrap チェーン** … 元の `mail.ParseAddress` の理由
- **スタック** … **あなたのコードのどこでその error を返したか**

この3つが役割分担します。25行目は「下位の err が無い」ので、**メッセージ + スタック（どこで required と判定したか）** が主な価値です。

### 3. 標準 `errors.New` との違い（このライブラリを使う理由の一つ）

標準の `errors.New` は **メッセージだけ** で、**自動では呼び出し元のファイル行は付かない**（自分で `fmt.Errorf` やログに場所を書く必要がある）。

cockroachdb/errors の `New` / `Wrap` は **error を作った時点でスタックを取る** 設計なので、**ログ・Sentry 等に詳細を出すとき**に向いています（README の Sentry / `GetSafeDetails` など）。

---

## 25行目で `New` が妥当な理由（おさらい）

ここでは **`mail.ParseAddress` がまだ失敗していない**（空文字を先に弾いている）ので、

- 引き継ぐ `err` が無い → **`Wrap` は使わない**
- 理由は自分で決めている → **`errors.New("email is required")`**

30行目は下位が `err` を返すので **`Wrap`**、25行目は **自分だけが理由を知っているので `New`**、という使い分けです。

---

## 実務上の注意（短く）

- API レスポンスに **`err.Error()` をそのまま返す** と、メッセージはユーザーに見える（今の controller はそのパターン）。スタックは通常 **クライアントには出さず、サーバログや `%+v` / Sentry 用**。
- メッセージにメールアドレス本体など **PII を入れない**（cockroachdb の `New` コメントでも、報告用に安全な文字列想定、と書かれています）。

「スタックが残る」= **デバッグ・運用で「どの行から返ったか」を追える** こと。アプリの正常系の挙動（SignUp が成功するか）は変わりません。
