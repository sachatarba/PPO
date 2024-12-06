package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type ICodeRepository interface {
	SaveCode(ctx context.Context, code string, clienID uuid.UUID) error
	GetCodeByClientID(ctx context.Context, clientID uuid.UUID) (string, error)
}

type CodeRepo struct {
	client *redis.Client
}

func NewCodeRepo(client *redis.Client) ICodeRepository {
	return &CodeRepo{client: client}
}

func (r *CodeRepo) SaveCode(ctx context.Context, code string, clienID uuid.UUID) error {
	data, err := json.Marshal(code)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, clienID.String(), data, time.Minute*10).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *CodeRepo) GetCodeByClientID(ctx context.Context, clientID uuid.UUID) (string, error) {
	code := r.client.Get(ctx, clientID.String())
	if code == nil {
		return "", fmt.Errorf("code not found")
	}

	// var str struct {
	// 	Code string `json:"code"`
	// }

	var str string

	bytes, err := code.Bytes()
	if err != nil {
		return "", fmt.Errorf("can't get bytes")
	}
	err = json.Unmarshal(bytes, &str)
	if err != nil {
		return "", fmt.Errorf("can't unmarshal")
	}

	log.Println("code unmarshal:", str)

	return str, nil
}
