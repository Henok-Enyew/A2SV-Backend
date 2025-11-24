package data

import (
	"context"
	"errors"
	"task6/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskService struct {
	collection *mongo.Collection
}

func NewTaskService(collection *mongo.Collection) *TaskService {
	return &TaskService{
		collection: collection,
	}
}

func (ts *TaskService) GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ts.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (ts *TaskService) GetTaskByID(id string) (models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var task models.Task
	err = ts.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, err
	}

	return task, nil
}

func (ts *TaskService) CreateTask(req models.CreateTaskRequest) (models.Task, error) {
	now := time.Now()
	status := req.Status
	if status == "" {
		status = "pending"
	}

	task := models.Task{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Status:      status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ts.collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (ts *TaskService) UpdateTask(id string, req models.UpdateTaskRequest) (models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{}
	if req.Title != "" {
		update["title"] = req.Title
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if !req.DueDate.IsZero() {
		update["due_date"] = req.DueDate
	}
	if req.Status != "" {
		update["status"] = req.Status
	}
	update["updated_at"] = time.Now()

	filter := bson.M{"_id": objectID}
	updateDoc := bson.M{"$set": update}

	result := ts.collection.FindOneAndUpdate(
		ctx,
		filter,
		updateDoc,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return models.Task{}, errors.New("task not found")
		}
		return models.Task{}, result.Err()
	}

	var task models.Task
	if err := result.Decode(&task); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (ts *TaskService) DeleteTask(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID format")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := ts.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}

	return nil
}

