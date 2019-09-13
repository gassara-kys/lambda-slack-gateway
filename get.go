package main

func getAlertCount() (*map[string]interface{}, error) {
	var result map[string]interface{}
	table, err := getAlertTable()
	if err != nil {
		return &result, err
	}

	var results []alertTable
	if err := table.Scan().All(&results); err != nil {
		return &result, err
	}
	result = map[string]interface{}{
		"AlertCount": len(results),
	}
	return &result, nil
}
