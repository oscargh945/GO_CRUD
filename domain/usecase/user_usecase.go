package usecase

import (
	"github.com/oscargh945/go-crud-graphql/domain/entities"
	"github.com/oscargh945/go-crud-graphql/domain/repositories"
	"github.com/oscargh945/go-crud-graphql/graph/model"
)

type UserUseCase struct {
	Repository repositories.UserRepository
}

func (u *UserUseCase) GetAllUsersUseCase() ([]*entities.User, error) {
	return u.Repository.GetAllUsers()
}

func (u *UserUseCase) GetUserUseCase(id string) (*entities.User, error) {
	return u.Repository.GetUser(id)
}

func (u *UserUseCase) PaginationSearchUsersUseCase(SearchUser *string, page, pageSize *int) ([]*entities.User, error) {
	return u.Repository.PaginationSearchUsers(SearchUser, page, pageSize)
}

func (u *UserUseCase) CreateUserUseCase(user entities.CreateUserInput) (*entities.User, error) {
	return u.Repository.CreateUser(user)
}

func (u *UserUseCase) UpdateUserUseCase(id string, user model.UpdateUserInput) *entities.User {
	return u.Repository.UpdateUser(id, user)
}

func (u *UserUseCase) SoftDeleteUserUseCase(id string) *model.DeleteUserResponse {
	return u.Repository.SoftDeleteUser(id)
}

func (u *UserUseCase) LoginUseCase(user model.LoginInput) (*model.LoginResponse, error) {
	return u.Repository.Login(user)
}

func (u *UserUseCase) RefreshUseCase(refreshToken string) (*model.LoginResponse, error) {
	return u.Repository.Refresh(refreshToken)
}
