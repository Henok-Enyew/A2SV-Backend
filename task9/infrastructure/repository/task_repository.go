package repository

import (
	"context"
	"errors"
	"task9/domain/entity"
	"task9/domain/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) repository.TaskRepository {
	return &taskRepository{
		collection: collection,
	}
}

func (r *taskRepository) Create(task *entity.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID := primitive.NewObjectID()
	task.ID = objectID.Hex()

	doc := bson.M{
		"_id":        objectID,
		"title":      task.Title,
		"description": task.Description,
		"due_date":   task.DueDate,
		"status":     task.Status,
		"created_at": task.CreatedAt,
		"updated_at": task.UpdatedAt,
	}

	_, err := r.collection.InsertOne(ctx, doc)
	return err
}

func (r *taskRepository) FindByID(id string) (*entity.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var doc bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return r.toEntity(doc), nil
}

func (r *taskRepository) FindAll() ([]*entity.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*entity.Task
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		tasks = append(tasks, r.toEntity(doc))
	}

	return tasks, nil
}

func (r *taskRepository) Update(task *entity.Task) error {
	objectID, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"due_date":    task.DueDate,
			"status":      task.Status,
			"updated_at":  task.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *taskRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

func (r *taskRepository) toEntity(doc bson.M) *entity.Task {
	id := doc["_id"].(primitive.ObjectID).Hex()
	
	var dueDate, createdAt, updatedAt time.Time
	if dt, ok := doc["due_date"].(primitive.DateTime); ok {
		dueDate = dt.Time()
	} else if t, ok := doc["due_date"].(time.Time); ok {
		dueDate = t
	}
	if dt, ok := doc["created_at"].(primitive.DateTime); ok {
		createdAt = dt.Time()
	} else if t, ok := doc["created_at"].(time.Time); ok {
		createdAt = t
	}
	if dt, ok := doc["updated_at"].(primitive.DateTime); ok {
		updatedAt = dt.Time()
	} else if t, ok := doc["updated_at"].(time.Time); ok {
		updatedAt = t
	}
	
	return &entity.Task{
		ID:          id,
		Title:       doc["title"].(string),
		Description: doc["description"].(string),
		DueDate:     dueDate,
		Status:      doc["status"].(string),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

