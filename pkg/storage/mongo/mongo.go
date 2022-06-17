package mongo

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Хранилище данных.
type Store struct {
	db *mongo.Client
}

const (
	dbName         = "go-news"
	collectionName = "posts"
)

// Конструктор объекта хранилища.
func New(connectionString string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: client,
	}
	return &s, err
}

func (s *Store) Posts() ([]storage.Post, error) {
	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var posts []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, cur.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.db.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return err
}
func (s *Store) UpdatePost(p storage.Post) error {
	collection := s.db.Database(dbName).Collection(collectionName)
	filter := bson.D{{Key: "ID", Value: p.ID}}
	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "Title", Value: p.Title},
			{Key: "Content", Value: p.Content},
			{Key: "AuthorID", Value: p.AuthorID},
			{Key: "AuthorName", Value: p.AuthorName},
			{Key: "PublishedAt", Value: p.PublishedAt},
		},
	}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return err
}
func (s *Store) DeletePost(storage.Post) error {
	return nil
}
