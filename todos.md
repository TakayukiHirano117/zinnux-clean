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

Aggregate実装
