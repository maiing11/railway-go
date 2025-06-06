package usecase

import (
	"context"
	"errors"
	"railway-go/internal/constant/model"
	"railway-go/internal/repository"
	"railway-go/internal/utils"
	"railway-go/internal/utils/token"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type UserSessionUC interface {
	Register(ctx context.Context, request *model.RegisterUserRequest) error
	Login(ctx context.Context, request *model.LoginUserRequest) (map[string]any, error)
	Logout(ctx context.Context, sessionID string) error
	RenewAccessToken(ctx context.Context, refreshToken string) (string, error)
	CreateGuestSession(ctx context.Context) (model.Session, error)
	GetSession(ctx context.Context, sessionID string) (model.Session, error)
	RegisterAdmin(ctx context.Context, request *model.RegisterAdminRequest) error
	GetUserIDFromSession(ctx context.Context, sessionID string) (*uuid.UUID, error)
}

type UserSessionUsecase struct {
	*UseCase
	TokenMaker token.Maker
	config     *viper.Viper
}

func NewUserSessionUsecase(
	useCase *UseCase,
	tokenMaker token.Maker,
	config *viper.Viper,
) *UserSessionUsecase {
	return &UserSessionUsecase{
		UseCase:    useCase,
		TokenMaker: tokenMaker,
		config:     config,
	}
}

func (uc *UserSessionUsecase) Register(ctx context.Context, request *model.RegisterUserRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = uc.Validate.Struct(request)
	if err != nil {

		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "invalid request body")
	}

	count, err := tx.CountUserByEmail(ctx, request.Email)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to count user from database")
	}

	if count > 0 {
		return utils.WrapError(fiber.StatusConflict, uc.Log, utils.Error, err, "email already exists")
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to hash password")
	}
	id := uuid.New()
	user := repository.CreateUserParams{
		ID:          id,
		Name:        request.Name,
		Email:       request.Email,
		Password:    hashedPassword,
		PhoneNumber: request.PhoneNumber,
	}

	if err = tx.CreateUser(ctx, user); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create user")
	}

	if err := tx.Commit(ctx); err != nil {
		return utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	uc.Log.Info("user registered successfully", zap.String("email", user.Email))
	return nil
}

/*
		Login handles the user login process. It validates the login request, checks the user credentials,
	 generates a session, creates access and refresh tokens, and sets the session ID in cookies.
*/

func nullUserRoleToString(role repository.NullUserRole) string {
	if role.Valid {
		return string(role.UserRole)
	}
	return "user"
}

func (uc *UserSessionUsecase) Login(ctx context.Context, request *model.LoginUserRequest) (map[string]any, error) {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if ctx == nil {
		return nil, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Error, err, "context is nil")
	}

	fiberCtx, ok := ctx.Value("fiberCtx").(*fiber.Ctx)
	if !ok || fiberCtx == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "fiber context is nil")
	}

	if err := uc.Validate.Struct(request); err != nil {
		return nil, utils.WrapError(fiber.StatusBadRequest, uc.Log, utils.Error, err, "invalid request body")
	}

	user, err := tx.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusUnauthorized, uc.Log, utils.Error, err, "user not found")
	}

	err = utils.CheckPassword(request.Password, user.Password)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusUnauthorized, uc.Log, utils.Error, err, "invalid password")
	}

	// generate session
	session := &model.Session{
		ID:           uuid.NewString(),
		RefreshToken: "",
		UserID:       &user.ID,
		Role:         nullUserRoleToString(user.Role),
		UserAgent:    fiberCtx.Get("User-Agent"),
		ClientIP:     fiberCtx.IP(),
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(uc.config.GetDuration("Token.RefreshTokenDuration")),
	}

	// commit
	if err := tx.Commit(ctx); err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to commit transaction")
	}

	// save session in redis
	if err := uc.Repo.CreateSession(ctx, session); err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create session")
	}

	// generate tokens with sesssion ID
	accessDuration := uc.config.GetDuration("Token.AccessTokenDuration")
	accessToken, accessPayload, err := uc.TokenMaker.CreateToken(&user.ID, session.ID, session.Role, accessDuration)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create access token")
	}
	// uc.Log.Info("generate access token", zap.Any("payload", accessToken))

	refreshDuration := uc.config.GetDuration("Token.RefreshTokenDuration")
	refreshToken, refreshPayload, err := uc.TokenMaker.CreateToken(&user.ID, session.ID, session.Role, refreshDuration)
	if err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to create refresh token")
	}

	// update session with refresh token
	session.RefreshToken = refreshToken
	if err := uc.Repo.UpdateSessionAccessToken(ctx, session.ID, accessToken, session.ExpiresAt); err != nil {
		return nil, utils.WrapError(fiber.StatusInternalServerError, uc.Log, utils.Error, err, "failed to update session with new access token")
	}

	// set session ID in cookies
	cookie := &fiber.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   true,
	}

	fiberCtx.Cookie(cookie)

	response := map[string]any{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]any{
			"id":    user.ID,
			"email": user.Email,
			"role":  session.Role,
		},
		"access_payload":  accessPayload,
		"refresh_payload": refreshPayload,
	}
	uc.Log.Info("Login successful", zap.Any("userID", user.ID))
	return response, nil

}

//	func (u *UserSessionUsecase) ForgotPassword(ctx context.Context, email string) error {
//		user, err := u.Repo.GetUserByEmail(ctx, email)
//	}
//

func (uc *UserSessionUsecase) RegisterAdmin(ctx context.Context, request *model.RegisterAdminRequest) error {
	tx, err := uc.Repo.BeginTransaction(ctx)
	if err != nil {
		uc.Log.Warn("failed to begin transaction", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("invalid request body", zap.Any("request", request))
		return fiber.ErrBadRequest
	}

	count, err := tx.CountUserByEmail(ctx, request.Email)
	if err != nil {
		uc.Log.Warn("failed to count user from database", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	if count > 0 {
		uc.Log.Warn("email already exists", zap.String("email", request.Email))
		return fiber.ErrConflict
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		uc.Log.Warn("failed to hash password", zap.Error(err))
		return fiber.ErrInternalServerError
	}
	role := repository.NullUserRole{UserRole: repository.UserRole(request.Role), Valid: true}

	id := uuid.New()
	user := repository.CreateUserParams{
		ID:          id,
		Name:        request.Name,
		Email:       request.Email,
		Password:    hashedPassword,
		Role:        role,
		PhoneNumber: request.PhoneNumber,
	}

	if err = tx.CreateUser(ctx, user); err != nil {
		uc.Log.Warn("failed to create user", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit(ctx); err != nil {
		uc.Log.Warn("failed to commit transaction", zap.Error(err))
		return fiber.ErrInternalServerError
	}

	uc.Log.Info("user registered successfully", zap.String("email", user.Email))
	return nil
}

func (uc *UserSessionUsecase) Logout(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		uc.Log.Error("Session ID is empty")
		return errors.New("invalid session")

	}
	err := uc.Repo.DeleteSession(ctx, sessionID)
	if err != nil {
		uc.Log.Error("Failed to log out", zap.String("sessionID", sessionID), zap.Error(err))
		return err
	}
	uc.Log.Info("User logged out successfully", zap.String("sessionID", sessionID))
	return nil
}
