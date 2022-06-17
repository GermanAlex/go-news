package mongo

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Хранилище данных.
type Store struct {
	db *mongo.Client
}

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
	return posts, nil
}

func (s *Store) AddPost(storage.Post) error {
	return nil
}
func (s *Store) UpdatePost(storage.Post) error {
	return nil
}
func (s *Store) DeletePost(storage.Post) error {
	return nil
}
