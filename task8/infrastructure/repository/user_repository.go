package repository

import (
	"context"
	"errors"
	"task8/domain/entity"
	"task8/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) repository.UserRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) Create(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID := primitive.NewObjectID()
	user.ID = objectID.Hex()

	doc := bson.M{
		"_id":      objectID,
		"username": user.Username,
		"password": user.Password,
		"role":     user.Role,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc bson.M
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.toEntity(doc), nil
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return r.toEntity(doc), nil
}

func (r *userRepository) Update(user *entity.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"username": user.Username,
			"password": user.Password,
			"role":     user.Role,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *userRepository) toEntity(doc bson.M) *entity.User {
	id := doc["_id"].(primitive.ObjectID).Hex()
	return &entity.User{
		ID:       id,
		Username: doc["username"].(string),
		Password: doc["password"].(string),
		Role:     doc["role"].(string),
	}
}


