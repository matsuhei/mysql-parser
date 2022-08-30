## mysql-parser

mysqlのtableをparseして、各種formatに置き換えることができます。
まだオプション系の対応はしていないので、これから対応予定です。
`.mysql-parser_template.yaml` に必要なenv情報が記載されているので、参考にしてください。

## 対応format
- [figjam](https://www.figma.com/community/widget/1102833776087940938/Database-Table)
- [WIP] go struct


## 実行方法
```shell
go run main.go mysql-parser ${table_name}
```
