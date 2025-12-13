package repository

import (
	"context"

	"github.com/adohong4/driving-license/internal/models"
	news "github.com/adohong4/driving-license/internal/news"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type newsRepo struct {
	db *sqlx.DB
}

func NewNewsRepo(db *sqlx.DB) news.Repository {
	return &newsRepo{db: db}
}

func (r *newsRepo) Create(ctx context.Context, n *models.News) (*models.News, error) {
	news := &models.News{}
	err := r.db.QueryRowxContext(ctx, createNewsQuery,
		n.Id, n.Code, n.Image, n.Title, n.Content, n.Category, n.Author, n.Type, n.Target,
		n.TargetUser, n.Tag, n.View, n.Status, n.Version, n.CreatorId, n.ModifierID,
		n.CreatedAt, n.UpdatedAt, n.Active,
	).StructScan(news)
	if err != nil {
		return nil, errors.Wrap(err, "newsRepo.Create")
	}
	return news, nil
}

func (r *newsRepo) Update(ctx context.Context, n *models.News) (*models.News, error) {
	news := &models.News{}
	err := r.db.QueryRowxContext(ctx, updateNewsQuery,
		n.Code, n.Image, n.Title, n.Content, n.Category, n.Author, n.Type, n.Target,
		n.TargetUser, n.Tag, n.Status, n.ModifierID, n.UpdatedAt, n.Id,
	).StructScan(news)
	if err != nil {
		return nil, errors.Wrap(err, "newsRepo.Update")
	}
	return news, nil
}

func (r *newsRepo) DeleteNews(ctx context.Context, db *models.News) (*models.News, error) {
	news := &models.News{}
	if err := r.db.QueryRowxContext(ctx, deleteNewsQuery,
		db.ModifierID, db.UpdatedAt, db.Id,
	).StructScan(news); err != nil {
		return nil, errors.Wrap(err, "NewsRepo.DeleteNews.StructScan")
	}
	return news, nil
}

func (r *newsRepo) FindById(ctx context.Context, id uuid.UUID) (*models.News, error) {
	news := &models.News{}
	err := r.db.GetContext(ctx, news, getNewsByIdQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "newsRepo.FindById")
	}

	_, _ = r.IncrementView(ctx, id)

	return news, nil
}

func (r *newsRepo) FindAll(ctx context.Context, pq *utils.PaginationQuery) (*models.NewsList, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, getTotalCountQuery); err != nil {
		return nil, errors.Wrap(err, "newsRepo.FindAll.total")
	}

	var items []*models.News
	err := r.db.SelectContext(ctx, &items, getAllNewsQuery, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "newsRepo.FindAll.select")
	}

	return &models.NewsList{
		TotalCount: total,
		TotalPages: utils.GetTotalPage(total, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), total, pq.GetSize()),
		News:       items,
	}, nil
}

func (r *newsRepo) IncrementView(ctx context.Context, id uuid.UUID) (int, error) {
	var view int
	err := r.db.GetContext(ctx, &view, incrementViewQuery, id)
	return view, err
}
