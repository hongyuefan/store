package initstore

import (
	"database/sql"

	"store/core"

	_ "github.com/go-sql-driver/mysql"
)

const (
	SelectCode   = "select a.class_id,a.code ,a.name,a.score, b.class from store_code a,store_class b where a.class_id=b.id"
	SqlSaveScore = "update store_code  set score = ? where code = ?"
)

var DB *sql.DB

func OnInitCore(core *core.Core) (err error) {

	DB, err = open(core.Cont.DBUrl)

	if err != nil {
		return
	}

	if err = selectCode(DB, core.CallbackInitMap); err != nil {
		return
	}
	return
}

func open(url string) (DB *sql.DB, err error) {
	return sql.Open("mysql", url)

}

func SaveScore(code string, score float64) error {
	return saveScore(DB, code, score)
}

func saveScore(db *sql.DB, code string, score float64) error {

	stmt, err := db.Prepare(SqlSaveScore)

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(score, code)

	return err
}

type callback func(int, string, float64, string, string) error

func selectCode(db *sql.DB, f callback) error {

	var (
		classId int
		code    string
		name    string
		class   string
		score   float64
	)

	stmt, err := db.Prepare(SelectCode)

	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.Scan(&classId, &code, &name, &score, &class); err != nil {
			return err
		}
		if err := f(classId, code, score, name, class); err != nil {
			return err
		}
	}

	return nil

}
