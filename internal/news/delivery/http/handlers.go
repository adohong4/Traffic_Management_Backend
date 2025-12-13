package http

import (
	"net/http"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/internal/news"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type newsHandlers struct {
	cfg    *config.Config
	newsUC news.UseCase
	logger logger.Logger
}

func NewsHandlers(cfg *config.Config, newsUC news.UseCase, log logger.Logger) news.Handlers {
	return &newsHandlers{cfg: cfg, newsUC: newsUC, logger: log}
}

// @Summary Create news
// @Description Create a new news article
// @Tags News
// @Accept json
// @Produce json
// @Param news body models.News true "News object"
// @Success 201 {object} models.News
// @Failure 400,401,500 {object} httpErrors.RestError
// @Security JWT
// @Router /news [post]
func (h *newsHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		n := &models.News{}
		if err := c.Bind(n); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		created, err := h.newsUC.Create(c.Request().Context(), n)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusCreated, created)
	}
}

// @Summary Update news by ID
// @Description Update an existing news
// @Tags News
// @Accept json
// @Produce json
// @Param id path string true "News ID"
// @Param news body models.News true "News object"
// @Success 200 {object} models.News
// @Failure 400,401,500 {object} httpErrors.RestError
// @Security JWT
// @Router /news/{id} [put]
func (h *newsHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.News{}
		if err = c.Bind(n); err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = id

		updated, err := h.newsUC.Update(c.Request().Context(), n)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, updated)
	}
}

// @Summary Delete news by ID
// @Description Soft delete news
// @Tags News
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} models.News
// @Failure 401,500 {object} httpErrors.RestError
// @Security JWT
// @Router /news/{id} [delete]
func (h *newsHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		UUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		n := &models.News{}
		if err = c.Bind(n); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		n.Id = UUID

		deletedNews, err := h.newsUC.Delete(ctx, n)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, deletedNews)
	}
}

// @Summary Get news by ID
// @Description Get single news (view count will increase)
// @Tags News
// @Produce json
// @Param id path string true "News ID"
// @Success 200 {object} models.News
// @Failure 404,500 {object} httpErrors.RestError
// @Router /news/{id} [get]
func (h *newsHandlers) FindById() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		news, err := h.newsUC.FindById(c.Request().Context(), id)
		if err != nil {
			return c.JSON(httpErrors.ErrorResponse(err))
		}
		return c.JSON(http.StatusOK, news)
	}
}

// @Summary Get all news (paginated)
// @Description Get list of news with pagination
// @Tags News
// @Produce json
// @Param page query int false "Page" default(1)
// @Param size query int false "Size" default(10)
// @Success 200 {object} models.NewsList
// @Router /news [get]
func (h *newsHandlers) FindAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		list, err := h.newsUC.FindAll(ctx, pq)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, list)
	}
}
