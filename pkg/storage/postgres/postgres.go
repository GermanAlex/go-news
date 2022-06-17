package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор объекта хранилища.
func New(connectionString string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
			p.id,
			p.title,
			p.content,
			p.author_id,
			a.name,
			p.created_at,
			p.published_at
		FROM posts p, authors a
		WHERE
			p.author_id = a.id
		ORDER BY id;
	`,
	)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (s *Store) AddPost(p storage.Post) error {
	var id int
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = s.addAuthor(p)

	if err != nil {
		return err
	}

	err = tx.QueryRow(context.Background(), `
		INSERT INTO posts (id, author_id, title, content, created_at, published_at)
		VALUES ($1, $2, $3, $4, $5, $6);
		`,
		p.ID,
		p.AuthorID,
		p.Title,
		p.Content,
		p.CreatedAt,
		p.PublishedAt,
	).Scan(&id)

	return err
}
func (s *Store) UpdatePost(p storage.Post) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	commandTag, err := tx.Exec(context.Background(), `
		UPDATE posts SET
			author_id = $1,
			title = $2,
			content = $3,
			published_at = $4
		WHERE id = $5;
		`,
		p.AuthorID,
		p.Title,
		p.Content,
		p.PublishedAt,
		p.ID)

	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}

	return err

}
func (s *Store) DeletePost(p storage.Post) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	commandTag, err := tx.Exec(context.Background(), `
		DELETE FROM posts
		where id = $1;
		`,
		p.ID)

	if commandTag.RowsAffected() != 1 {
		return errors.New("no row found to delete")
	}

	return err
}

func (s *Store) addAuthor(p storage.Post) error {
	var id int
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), `
		SELECT id
		FROM authors
		WHERE id = $1
		`,
		p.AuthorID,
	).Scan(&id)

	if err == pgx.ErrNoRows {
		err = tx.QueryRow(context.Background(), `
			INSERT INTO authors (id, name)
			VALUES ($1, $2);
			`,
			p.AuthorID,
			p.AuthorName,
		).Scan(&id)
	}

	return err
}

/*var posts = []storage.Post{
	{
		ID:      1,
		Title:   "Effective Go",
		Content: "Go is a new language. Although it borrows ideas from existing languages, it has unusual properties that make effective Go programs different in character from programs written in its relatives. A straightforward translation of a C++ or Java program into Go is unlikely to produce a satisfactory result—Java programs are written in Java, not Go. On the other hand, thinking about the problem from a Go perspective could produce a successful but quite different program. In other words, to write Go well, it's important to understand its properties and idioms. It's also important to know the established conventions for programming in Go, such as naming, formatting, program construction, and so on, so that programs you write will be easy for other Go programmers to understand.",
	},
	{
		ID:      2,
		Title:   "The Go Memory Model",
		Content: "The Go memory model specifies the conditions under which reads of a variable in one goroutine can be guaranteed to observe values produced by writes to the same variable in a different goroutine.",
	},
}*/
