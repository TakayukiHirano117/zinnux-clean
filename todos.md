Aggregate実装
学習用にしてたから.env入れてたが消したい

## projectユースケース作成

backlog的な何かにしようこれ

- [ ] project一覧
  - [ ] 課題やユーザーがprojectに所属するイメージ。なのでprojectがユーザーIDとか持つとかはしない
  - [ ] データモデル作ってマイグレーションする
  - [ ] controller
    - [ ] getAllProjects
      - [ ] 所属している人が見れるかどうか認可の処理必要
  - [ ] usecase
    - [ ] getAllProjectsユースケース
    - [ ] ドメイン・インフラ使って組み立てるだけ
  - [ ] infra
    - [ ] repository作成する
  - [ ] domain
    - [ ] entity
    - [ ] vo
- [ ] project詳細
- [ ] project作成
- [ ] project削除
- [ ] project編集

## バリデーション

## iota or constの定数で状態管理

## エラーハンドリング

- [ ] バリデーションエラーの取得と返し方
- [ ] ドメインエラーの実装と返し方

## ロガー実装

## gRPC導入

## Docker環境構築

https://docs.docker.com/guides/golang/ を参考に作る

## マイグレートツール導入

https://gorm.io/ja_JP/docs/migration.html
gormだと弱い

## makefile修正

GO_ENVべた書きしているので動的に変えたい。
できればdev以外はy/Nでチェック走るようにしたい。
関係ないがCI/CDでも同じようなことしたい。

## Done

HTTPバリデーション

- [x] go-playground/validatorを使う
- [x] main.goでe.Validator = &CustomValidator{validator: validator.New()}を登録する
- [x] controllerでc.Validate(&signUpRequest);でバリデーションかける

##

## ユーザー登録API改修

DTOはユースケースに置く

- [x] RequestDTO作成
- [x] ResponseDTO作成
- [x] Entity作成
  - [x] Entity単位でRepositoryで保存するように修正
    - [x] repositoryでuser_entiry受け取る
    - [x] データモデルにつめかえて保存する
- [x] データモデルはそのまま
- [x] Controller, UseCaseでRequestの型とDTOの型を分ける
  - [x] Requestの型がDTOに依存するので。そうなるとUseCaseのDTOのプロパティの数が増減するとそれにcontrollerのrequestのプロパティ数も引っ張られる。密結合になっちゃう

vo実装

- [x] user email
- [x] uuid
  - [x] コンストラクタ
    - [x] DBから作るやつ
    - [x] 新規作成するやつ
