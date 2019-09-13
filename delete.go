package main

import (
	"log"
)

func deleteAlert() (*map[string]interface{}, error) {
	var result map[string]interface{}
	table, err := getAlertTable()
	if err != nil {
		return &result, err
	}
	hashKey, err := getAlertHashKey()
	if err != nil {
		return &result, err
	}

	var records []alertTable
	if err := table.Scan().All(&records); err != nil {
		return &result, err
	}
	for _, record := range records {
		if err := table.Delete(hashKey, record.Timestamp).Run(); err != nil {
			return &result, err
		}
	}
	log.Printf("[DELETE]%d alert deleted.", len(records))
	result = map[string]interface{}{
		"DeleteCount": len(records),
	}
	return &result, nil
}
