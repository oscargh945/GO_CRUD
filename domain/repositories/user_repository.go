package repositories

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oscargh945/go-crud-graphql/domain"
	"github.com/oscargh945/go-crud-graphql/domain/entities"
	"github.com/oscargh945/go-crud-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type UserRepository struct {
	Client *mongo.Database
}

func (r *UserRepository) GetAllUsers() ([]*entities.User, error) {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var users []*entities.User
	filter := bson.M{"softDeleted": false}
	cur, err := collec.Find(ctx, filter)
	if err != nil {
		log.Fatal("No sirve esta kgaa", err)
	}
	if err = cur.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetUser(id string) (*entities.User, error) {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id, "softDeleted": false}
	var users *entities.User
	err := collec.FindOne(ctx, filter).Decode(&users)
	if err != nil {
		return nil, fmt.Errorf("El usuario esta eliminado")
	}

	if users.SoftDeleted {
		return nil, fmt.Errorf("El usuario esta eliminado")
	}

	return users, nil
}

func (r *UserRepository) PaginationSearchUsers(SearchUser *string, page, pageSize *int) ([]*entities.User, error) {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pagination := (*page - 1) * *pageSize
	searchPattern := bson.M{"$regex": *SearchUser, "$options": "i"}
	filter := bson.M{"name": searchPattern, "softDeleted": false}
	var option = options.Find().SetSort(bson.D{{Key: "name", Value: 1}}).SetLimit(int64(*pageSize)).SetSkip(int64(pagination))
	cur, err := collec.Find(ctx, filter, option)
	if err != nil {
		return nil, fmt.Errorf("Error al encontrar los usuarios")
	}
	defer cur.Close(ctx)
	var users []*entities.User
	for cur.Next(ctx) {
		var user *entities.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, fmt.Errorf("Error")
		}
		users = append(users, user)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar users")
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(id string, userInfo model.UpdateUserInput) *entities.User {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateUserInfo := bson.M{}
	if userInfo.Name != nil {
		updateUserInfo["name"] = userInfo.Name
	}
	if userInfo.Email != nil {
		updateUserInfo["email"] = userInfo.Email
	}
	if userInfo.Phone != nil {
		updateUserInfo["phone"] = userInfo.Phone
	}

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateUserInfo}
	option := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := collec.FindOneAndUpdate(ctx, filter, update, option)

	var user entities.User
	if err := result.Decode(&user); err != nil {
		log.Fatal("No sirve esta kgaaa", err)
	}
	return &user
}

func (r *UserRepository) CreateUser(userInfo entities.CreateUserInput) (*entities.User, error) {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if r.ExistingEmailUser(userInfo.Email) {
		return nil, fmt.Errorf("El email ingresado ya esta en uso")
	}
	if r.ExistingPhoneUser(userInfo.Phone) {
		return nil, fmt.Errorf("El phone ingresado ya esta en uso")
	}
	hash, err := domain.HashPassword(userInfo.Password)
	if err != nil {
		log.Fatal("No se hasheo la password")
	}
	inserg, err := collec.InsertOne(ctx, bson.M{"name": userInfo.Name, "email": userInfo.Email, "phone": userInfo.Phone, "password": hash, "softDeleted": userInfo.SoftDeleted})

	if err != nil {
		log.Fatal("No se inserto el usuario", err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnGetAllUsers := entities.User{Id: insertedID, Name: userInfo.Name, Email: userInfo.Email, Phone: userInfo.Phone, SoftDeleted: userInfo.SoftDeleted}
	return &returnGetAllUsers, nil
}

func (r *UserRepository) SoftDeleteUser(id string) *model.DeleteUserResponse {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	//_, err := userCollec.DeleteOne(ctx, filter)
	_, err := collec.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"softDeleted": true}})
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteUserResponse{DeleteUserID: id}
}

func (r *UserRepository) ExistingEmailUser(email string) bool {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"email": email}
	count, err := collec.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal("No se encontro un usuario con el email ingresado", err)
	}
	return count > 0
}

func (r *UserRepository) ExistingPhoneUser(phone string) bool {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{"phone": phone}
	count, err := collec.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal("No se encontro un usuario con el phone ingresado", err)
	}
	return count > 0
}

func (r *UserRepository) Login(userInfo model.LoginInput) (*model.LoginResponse, error) {
	collec := r.Client.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user entities.User

	filter := bson.M{"email": userInfo.Email}
	err := collec.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("Email no existente")
	}

	// Validate the password
	isValid := domain.ValidPassword(userInfo.Password, user.Password)
	if !isValid {
		return nil, fmt.Errorf("La contraseña es incorrecta")
	}

	// Generate the tokens
	tokens, err := r.GenerateTokens(user)
	if err != nil {
		return nil, err
	}

	response := model.LoginResponse{
		User:         &user,
		TokenAccess:  tokens["access"],
		TokenRefresh: tokens["refresh"],
	}
	return &response, nil
}

func (r *UserRepository) GenerateTokens(user entities.User) (map[string]string, error) {
	tokens := make(map[string]string)

	jti := uuid.New()
	rti := uuid.New()

	accessDuration := time.Minute * 24
	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Id,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"exp":   time.Now().Add(accessDuration).Unix(),
		"jti":   jti,
		"rti":   rti,
	})

	refreshDuration := time.Minute * 25
	tokenRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.Id,
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
		"exp":   time.Now().Add(refreshDuration).Unix(),
		"jti":   rti,
		"rti":   jti,
	})

	tokenAccessString, err := tokenAccess.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}
	tokenRefreshString, err := tokenRefresh.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return nil, err
	}

	tokens["access"] = tokenAccessString
	tokens["refresh"] = tokenRefreshString

	return tokens, nil
}

func (r *UserRepository) ValidateJWT(tokenString string, secretKey []byte) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		response := map[string]interface{}{
			"sub":   claims["sub"].(string),
			"name":  claims["name"].(string),
			"email": claims["email"].(string),
			"phone": claims["phone"].(string),
		}
		return response, nil
	} else {
		return nil, fmt.Errorf("Token invalido")
	}
}

func (r *UserRepository) Refresh(refreshToken string) (*model.LoginResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, fmt.Errorf("Token invalido")
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, fmt.Errorf("token Refresh expirado")
	}

	var user entities.User
	tokens, _ := r.GenerateTokens(user)
	response := &model.LoginResponse{
		User:         &user,
		TokenAccess:  tokens["access"],
		TokenRefresh: tokens["refresh"],
	}
	return response, nil
}
