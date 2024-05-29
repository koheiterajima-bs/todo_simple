package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

// 構造体の宣言(mapではそれぞれの型が同じでなくてはならないため)
// 構造体のフィールドとjsonのキーを紐づけ
type TodoList struct {
	ID          uuid.UUID `json:"ID"`
	Name        string    `json:"Name"`
	TodoContent string    `json:"TodoContent"`
}

// 【jsonファイル読み込み】
func readJson() []TodoList {
	file, err := os.ReadFile("./todoList.json") // os.ReadFileは指定されたファイルの内容を全て読み込み、バイトスライスとして返す。
	if err != nil {
		fmt.Println("jsonファイルが読み込めませんでした。")
	}
	jsonData := make([]TodoList, 0)       // jsonデータを格納するスライスの変数を宣言
	err = json.Unmarshal(file, &jsonData) // json.Unmarshalでjsonデータを解析し、Goのデータ構造に変換
	if err != nil {
		fmt.Println("jsonファイルデータを配列に格納できませんでした。")
	}

	return jsonData // addSliceにスライスを渡す
}

// 【追加されたスライスをjsonへ格納する】
func writeJson(jsonData []TodoList) {
	output, err := json.MarshalIndent(jsonData, "", "\t\t") // Goのデータ構造をjson形式の文字列に変換
	if err != nil {
		fmt.Println("jsonファイルへ書き込みできませんでした。")
	}

	// ファイルに書きこみ
	err = os.WriteFile("./todolist.json", output, 0600)
	if err != nil {
		fmt.Println("ファイルへの書き込みに失敗しました:", err)
	}
}

// 【新しいTodoをスライスに追加】
func addSlice(jsonData []TodoList) []TodoList {
	fmt.Println("Todoを作成します。")

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
	todo := TodoList{
		ID:          uuid,
		Name:        name,
		TodoContent: todoContent,
	}

	// 既存のjsonDataに新たに追加
	jsonData = append(jsonData, todo)

	fmt.Println("Todo追加が完了しました。")

	return jsonData // writeJsonへスライスを渡す
}

// 【スライスの削除】
func deleteSlice(jsonData []TodoList) []TodoList {
	fmt.Println("どのTodoを削除しますか？\n削除したいTodoのIDを貼り付けてください。")

	// jsonDataの配列一覧を出す
	for _, obj := range jsonData {
		fmt.Println(obj)
	}

	// scanにて一旦入力させるためにstring型に設定(uuid.UUID型は認識できない)
	var deleteIDStr string
	fmt.Scan(&deleteIDStr)

	// uuid.Parse()は文字列型をuuid.UUID型に変えるための関数
	deleteID, err := uuid.Parse(deleteIDStr)
	if err != nil {
		fmt.Println("無効なIDです。")
	}

	// 一致するIDがあるかどうかを判定する変数、以下のfor文にて一致すればtrueに変更
	judgeDeleteID := false

	for i, obj := range jsonData {
		if obj.ID == deleteID {
			jsonData = append(jsonData[:i], jsonData[i+1:]...) // i番目のオブジェクトを削除
			judgeDeleteID = true
			break
		}
	}

	if judgeDeleteID {
		fmt.Println("Todoを削除しました。")
	} else {
		fmt.Println("一致するIDはありませんでした。")
	}

	// 【オブジェクトの番号を指定して削除するver】
	// var deleteNumber int
	// fmt.Scan(&deleteNumber)
	// i := deleteNumber - 1
	// jsonData = append(jsonData[:i], jsonData[i+1:]...)
	// fmt.Println("i番目のTodo削除が完了しました。")

	// 【スライスを空にして、全削除ver】
	// jsonData = jsonData[:0]
	// fmt.Println(jsonData)

	return jsonData
}

// 【スライスの更新】
func updateSlice(jsonData []TodoList) []TodoList {
	// jsonDataの配列一覧を出す
	for _, obj := range jsonData {
		fmt.Println(obj)
	}

	fmt.Println("どのTodoを更新しますか？\n更新したいTodoのIDを貼り付けてください。")
	// uuid.UUID型だと、fmt.Scan()が型に対応していないため
	var updateIDStr string
	fmt.Scan(&updateIDStr)

	// uuid.Parse()にてuuid.UUID型に変換
	updateID, err := uuid.Parse(updateIDStr)
	if err != nil {
		fmt.Println("無効なIDです。")
	}

	for i, obj := range jsonData {
		if obj.ID == updateID {
			fmt.Println("IDが確認できました。")

			fmt.Println("Todoの内容を記載してください。")
			var updateTodoContent string
			fmt.Scan(&updateTodoContent)

			jsonData[i].TodoContent = updateTodoContent // TodoContentの更新
			fmt.Println("Todoの内容を更新しました。")
		}
		fmt.Println("一致するIDがありませんでした。")
		break // 上記のPrintlnを一致しなかったID分だけ繰り返したくないが、[the surrounding loop is unconditionally terminated (SA4004)]とエラーが出る。
	}

	// 【オブジェクトの番号を指定して更新するver】
	// fmt.Println("どのTodoを更新しますか？\n更新したいTodoの番号を入力してください。")
	// var updateNumber int
	// fmt.Scan(&updateNumber)
	// fmt.Println("Todoの内容を記載してください。")
	// var updateTodoContent string
	// fmt.Scan(&updateTodoContent)
	// i := updateNumber - 1
	// jsonData[i].TodoContent = updateTodoContent
	// fmt.Println("Todo更新が完了しました。")

	return jsonData
}

func main() {
	// コマンドラインからの入力を受け付ける
	var command string
	fmt.Print("何をしますか？(Todo作成/Todo更新/Todo削除)")
	// 参照渡しでメモリアドレスを受け取る
	// fmt.Scan()は、引数として値を受け取るが、その引数には変数のアドレスが必要
	fmt.Scan(&command)

	switch command {
	case "Todo作成":
		jsonData := readJson()
		sliceAdd := addSlice(jsonData)
		writeJson(sliceAdd)
	case "Todo更新":
		jsonData := readJson()
		sliceUpdate := updateSlice(jsonData)
		writeJson(sliceUpdate)
	case "Todo削除":
		jsonData := readJson()
		sliceDelete := deleteSlice(jsonData)
		writeJson(sliceDelete)
	default:
		fmt.Println("エラー発生しました、適切なコマンド入力をしてください。")
	}
}
