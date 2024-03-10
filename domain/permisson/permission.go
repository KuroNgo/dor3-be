package permisson

import "go.mongodb.org/mongo-driver/bson/primitive"

type Permission struct {
	PermissionID primitive.ObjectID `bson:"_id" json:"_id"`
	Name         string             `bson:"name" json:"name"`
	Method       string             `bson:"method" json:"method"`
	RoutePath    string             `bson:"route_path" json:"route_path"`
}

type PermissionInput struct {
	Name      string `bson:"name" json:"name"`
	Method    string `bson:"method" json:"method"`
	RoutePath string `bson:"route_path" json:"route_path"`
}
