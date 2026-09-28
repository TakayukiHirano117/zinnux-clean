# Entity保存と createdAt

**最終更新**: 2026-09-28

## 要点

- Repository の入口は Entity。`created_at` は今の `UserEntity` には持たない
- 永続化の時刻は Interface Adapter（repository 実装）が `model.User` に詰め替えるときに付ける
- GORM は `CreatedAt` がゼロ値なら `Create` 時に現在時刻を入れる（非ゼロなら上書きしない）
- 業務ルールとして「いつ作られたか」が必要なら Entity の生成時に持たせ、repository はその値をコピーする

## 補足

- 根拠: Uncle Bob *The Clean Architecture*（Entity は業務ルール、DB 形式への変換は Interface Adapters）
- 根拠: https://gorm.io/docs/conventions.html の CreatedAt / UpdatedAt
- 対象コード: `domain/model/user/user_entity.go`, `model/user.go`, `repository/user_repository.go`

## 元の文脈

- 質問: Entity 単位で repository 保存するとき createdAt はどうするか
