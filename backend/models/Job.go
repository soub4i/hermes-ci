package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Job struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Workflow     string             `json:"workflow" bson:"workflow"`
	Repository   string             `json:"repository" bson:"repository"`
	RepositoryId string             `json:"repository_id" bson:"repository_id"`
	Logs         []string           `json:"logs,omitempty" bson:"logs,omitempty"`
	Finished     bool               `json:"finished" bson:"finished"`
	Owner        OwnerType          `bson:"owner,omitempty" json:"owner,omitempty"`
	CreateAt     time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdateAt     time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type OwnerType struct {
	Name   string `bson:"name,omitempty" json:"name,omitempty"`
	Login  string `bson:"login,omitempty" json:"login,omitempty"`
	Avatar string `bson:"avatar,omitempty" json:"avatar,omitempty"`
}
