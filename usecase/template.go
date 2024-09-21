package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

func (uc *useCase) renderTemplate(msg string, args any) (result string, err error) {
	var data map[string]interface{}
	if err = json.Unmarshal([]byte(msg), &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal template: %w", err)
	}

	if err = uc.processMap(data, args); err != nil {
		return "", fmt.Errorf("failed to process template: %w", err)
	}

	resultBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}
	result = string(resultBytes)

	return result, nil
}

func (uc *useCase) processMap(data map[string]interface{}, args any) error {
	for key, value := range data {
		switch v := value.(type) {
		case string:
			tmpl, err := template.New(key).Parse(v)
			if err != nil {
				return err
			}
			var tpl bytes.Buffer
			if err = tmpl.Execute(&tpl, args); err != nil {
				return err
			}
			data[key] = tpl.String()

		case map[string]interface{}:
			if err := uc.processMap(v, args); err != nil {
				return err
			}

		case []interface{}:
			for i, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if err := uc.processMap(itemMap, args); err != nil {
						return err
					}
					v[i] = itemMap
				}
			}
		}
	}
	return nil
}
