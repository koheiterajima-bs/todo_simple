## SQL
- SQLはpostgreSQLを使う
```sh
-- psqlにログイン
psql -U postgres

-- パスワードを入力
Password for user postgres:

-- データベースの作成
postgres=# CREATE DATABASE todos;

-- データベースの一覧を確認
postgres=# \l

-- 削除するとき
postgres=# DROP DATABASE データベース名

-- テーブルの作成
postgres=# CREATE TABLE todos
postgres-# (
postgres(# id serial PRIMARY KEY,
postgres(# name text NOT NULL,
postgres(# todo text NOT NULL
postgres(# );

-- テーブルにtodoを登録
postgres=# INSERT INTO todo_list (id, name, todo) VALUES (1, 'terajima', 'todoリスト作成');

-- テーブルの確認
postgres=# select * from todo_list;

-- テーブルの削除
postgres=# DROP TABLE IF EXISTS todo_list;

-- postgresを閉じる
postgres=# \q

```
**参考サイト**
[誰でも分かる！PostgreSQLでDB構築！](https://qiita.com/hiroyuki_mrp/items/10322eeb29bb8e35987f)


## 5/20
日本語がうまく認識されない…
おそらく以下のエラーが原因…
```sh
ERROR:  新しい照合順序(ja_JP.UTF-8)はテンプレートデータベースの照合順序(Japanese_Japan.utf8)と互換性がありません
```
## 5/21
**todo作成参考サイト**
[【Go言語】DB接続サンプル（PostgreSQL）](https://lelelemon.hatenablog.com/entry/2024/01/15/075032)
とりあえず、DBに接続して、入力画面まで進んだ

## 5/22
データベースの作成まではできたが、データベースへの追加ができない
```sh
$ go run main.go
DB接続ができました。
何をしますか？ (Todo作成/Todo更新/Todo削除): Todo作成
Todo作成するよ
名前を入力してください。hoge
Todoの内容を入力してください。hogehoge
2024/05/22 22:51:50 pq: ユーザー"user"のパスワード認証に失敗しました
exit status 1
```

そもそも
```sh
$ psql
Password for user koheiterajimabs:
psql: error: connection to server at "localhost" (::1), port 5432 failed: FATAL:  ユーザー"koheiterajimabs"のパスワード認証に失敗しました
<!-- では入れず、 -->
psql -U postgres

<!-- では入れるのもよくわからない、、 -->
```

## 5/24
以下のコードを参考に構造体、メソッド、インターフェースあたりを作成したいかも
```go
package main

import "fmt"

// 構造体の宣言
type Blog struct {
	title   string
	content string
}
type Blog2 struct {
	title     string
	paragraph []string
}

// メソッドの宣言
func (b Blog) GetFullArticle() string {
	return b.title + "\n" + "------------" + "\n" + b.content
}

func (b Blog2) GetFullArticle() string {
	article := b.title + "\n" + "------------" + "\n"

	for _, paragraph := range b.paragraph {
		article += paragraph + "\n\n"
	}

	return article
}

// 各メソッドを表示させる関数
func displayBlog(b Blog) {
	fmt.Println(b.GetFullArticle())
}

func displayBlog2(b Blog2) {
	fmt.Println(b.GetFullArticle())
}

func main() {
	// 構造体をインスタンス化し、変数に代入している
	blog := Blog{"titleですよ", "contentですよ"}
	displayBlog(blog)

	blog2 := Blog2{"titleですよ", []string{"a", "b"}}
	displayBlog2(blog2)
}

```

## 5/25
SQLからjson形式に変更し、jsonファイルへの書き込みができるようになった。
しかし、現状上書きになってしまう。下に更新される形式にしたい。
また、mapではすべてのキーと値はそれぞれ同じ値でなくてはならないことを初めて知った。
```go
todo := map[string]string {
	"id": uuid,
	"name": name,
	"todoContent": todoContent,
}

//これはできない(idをint型にしたい)
todo := map[string]string|int {
	"id": uuid,
	"name": name,
	"todoContent": todoContent,
}
```
この場合、構造体を用いることで、異なる型を共存できる
```go
type todoList struct {
	ID uuid.UUID
	Name string
	TodoContent string
}
```

**参考サイト**
(GoのJSON操作【プログラミング初心者向け教材】)[https://tokitsubaki.com/go-json-manipulation/411/#toc1]

for rangeを使って既存のデータに追加した
→これをjson形式に変換し、jsonファイルに戻さなくては

## 5/26
- jsonファイルの読み込み、書き込みをmapとスライスのどちらで行うか悩んだが、スライスの方がコードの可読性が高いと感じた。
- スライスの場合
```go
// 【jsonファイル読み込み】
func readJson() []todoList {
	file, err := os.ReadFile("./todoList.json") // os.ReadFileは指定されたファイルの内容を全て読み込み、バイトスライスとして返す。
	if err != nil {
		fmt.Println("jsonファイルが読み込めませんでした。")
	}
	jsonData := make([]todoList, 0)       // jsonデータを格納するスライスの変数を宣言
	err = json.Unmarshal(file, &jsonData) // json.Unmarshalでjsonデータを解析し、Goのデータ構造に変換
	if err != nil {
		fmt.Println("jsonファイルデータを配列に格納できませんでした。")
	}

	return jsonData // addSliceにスライスを渡す
}

// 【新しいTodoをスライスに追加】
func addSlice(jsonData []todoList) []todoList {
	var name string
	var todoContent string

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	fmt.Print("名前を入力してください。\t")
	fmt.Scan(&name)

	fmt.Print("Todoを入力してください。\t")
	fmt.Scan(&todoContent)

	// 構造体をインスタンス化し、変数に代入
	todo := todoList{
		ID:          uuid,
		Name:        name,
		TodoContent: todoContent,
	}

	// 既存のjsonDataに新たに追加
	jsonData = append(jsonData, todo)

	fmt.Println(jsonData)

	return jsonData // writeJsonへスライスを渡す
}
```

- マップの場合
```go
func readJson() []map[string]interface{} {
	file, err := os.ReadFile("./todoList.json") // os.ReadFileは指定されたファイルの内容を全て読み込み、バイトスライスとして返す。
	if err != nil {
		fmt.Println("jsonファイルが読み込めませんでした。")
	}
	jsonData := []map[string]interface{}{} // jsonデータを格納するマップの変数宣言
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		fmt.Println("jsonファイルデータをマップに格納できませんでした。")
	}

	return jsonData
}

// 【新しいTodoをスライスに追加】
func addSlice(jsonData []map[string]interface{}) map[string]interface{} {
	var name string
	var todoContent string

	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	fmt.Print("名前を入力してください。\t")
	fmt.Scan(&name)

	fmt.Print("Todoを入力してください。\t")
	fmt.Scan(&todoContent)

	// 新しいオブジェクトの作成
	newObject := map[string]interface{}{
		"ID":          uuid,
		"Name":        name,
		"TodoContent": todoContent,
	}

	// 既存のjsonDataに新しいオブジェクトを追加
	jsonData := append(jsonData, newObject)

	fmt.Println(jsonData)

	return jsonData // writeJsonへスライスを渡す
}

		jsonData := readJson() // 消す
		addSlice(jsonData)     // 消す
```
addSliceのmap版での書き方がわからん、、

## 5/27
todoの一括削除まではできるようになった。
が特定のidのオブジェクトの消し方がわからない、、
スライスを削除する際のappend関数がわかっていない。
**参考サイト**
(Go 言語でスライスから要素を消すには)[https://zenn.dev/mattn/articles/31dfed3c89956d]
```go
package main

import "fmt"

func main() {
	a := []int{1, 2, 3, 4, 5}

	i := 2
	// a[:2]は0番目から2-1番目の要素を取得、a[3:]は3番目から最後の要素まで取得
	// ...はスライスの展開演算子で、スライスを展開して要素を1つずつ追加する
	// a = append(a[:i], a[i+1:]...)
	a = append(a[:i], a[i+2:]...)
	fmt.Println(a)
}
```
特定の番号のものは削除できたが、IDを指定はできなかった。
また、存在しない番号を入力するとパニックを起こしてしまう。
**参考サイト**
(golang スライスとは 使い方、注意点について解説)[https://note.com/webdrawer/n/n3c0afb015e28]
更新機能も付ける

## 5/28
### 構造体1文字目を大文字にする

### updateSliceについて、IDが一致しなければfor文を即座に終了させたいが、エラーが出る。
```sh
the surrounding loop is unconditionally terminated (SA4004)
```

### breakとreturnの違いがわかった
#### break
- ループやswitch、if文の中で使用でき、それらから直ちに脱出する。
- その関数からは脱出しない。
#### return
- 関数の実行を即座に終了する。
- 関数内の任意の場所にて使用できる。