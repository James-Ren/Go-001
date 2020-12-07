package dao

import (
	"context"
	"database/sql"

	"github.com/James-Ren/Go-000/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var db *sql.DB

// ErrRecordNotFound indicates cannot find record
var ErrRecordNotFound = errors.New("Not Found")

func init() {
	//Init db
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/db?charset=utf8&timeout=1s&readTimeout=1s")
	if err != nil {
		// 假设有强依赖，发生错误，直接panic
		panic(err)
	}

}

func GetArticle(ctx context.Context, id int) (*model.Article, error) {
	article := &model.Article{}
	row := db.QueryRowContext(ctx, "select id,title,content from articles where id=?", id)
	err := row.Scan(&article.ID, &article.Title, &article.Content)
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(ErrRecordNotFound, "No corresponding article")
	}
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get article")
	}
	return article, nil
}
