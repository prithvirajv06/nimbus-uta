package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GenericRepository provides generic CRUD operations for a MongoDB collection.
type GenericRepository[T any] struct {
	Coll *mongo.Collection
	Ctx  context.Context
}

// NewGenericRepository creates a new GenericRepository.
func NewGenericRepository[T any](ctx context.Context, db *mongo.Database, collectionName string) *GenericRepository[T] {
	coll := db.Collection(collectionName)
	return &GenericRepository[T]{Coll: coll, Ctx: ctx}
}
func (r *GenericRepository[T]) InsertOne(doc T) (*mongo.InsertOneResult, error) {
	return r.Coll.InsertOne(r.Ctx, doc)
}

func (r *GenericRepository[T]) FindOne(filter interface{}, opts ...*options.FindOneOptions) (*T, error) {
	var result T
	err := r.Coll.FindOne(r.Ctx, filter, opts...).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *GenericRepository[T]) UpdateOne(filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return r.Coll.UpdateOne(r.Ctx, filter, update, opts...)
}

func (r *GenericRepository[T]) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return r.Coll.DeleteOne(r.Ctx, filter, opts...)
}

func (r *GenericRepository[T]) FindMany(filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := r.Coll.Find(r.Ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.Ctx)

	var results []T
	for cursor.Next(r.Ctx) {
		var elem T
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GenericRepository[T]) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return r.Coll.CountDocuments(r.Ctx, filter, opts...)
}

func (r *GenericRepository[T]) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) ([]T, error) {
	cursor, err := r.Coll.Aggregate(r.Ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.Ctx)

	var results []T
	for cursor.Next(r.Ctx) {
		var elem T
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, elem)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *GenericRepository[T]) CreateIndex(model mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return r.Coll.Indexes().CreateOne(r.Ctx, model, opts...)
}

func (r *GenericRepository[T]) DropIndex(name string, opts ...*options.DropIndexesOptions) (bson.Raw, error) {
	return r.Coll.Indexes().DropOne(r.Ctx, name, opts...)
}

func (r *GenericRepository[T]) ListIndexes(opts ...*options.ListIndexesOptions) ([]bson.M, error) {
	cursor, err := r.Coll.Indexes().List(r.Ctx, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(r.Ctx)

	var indexes []bson.M
	for cursor.Next(r.Ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			return nil, err
		}
		indexes = append(indexes, index)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return indexes, nil
}
