package repository

import (
	"context"
	"errors"
	"task9/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepositoryMongo struct {
	collection *mongo.Collection
}

func NewTaskRepositoryMongo(collection *mongo.Collection) domain.TaskRepository {
	return &TaskRepositoryMongo{collection: collection}
}

func (r *TaskRepositoryMongo) GetAll() ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []domain.Task
	for cursor.Next(ctx) {
		var taskDoc bson.M
		if err := cursor.Decode(&taskDoc); err != nil {
			continue
		}
		task := r.mapToDomain(taskDoc)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepositoryMongo) GetByID(id string) (domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var taskDoc bson.M
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&taskDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, errors.New("task not found")
		}
		return domain.Task{}, err
	}

	return r.mapToDomain(taskDoc), nil
}

func (r *TaskRepositoryMongo) Create(task domain.Task) (domain.Task, error) {
	objectID := primitive.NewObjectID()
	task.ID = objectID.Hex()

	doc := r.mapToDocument(task)
	doc["_id"] = objectID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (r *TaskRepositoryMongo) Update(id string, task domain.Task) (domain.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{}
	if task.Title != "" {
		update["title"] = task.Title
	}
	if task.Description != "" {
		update["description"] = task.Description
	}
	if !task.DueDate.IsZero() {
		update["due_date"] = task.DueDate
	}
	if task.Status != "" {
		update["status"] = task.Status
	}
	update["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID}
	updateDoc := bson.M{"$set": update}

	result := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		updateDoc,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return domain.Task{}, errors.New("task not found")
		}
		return domain.Task{}, result.Err()
	}

	var taskDoc bson.M
	if err := result.Decode(&taskDoc); err != nil {
		return domain.Task{}, err
	}

	return r.mapToDomain(taskDoc), nil
}

func (r *TaskRepositoryMongo) Delete(id string) error {
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

func (r *TaskRepositoryMongo) mapToDomain(doc bson.M) domain.Task {
	task := domain.Task{}
	if id, ok := doc["_id"].(primitive.ObjectID); ok {
		task.ID = id.Hex()
	}
	if title, ok := doc["title"].(string); ok {
		task.Title = title
	}
	if desc, ok := doc["description"].(string); ok {
		task.Description = desc
	}
	if dueDate, ok := doc["due_date"].(primitive.DateTime); ok {
		task.DueDate = dueDate.Time()
	} else if dueDate, ok := doc["due_date"].(time.Time); ok {
		task.DueDate = dueDate
	}
	if status, ok := doc["status"].(string); ok {
		task.Status = status
	}
	if createdAt, ok := doc["created_at"].(primitive.DateTime); ok {
		task.CreatedAt = createdAt.Time()
	} else if createdAt, ok := doc["created_at"].(time.Time); ok {
		task.CreatedAt = createdAt
	}
	if updatedAt, ok := doc["updated_at"].(primitive.DateTime); ok {
		task.UpdatedAt = updatedAt.Time()
	} else if updatedAt, ok := doc["updated_at"].(time.Time); ok {
		task.UpdatedAt = updatedAt
	}
	return task
}

func (r *TaskRepositoryMongo) mapToDocument(task domain.Task) bson.M {
	doc := bson.M{
		"title":       task.Title,
		"description": task.Description,
		"due_date":    task.DueDate,
		"status":      task.Status,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	}
	return doc
}

