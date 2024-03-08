package news

import (
	"classting/domain"
	"classting/pkg/cerrors"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type newsRepository struct {
	sqlDB *sql.DB
}

func NewNewsRepository(sqlDB *sql.DB) *newsRepository {
	return &newsRepository{
		sqlDB: sqlDB,
	}
}

var _ domain.NewsRepository = (*newsRepository)(nil)

func (n newsRepository) CreateNews(ctx context.Context, news domain.News) (int, error) {
	const op cerrors.Op = "news/schoolRepository/CreateNews"

	result, err := n.sqlDB.ExecContext(ctx, createNewsQuery, news.SchoolID, news.UserID, news.Title)
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return int(userID), nil
}

func (n newsRepository) ListNews(ctx context.Context, params domain.ListNewsParams) ([]domain.News, error) {
	const op cerrors.Op = "news/newsRepository/ListNews"

	var news []domain.News

	query := fmt.Sprintf(listNewsQuery, params.AfterCursor())

	rows, err := n.sqlDB.QueryContext(ctx, query, params.UserID)
	if err != nil {
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.News
		err := rows.Scan(
			&product.ID,
			&product.CreateDate,
			&product.UpdateDate,
			&product.DeleteDate,
			&product.SchoolID,
			&product.UserID,
			&product.Title,
		)
		if err != nil {
			return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
		}
		news = append(news, product)
	}

	return news, nil
}

func (n newsRepository) UpdateNews(ctx context.Context, news domain.News) error {
	const op cerrors.Op = "news/newsRepository/UpdateNews"

	updateNewsQuery := `UPDATE news SET title = ? WHERE id = ?`
	_, err := n.sqlDB.ExecContext(ctx, updateNewsQuery, news.Title, news.ID)
	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return nil
}

func (n newsRepository) FindNewsByID(ctx context.Context, newsID int) (*domain.News, error) {
	const op cerrors.Op = "news/newsRepository/FindNewsByID"

	var news domain.News

	err := n.sqlDB.QueryRowContext(ctx, findNewsByIDQuery, newsID).Scan(
		&news.ID,
		&news.CreateDate,
		&news.UpdateDate,
		&news.DeleteDate,
		&news.SchoolID,
		&news.UserID,
		&news.Title,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return &news, nil
}

func (n newsRepository) DeleteNews(ctx context.Context, newsID int) error {
	const op cerrors.Op = "news/newsRepository/DeleteNews"

	_, err := n.sqlDB.ExecContext(ctx, deleteNewsQuery, time.Now().UTC(), newsID)
	if err != nil {
		return cerrors.E(op, cerrors.Internal, err, "서버 에러가 발생했습니다.")
	}

	return nil
}
