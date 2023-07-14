package repositories

import (
	"context"
	"fmt"
	"github.com/oscargh945/go-crud-graphql/domain/entities"
	"github.com/oscargh945/go-crud-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
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
	//searchPattern := bson.M{"$regex": SearchUser, "$name": ""}
	pagination := (*page - 1) * *pageSize
	//findOptions := options.Find()
	//findOptions.SetLimit(int64(pageSize))
	//findOptions.SetSkip(int64(pagination))
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
	inserg, err := collec.InsertOne(ctx, bson.M{"name": userInfo.Name, "email": userInfo.Email, "phone": userInfo.Phone, "softDeleted": userInfo.SoftDeleted})

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
