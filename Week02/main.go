package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/James-Ren/Go-000/dao"
	"github.com/James-Ren/Go-000/service"
	"github.com/pkg/errors"
)

func main() {
	//模拟 api层（controller）
	articleID := 30
	article, err := service.GetArticle(context.Background(), articleID)
	if errors.Is(err, dao.ErrRecordNotFound) {
		// 返回 404 并打印日志
		log.Printf("404: %v", err)
		return
	}
	if err != nil {
		// 返回 500 并打印日志
		log.Printf("500: %+v", err)
		return
	}
	//返回200 并打印数据
	resp, err := json.Marshal(article)
	if err != nil {
		// 返回 500 并打印日志
		log.Printf("500: %v", err)
		return
	}
	fmt.Println(string(resp))
}
