package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/go-redis/redis/v8"
)

type RedisDatabase struct {
	client *redis.Client
}

func NewRedisDatabase() *RedisDatabase {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisDatabase{
		client: client,
	}
}

func (db *RedisDatabase) Save(shortURL, originalURL string) error {
	ctx := context.Background()
	return db.client.Set(ctx, shortURL, originalURL, 0).Err()
}

func (db *RedisDatabase) Get(shortURL string) (string, error) {
	ctx := context.Background()
	return db.client.Get(ctx, shortURL).Result()
}

func (db *RedisDatabase) Close() error {
	return db.client.Close()
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase() *PostgresDatabase {
	db, err := sql.Open("postgres", "postgres://postgres:@localhost/urls?sslmode=disable")
	if err != nil {
		panic(err)
	}

	return &PostgresDatabase{
		db: db,
	}
}

func (db *PostgresDatabase) Save(shortURL, originalURL string) error {
	_, err := db.db.Exec("INSERT INTO urls (short_url, original_url) VALUES ($1, $2)", shortURL, originalURL)
	return err
}

func (db *PostgresDatabase) Get(shortURL string) (string, error) {
	var originalURL string
	err := db.db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("Short URL not found")
		}
		return "", err
	}
	return originalURL, nil
}

func (db *PostgresDatabase) Close() error {
	return db.db.Close()
}
