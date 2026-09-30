## ユーザー登録API改修
DTOはユースケースに置く
- [x] RequestDTO作成
- [x] ResponseDTO作成
- [x] Entity作成
  - [ ] Entity単位でRepositoryで保存するように修正
    - [ ] repositoryでuser_entiry受け取る
    - [ ] データモデルにつめかえて保存する
- [ ] データモデルはそのまま
- [ ] Controller, UseCaseでRequestの型とDTOの型を分ける
  - [ ] Requestの型がDTOに依存するので。そうなるとUseCaseのDTOのプロパティの数が増減するとそれにcontrollerのrequestのプロパティ数も引っ張られる。密結合になっちゃう

vo実装
- [ ] user email
- [ ] uuid
  - [ ] コンストラクタ
    - [ ] DBから作るやつ
    - [ ] 新規作成するやつ

Aggregate実装


## バリデーション
HTTPバリデーション
- [ ] go-playground/validatorを使う
- [ ] main.goでe.Validator = &CustomValidator{validator: validator.New()}を登録する
- [ ] controllerでc.Validate(&signUpRequest);でバリデーションかける

## エラーハンドリング
- [ ] バリデーションエラーの取得と返し方
- [ ] ドメインエラーの実装と返し方

## ロガー実装
