package repository

import (
	"context"
	"errors"
	"task8/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct {
	collection *mongo.Collection
}

func NewUserRepositoryMongo(collection *mongo.Collection) domain.UserRepository {
	return &UserRepositoryMongo{collection: collection}
}

func (r *UserRepositoryMongo) Create(user domain.User) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return domain.User{}, errors.New("username already exists")
	}
	if err != mongo.ErrNoDocuments {
		return domain.User{}, err
	}

	objectID := primitive.NewObjectID()
	user.ID = objectID.Hex()

	doc := r.mapToDocument(user)
	doc["_id"] = objectID

	_, err = r.collection.InsertOne(ctx, doc)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepositoryMongo) GetByUsername(username string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userDoc bson.M
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&userDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return r.mapToDomain(userDoc), nil
}

func (r *UserRepositoryMongo) GetByID(id string) (domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, errors.New("invalid user ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userDoc bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&userDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}

	return r.mapToDomain(userDoc), nil
}

func (r *UserRepositoryMongo) UpdateRole(username string, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"role": role}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepositoryMongo) IsFirstUser() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *UserRepositoryMongo) mapToDomain(doc bson.M) domain.User {
	user := domain.User{}
	if id, ok := doc["_id"].(primitive.ObjectID); ok {
		user.ID = id.Hex()
	}
	if username, ok := doc["username"].(string); ok {
		user.Username = username
	}
	if password, ok := doc["password"].(string); ok {
		user.Password = password
	}
	if role, ok := doc["role"].(string); ok {
		user.Role = role
	}
	return user
}

func (r *UserRepositoryMongo) mapToDocument(user domain.User) bson.M {
	doc := bson.M{
		"username": user.Username,
		"password": user.Password,
		"role":     user.Role,
	}
	return doc
}

