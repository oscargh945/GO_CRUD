package graph

import "github.com/oscargh945/go-crud-graphql/domain/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserUseCase usecase.UserUseCase
}
