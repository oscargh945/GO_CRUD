// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/oscargh945/go-crud-graphql/domain/entities"
)

type DeleteUserResponse struct {
	DeleteUserID string `json:"deleteUserId"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         *entities.User `json:"user"`
	TokenAccess  string         `json:"tokenAccess"`
	TokenRefresh string         `json:"tokenRefresh"`
}

type UpdateUserInput struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}
