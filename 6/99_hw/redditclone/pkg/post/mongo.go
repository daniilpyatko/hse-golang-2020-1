package post

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo is located here and not in the separate folder for convenience purposes

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type CollectionHelper interface {
	Find(context.Context, interface{}, ...*options.FindOptions) (CursorHelper, error)
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) error
	UpdateOne(context.Context, interface{}, interface{}) error
	DeleteOne(context.Context, interface{}) error
}

type CursorHelper interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context)
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	res *mongo.SingleResult
}

type mongoCursor struct {
	cur *mongo.Cursor
}

func (c *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorHelper, error) {
	r, err := c.coll.Find(ctx, filter, opts...)
	return &mongoCursor{cur: r}, err
}

func (c *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	f := c.coll.FindOne(ctx, filter)
	return &mongoSingleResult{res: f}
}

func (c *mongoCollection) InsertOne(ctx context.Context, filter interface{}) error {
	_, err := c.coll.InsertOne(ctx, filter)
	return err
}

func (c *mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) error {
	_, err := c.coll.UpdateOne(ctx, filter, update)
	return err
}

func (c *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) error {
	_, err := c.coll.DeleteOne(ctx, filter)
	return err
}

func (s *mongoSingleResult) Decode(v interface{}) error {
	return s.res.Decode(v)
}

func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.cur.Next(ctx)
}

func (c *mongoCursor) Decode(v interface{}) error {
	return c.cur.Decode(v)
}

func (c *mongoCursor) Close(ctx context.Context) {
	c.cur.Close(ctx)
}
