# テスト方針メモ（未実装）

**状態:** リポジトリ内に `*_test.go` はまだ無い。本メモは「何を・どの順で・どう書くか」の方針のみ。実装は別タスク。

**関連:** ユーザー登録改修（DTO / Entity / Repository）、`.cursor/notes/error-handling-memo.md`（cockroachdb/errors）

---

## ゴール

- リファクタや VO 化の前後で **挙動が壊れていないこと** を短いテストで担保する
- 最初から E2E や全パッケージ一括実行は狙わない。**変更したパッケージだけ** `go test` する運用でよい

---

## レイヤー別の優先順位

| 順  | 対象                                          | 種別         | 理由                                                                         |
| --- | --------------------------------------------- | ------------ | ---------------------------------------------------------------------------- |
| 1   | `domain/model/user`（`NewUserEntity`）        | 単体         | DB・HTTP 不要。バリデーションがここに集約されている                          |
| 2   | `usecase`（`SignUp`）                         | 単体         | `IUserRepository` / `IPasswordGenerator` が interface なのでフェイクで切れる |
| 3   | `controller`（`SignUp`）                      | HTTP 単体    | Echo + `net/http/httptest` でステータスと JSON を確認                        |
| 4   | `infra/shared`（`passwordGenerator.Execute`） | 単体（任意） | bcrypt の正常系だけ。コスト定数は usecase の `DEFAULT_COST` と整合           |
| 5   | `infra/repository`（`CreateUser` 等）         | 統合         | GORM + Postgres 依存。Docker の DB など別途方針が要る                        |
| 6   | `Login` / JWT / Cookie                        | 後回し       | `SECRET` 環境変数・bcrypt 比較・JWT 署名が絡み、SignUp より重い              |

**Repository 統合テスト:** 着手前に `CreateUser` の GORM 呼び出し（ポインタ渡し等）が意図どおりかコードレビューしてから。テストで赤→実装バグ、の順でもよい。

---

## Go のテストの決まり（このプロジェクト用）

- ファイル名: `<対象>_test.go`（例: `user_entity_test.go`）
- 関数: `TestXxx` で始める。ケース分けは **テーブル駆動**（`[]struct{ name, ... }`）を基本
- パッケージ:
  - **同じパッケージ**（`package user`）… 公開 API だけ触るならこれで十分
  - **外部テスト**（`package user_test`）… 公開 API のみ触りたいとき（今は必須ではない）
- 実行例（実装後）:
  - `go test ./domain/model/user/...`
  - `go test ./usecase/...`
  - `go test ./controller/...`
- 依存追加: 最初は **標準 `testing` + 手書きフェイク** で足りる。`testify` / `mockgen` は必要になってから検討

---

## 1. Entity（`NewUserEntity`）— 先に書く候補ケース

現状のルール（`user_entity.go`）に沿った例:

| ケース     | 入力                   | 期待                                                            |
| ---------- | ---------------------- | --------------------------------------------------------------- |
| 正常       | 有効 UUID + 有効メール | error なし、`Email()` は正規化後（`ParseAddress` 後の Address） |
| ID 不正    | `uuid.Nil`             | error（メッセージ: user id must be a valid UUID）               |
| メール空   | 空白のみ / 空文字      | error（email is required）                                      |
| メール不正 | パース不能な文字列     | error（Wrap 付き invalid email address）                        |

**注意:** 将来 Email を VO に切り出す場合、同じケースを VO のテストに移すか、Entity テストを VO 経由に差し替える。

---

## 2. Usecase（`SignUp`）— フェイクで切る

**依存:**

- `repository.IUserRepository` … `CreateUser(*user.UserEntity, []byte) error`
- `shared.IPasswordGenerator` … `Execute(string, int) ([]byte, error)`

**手書きフェイクの例（方針）:**

- `CreateUser` 内で「渡された Entity の ID / Email」「passwordHash」をフィールドに保存し、テスト末尾で assert
- `Execute` は固定の `[]byte("fake-hash")` を返す（bcrypt の中身は usecase テストの対象外）

**ケース例:**

| ケース          | フェイクの振る舞い                       | 期待                                         |
| --------------- | ---------------------------------------- | -------------------------------------------- |
| 正常            | Execute 成功、CreateUser 成功            | DTO の ID が Entity と一致、CreateUser が1回 |
| Entity 失敗     | 不正メール等（`NewUserEntity` が error） | error、CreateUser は呼ばれない               |
| Execute 失敗    | Execute が error                         | error、CreateUser は呼ばれない               |
| CreateUser 失敗 | CreateUser が error                      | error                                        |

**UUID:** usecase 内で `uuid.New()` のため、期待 ID を固定値 assert しない。フェイクが受け取った ID とレスポンス ID の一致で見る。

**Login:** 別メモ化または SignUp 安定後。`GetUserByEmail` + bcrypt + `os.Getenv("SECRET")` の準備が要る。

---

## 3. Controller（`SignUp`）— httptest

**構成:**

- `usecase.IUserUsecase` のフェイクを `NewUserController` に注入
- `echo.New()` + `httptest.NewRequest` / `httptest.NewRecorder`
- Body: `{"email":"...","password":"..."}`

**ケース例:**

| ケース       | 条件                | 期待 HTTP          |
| ------------ | ------------------- | ------------------ |
| Bind 失敗    | 不正 JSON 等        | 400                |
| Usecase 成功 | フェイクが DTO 返却 | 201 + JSON に `id` |
| Usecase 失敗 | フェイクが error    | 500                |

HTTP 都合のバリデーション（controller TODO）は入ったタイミングでケース追加。

**Login / Logout:** Cookie・`API_DOMAIN` 依存。SignUp の controller テストの後。

---

## 4. リファクタとテストの順（ユーザー登録まわり）

- **機能追加・VO 化・Repository 修正** と **リファクタのみ** を混ぜない
- ユーザーが「リファクタ」と明示したとき: 対象の **現状挙動をテストで固定 → リファクタ → 同じテスト** の順（Rails 側ルールと同趣旨）
- 今回の Entity バリデーション追加のような **仕様変更** は、Entity テストを先に（または同時に）足してから usecase / controller を触る

---

## 意図的に今はやらないこと

- 全パッケージ `go test ./...` を CI 必須にする（Makefile に test ターゲットも未整備）
- Repository の Postgres 統合テスト（方針決定まで保留）
- Login 一式の網羅
- mock コード生成ツールの導入

---

## 実装タスクにするときのチェックリスト

- [ ] `domain/model/user/user_entity_test.go`
- [ ] `usecase/user_usecase_test.go`（SignUp のみ）
- [ ] `controller/user_controller_test.go`（SignUp のみ）
- [ ] （任意）`infra/shared/password_generator_test.go`
- [ ] （別イシュー）repository 統合テスト方針（DB 接続・マイグレーション・データクリーンアップ）
