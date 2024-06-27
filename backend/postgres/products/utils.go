package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getMongoCollection(srv *Service, collectioName string) *mongo.Collection {
	return srv.MongoDB.Database(srv.Config.MONGO_DATABASE).Collection(collectioName)
}

func getMongoRecordsCount(collection *mongo.Collection, filter interface{}) (int64, error) {
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getMongoRecord(collection *mongo.Collection, filter interface{}) (map[string]interface{}, error) {
	var record map[string]interface{}
	err := collection.FindOne(context.TODO(), filter).Decode(&record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func getMongoRecords(collection *mongo.Collection, filter interface{}) ([]map[string]interface{}, error) {
	var records []map[string]interface{}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var record map[string]interface{}
		err := cursor.Decode(&record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func prepareMongoFilters(filters []map[string]interface{}) bson.M {
	if filters == nil {
		return bson.M{}
	}
	var andConditions []bson.M
	for _, filter := range filters {
		key := filter["key"].(string)
		value := filter["value"]

		var orConditions []bson.M

		switch v := value.(type) {
		case string:
			if key == "id" || key == "_id" {
				temp, _ := primitive.ObjectIDFromHex(value.(string))
				andConditions = append(andConditions, bson.M{key: bson.M{"$eq": temp}})
			} else {
				andConditions = append(andConditions, bson.M{key: bson.M{"$eq": v}})
			}
			// orConditions = append(orConditions, bson.M{key: bson.M{"$eq": v}})
		case []string:
			for _, val := range v {
				orConditions = append(orConditions, bson.M{key: bson.M{"$eq": val}})
			}
		}

		if len(orConditions) > 0 {
			andConditions = append(andConditions, bson.M{"$or": orConditions})
		}
	}

	if len(andConditions) > 0 {
		return bson.M{"$and": andConditions}
	}
	return bson.M{}
}
