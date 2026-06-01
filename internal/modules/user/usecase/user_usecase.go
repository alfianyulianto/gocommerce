package usecase

import (
	"context"
	"math"

	"github.com/alfianyulianto/gocommerce/internal/modules/user/dto"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/entity"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/repository"
	"github.com/alfianyulianto/gocommerce/internal/modules/user/search"
	"github.com/alfianyulianto/gocommerce/pkg/response"
	"github.com/alfianyulianto/gocommerce/pkg/validation"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	Create(ctx context.Context, request *dto.CreateUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, request *dto.UpdateUserRequest) (*dto.UserResponse, error)
	FindById(ctx context.Context, id string) (*dto.UserResponse, error)
	Delete(ctx context.Context, id string) error
	FindAll(ctx context.Context, request *dto.UserFilterRequest) ([]dto.UserResponse, *response.Pagination, error)
	Search(ctx context.Context, request *dto.SearchUserRequest) ([]dto.UserResponse, *response.Pagination, error)
}

type userUseCase struct {
	Repository repository.UserRepository
	Log        *logrus.Logger
	ESSearch   search.ESSearch
}

func NewUserUseCase(repository repository.UserRepository, log *logrus.Logger, ESSearch search.ESSearch) UserUseCase {
	return &userUseCase{Repository: repository, Log: log, ESSearch: ESSearch}
}

func (c *userUseCase) Create(ctx context.Context, request *dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.NewString(),
		Name:     request.Name,
		Email:    request.Email,
		Password: string(password),
		Phone:    request.Phone,
	}

	if err = c.Repository.Create(ctx, user); err != nil {
		return nil, err
	}

	go func() {
		err = c.ESSearch.Index(ctx, user)
		if err != nil {
			c.Log.WithField("user", user).WithError(err).Error("Failed to index create user in elasticsearch")
		}
	}()

	return dto.ToUserResponse(user), nil
}

func (c *userUseCase) Update(ctx context.Context, request *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, err
	}

	user := new(entity.User)
	if err := c.Repository.FindById(ctx, user, request.ID); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	user.Name = request.Name
	user.Email = request.Email

	if request.Password != nil {
		password, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(password)
	}

	if request.Phone != nil {
		user.Phone = request.Phone
	}

	if err := c.Repository.Update(ctx, user); err != nil {
		return nil, err
	}

	go func() {
		err := c.ESSearch.Index(ctx, user)
		if err != nil {
			c.Log.WithField("user", user).WithError(err).Error("Failed to index update user in elasticsearch")
		}
	}()

	return dto.ToUserResponse(user), nil
}

func (c *userUseCase) FindById(ctx context.Context, id string) (*dto.UserResponse, error) {
	user := new(entity.User)
	if err := c.Repository.FindById(ctx, user, id); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return dto.ToUserResponse(user), nil
}

func (c *userUseCase) Delete(ctx context.Context, id string) error {
	user := new(entity.User)
	if err := c.Repository.FindById(ctx, user, id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if err := c.Repository.Delete(ctx, user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	go func() {
		err := c.ESSearch.Delete(ctx, id)
		if err != nil {
			c.Log.WithField("user", user).WithError(err).Error("Failed to delete user in elasticsearch")
		}
	}()

	return nil
}

func (c *userUseCase) FindAll(ctx context.Context, request *dto.UserFilterRequest) ([]dto.UserResponse, *response.Pagination, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, nil, err
	}

	filter := dto.ToUserFilter(request)
	users, total, err := c.Repository.FindAll(ctx, filter)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponse := dto.ToUserResponse(&user)
		responses = append(responses, *userResponse)
	}

	totalPage := int(math.Ceil(float64(total) / float64(request.PerPage)))

	pagination := &response.Pagination{
		CurrentPage: request.Page,
		PerPage:     request.PerPage,
		TotalItem:   total,
		TotalPage:   totalPage,
		HasNext:     request.Page < totalPage,
		HasPrev:     request.Page > 1 && request.Page <= totalPage,
	}

	return responses, pagination, nil
}

func (c *userUseCase) Search(ctx context.Context, request *dto.SearchUserRequest) ([]dto.UserResponse, *response.Pagination, error) {
	if err := validation.Validate.Struct(request); err != nil {
		return nil, nil, err
	}

	users, total, err := c.ESSearch.Search(ctx, request)
	if err != nil {
		return nil, nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(request.PerPage)))

	pagination := &response.Pagination{
		CurrentPage: request.Page,
		PerPage:     request.PerPage,
		TotalItem:   total,
		TotalPage:   totalPage,
		HasNext:     request.Page < totalPage,
		HasPrev:     request.Page > 1 && request.Page <= totalPage,
	}

	return users, pagination, nil
}
