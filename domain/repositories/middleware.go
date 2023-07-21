package repositories

import (
	"context"
	"errors"
	"net/http"
	"os"
)

type UserCtx struct {
	Id    string
	Name  string
	Email string
	Phone string
}

const currentUserKey = "currentUser"

func (r *UserRepository) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		auth := request.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, request)
			return
		}

		userData, err := r.ValidateJWT(auth, []byte(os.Getenv("SECRET")))
		if err != nil || userData == nil {
			next.ServeHTTP(w, request)
			return
		}
		user := UserCtx{
			Id:    userData["sub"].(string),
			Name:  userData["name"].(string),
			Email: userData["email"].(string),
			Phone: userData["phone"].(string),
		}
		ctx := context.WithValue(request.Context(), currentUserKey, user)
		next.ServeHTTP(w, request.WithContext(ctx))

	})
}

func (r *UserRepository) GetCurrentUserFromCTX(ctx context.Context) (*UserCtx, error) {
	errNoUserInContext := errors.New("No hay ningun usuario en contexto")
	if ctx.Value(currentUserKey) == nil {
		return nil, errNoUserInContext
	}
	user, ok := ctx.Value(currentUserKey).(UserCtx)
	if !ok || user.Id == "" || user.Email == "" {
		return nil, errNoUserInContext
	}
	return &user, nil
}
