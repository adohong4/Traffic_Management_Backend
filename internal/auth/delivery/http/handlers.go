package http

import (
	"net/http"
	"strconv"

	"github.com/adohong4/driving-license/config"
	"github.com/adohong4/driving-license/internal/auth"
	"github.com/adohong4/driving-license/internal/models"
	"github.com/adohong4/driving-license/pkg/httpErrors"
	"github.com/adohong4/driving-license/pkg/logger"
	"github.com/adohong4/driving-license/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Auth handlers
type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
	logger logger.Logger
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase, log logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC, logger: log}
}

// CreateUser godoc
// @Summary      Create a new user (admin only)
// @Description  Creates a new user account. Typically used by administrators.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User information"
// @Success      200   {object}  models.UserWithToken
// @Failure      400   {object}  httpErrors.RestError  "Invalid input"
// @Failure      500   {object}  httpErrors.RestError  "Server error"
// @Router       /auth/create [post]
// @Security     BearerAuth
func (h *authHandlers) CreateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &models.User{}
		// Chỉ bind dữ liệu từ request
		if err := c.Bind(user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Chuẩn bị user (sinh Id, HashPassword, v.v.)
		if err := user.PrepareCreate(); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Validate sau khi PrepareCreate
		if err := utils.ValidateStruct(c.Request().Context(), user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		userWithToken, err := h.authUC.CreateUser(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		h.logger.Infof("User created successfully, ID: %s", userWithToken.User.Id)
		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with identity number and password
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      models.User  true  "Login credentials"
// @Success      200      {object}  models.UserWithToken
// @Failure      400      {object}  httpErrors.RestError  "Invalid request"
// @Failure      401      {object}  httpErrors.RestError  "Invalid credentials"
// @Failure      500      {object}  httpErrors.RestError
// @Router       /auth/login [post]
func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		IdentityNO string `json:"identity_no" db:"identity_no" validate:"required,lte=20"`
		Password   string `json:"password,omitempty" db:"password" validate:"omitempty,gte=6"`
	}
	return func(c echo.Context) error {
		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		ctx := c.Request().Context()
		userWithToken, err := h.authUC.Login(ctx, &models.User{
			IdentityNo: login.IdentityNO,
			Password:   login.Password,
		})
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, userWithToken)
	}
}

// Logout godoc
// @Summary      Logout user
// @Description  Invalidate the current session and remove session cookie
// @Tags         Auth
// @Accept 		 json
// @Produce      json
// @Success      200  {string}  string  "ok"
// @Failure      401  {object}  httpErrors.RestError
// @Router       /auth/logout [post]
// @Security     BearerAuth
func (h *authHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		utils.DeleteSessionCookie(c, h.cfg.Session.Name)

		return c.NoContent(http.StatusOK)
	}
}

// Update godoc
// @Summary      Update user information
// @Description  Update an existing user by ID (admin or own account)
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        id    path     string       true  "User ID (UUID)"
// @Param        user  body     models.User  true  "Updated user data"
// @Success      200   {object}  models.User
// @Failure      400   {object}  httpErrors.RestError
// @Failure      403   {object}  httpErrors.RestError  "Forbidden"
// @Failure      404   {object}  httpErrors.RestError  "User not found"
// @Failure      500   {object}  httpErrors.RestError
// @Router       /auth/{id} [put]
// @Security     BearerAuth
func (h *authHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		Id, err := uuid.Parse(c.Param("id"))

		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user := &models.User{}
		user.Id = Id

		if err = utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		updatedUser, err := h.authUC.Update(ctx, user)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Retrieve a user profile by user ID
// @Tags         Auth
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  models.User
// @Failure      400  {object}  httpErrors.RestError  "Invalid ID"
// @Failure      404  {object}  httpErrors.RestError  "User not found"
// @Failure      500  {object}  httpErrors.RestError
// @Router       /auth/{id} [get]
// @Security     BearerAuth
func (h *authHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		Id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		user, err := h.authUC.GetByID(ctx, Id)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, user)
	}
}

// Delete godoc
// @Summary      Delete user account
// @Description  Permanently delete a user account (admin or self-deletion)
// @Tags         Auth
// @Produce      json
// @Param        id       path      string  true   "User ID (UUID)"
// @Param        version  query     int     false  "Optimistic lock version (optional)"
// @Success      200      {object}  object{message=string}  "User deleted successfully"
// @Failure      400      {object}  httpErrors.RestError
// @Failure      403      {object}  httpErrors.RestError  "Forbidden"
// @Failure      409      {object}  httpErrors.RestError  "Conflict (version mismatch)"
// @Failure      500      {object}  httpErrors.RestError
// @Router       /auth/{id} [delete]
// @Security     BearerAuth
func (h *authHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		// Parse user ID từ tham số đường dẫn
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError(errors.Wrap(err, "invalid user ID")))
		}

		// Lấy modifierId từ context (giả định user hiện tại được lưu trong context)
		currentUser, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		// Giả định version được gửi qua query parameter hoặc body, nếu không có thì mặc định là 0
		versionStr := c.QueryParam("version")
		var version int
		if versionStr != "" {
			v, err := strconv.Atoi(versionStr)
			if err != nil {
				utils.LogResponseError(c, h.logger, err)
				return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError(errors.Wrap(err, "invalid version")))
			}
			version = v
		}

		// Gọi authUC.Delete để xóa người dùng
		if err := h.authUC.Delete(ctx, id, currentUser.Id, version); err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		// Trả về phản hồi thành công
		return c.JSON(http.StatusOK, map[string]string{"message": "user deleted successfully"})
	}
}

// FindByIdentityNO godoc
// @Summary      Search users by identity number
// @Description  Find users matching the given identity number with pagination
// @Tags         Auth
// @Produce      json
// @Param        identity_no  query     string  true   "Identity number (full or partial)"
// @Param        page         query     int     false  "Page number"      default(1)
// @Param        size         query     int     false  "Page size"        default(10)
// @Param        orderBy      query     string  false  "Sort field"
// @Success      200          {object}  models.UsersList
// @Failure      400          {object}  httpErrors.RestError
// @Failure      500          {object}  httpErrors.RestError
// @Router       /auth/find [get]
// @Security     BearerAuth
func (h *authHandlers) FindByIdentityNO() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		if c.QueryParam("identity_no") == "" {
			utils.LogResponseError(c, h.logger, httpErrors.NewBadRequestError("identity_no is required"))
			return c.JSON(http.StatusBadRequest, httpErrors.NewBadRequestError("identity_no is required"))
		}

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		response, err := h.authUC.FindByIdentity(ctx, c.QueryParam("identity_no"), paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, response)
	}
}

// GetUsers godoc
// @Summary      List all users (paginated)
// @Description  Retrieve a paginated list of all users (admin only)
// @Tags         Auth
// @Produce      json
// @Param        page     query     int     false  "Page number"  default(1)
// @Param        size     query     int     false  "Page size"    default(10)
// @Param        orderBy  query     string  false  "Sort field"
// @Success      200      {object}  models.UsersList
// @Failure      400      {object}  httpErrors.RestError
// @Failure      403      {object}  httpErrors.RestError  "Forbidden"
// @Failure      500      {object}  httpErrors.RestError
// @Router       /auth/all [get]
// @Security     BearerAuth
func (h *authHandlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		paginationQuery, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		usersList, err := h.authUC.GetUsers(ctx, paginationQuery)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, usersList)
	}
}

// GetMe godoc
// @Summary      Get current authenticated user
// @Description  Returns the profile of the currently logged-in user
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  models.User
// @Failure      401  {object}  httpErrors.RestError  "Unauthorized"
// @Failure      500  {object}  httpErrors.RestError
// @Router       /auth/me [get]
// @Security     BearerAuth
func (h *authHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			return utils.ErrResponseWithLog(c, h.logger, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		return c.JSON(http.StatusOK, user)
	}
}
