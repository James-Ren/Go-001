package service

import (
	"context"

	"github.com/James-Ren/Go-000/dao"
	"github.com/James-Ren/Go-000/model"
)

func GetArticle(ctx context.Context, id int) (*model.Article, error) {
	return dao.GetArticle(ctx, id)
}
