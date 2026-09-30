# GoのLinter警告「const is unused」が発生する仕組みと解消法

## 発生した警告

`domain/model/user/user_entity.go` にて以下の Linter 警告（`unused`）が発生した。

```text
[WARNING] L13:7 - const maxUnCompletedTaskCount is unused (unused) (go-golangci-lint-v2)
[WARNING] L22:23 - func (*UserEntity).canCreateTask is unused (unused) (go-golangci-lint-v2)
```

対象コード:

```go
const maxUnCompletedTaskCount = 20

// ...

func (ue *UserEntity) canCreateTask(uncompletedTaskCount int) bool {
	canCreate := uncompletedTaskCount < maxUnCompletedTaskCount
	return canCreate
}
```

---

## なぜ使っているはずの定数が「unused」になるのか？

コード上では `canCreateTask` の中で `maxUnCompletedTaskCount` を参照している。  
それにもかかわらず未使用警告が出る原因は、**参照元である `canCreateTask` 自体がデッドコード（到達不能コード）と判定されているため**。

### 判定のメカニズム

1. **`canCreateTask` が非公開（private）かつ未呼び出し**
   - Go では先頭小文字の識別子は非公開（package private）となる。
   - `user` パッケージ内のどこからも `canCreateTask` が呼び出されていないため、Linter は「この関数自体が誰にも使われていない不要コード」と判定する。
2. **デッドコード内でしか使われていない定数も「未使用」になる**
   - `maxUnCompletedTaskCount` も先頭小文字（非公開）。
   - この定数を参照しているのは、すでにデッドコードと判定された `canCreateTask` のみ。
   - Linter は「外部の有効な処理から到達する経路が存在しない」と判断し、連鎖して定数側にも `unused` 警告を出す。

---

## 一次ソース

Go の静的解析ツール `staticcheck`（内部チェッカー: `unused` / ルール名: `U1000`）の公式仕様より:

- 一次ソース: [dominikh/go-tools - unused/unused.go](https://github.com/dominikh/go-tools/blob/master/unused/unused.go)
- 仕様:
  - 公開（Export）された識別子、または `main` / `init` 関数をコールグラフ探索の起点（ルートノード）とする。
  - 起点から到達できない非公開の関数・メソッドは未使用（Unused）と判定する。
  - 到達不能な関数からのみ参照されている非公開定数・変数も、探索グラフ上で到達不能（Unused）となる。

---

## TS / PHP との比較

- **TypeScript (ESLint: `@typescript-eslint/no-unused-vars`)**:
  クラス内の private メソッドがどこからも呼ばれていない場合、およびその中でのみ参照されている定数・変数が同様に unused 警告の対象になる。
- **PHP (PHPStan / Psalm)**:
  `Unused private method` や `Unused private constant` として静的解析ツールに検知される挙動と同じ。

---

## 解消パターン

1. **外部（usecase など）から呼ぶ場合**:
   メソッド先頭を大文字にして公開（Export）する。
   ```go
   func (ue *UserEntity) CanCreateTask(uncompletedTaskCount int) bool {
       return uncompletedTaskCount < maxUnCompletedTaskCount
   }
   ```
2. **定数を外部から参照する場合**:
   定数先頭を大文字にして公開する。
   ```go
   const MaxUnCompletedTaskCount = 20
   ```
3. **Entity 内部の公開メソッドから呼ぶ場合**:
   `NewUserEntity` や将来の公開メソッド内から `ue.canCreateTask(...)` を呼び出す実装になれば、到達可能パスがつながり警告は自然消滅する。
