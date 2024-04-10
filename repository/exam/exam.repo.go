package exam_repository

import "clean-architecture/infrastructor/mongo"

type examRepository struct {
	database mongo.Database
}
