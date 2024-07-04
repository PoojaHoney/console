package main

import "github.com/google/uuid"

func preparePostgresFilters(filters []map[string]interface{}) map[string]interface{} {
	if filters == nil {
		return nil
	}
	allFilters := map[string]interface{}{}
	andConditions := map[string]interface{}{}
	for _, filter := range filters {
		key := filter["key"].(string)
		value := filter["value"]

		orConditions := map[string]interface{}{}

		switch v := value.(type) {
		case string:
			if key == "id" || key == "_id" {
				temp, _ := uuid.Parse(value.(string))
				andConditions[key] = temp
			} else {
				andConditions[key] = v
			}
			// orConditions = append(orConditions, bson.M{key: bson.M{"$eq": v}})
		case []string:
			for _, val := range v {
				orConditions[key] = val
			}
		}

		if len(orConditions) > 0 || orConditions != nil {
			allFilters["or"] = orConditions
		}
	}

	if len(andConditions) > 0 {
		allFilters["and"] = andConditions
		return allFilters
	}
	return nil
}