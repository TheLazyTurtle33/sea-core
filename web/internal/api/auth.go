package api

import (
	"encoding/json"
	"os"
	"slices"

	"github.com/TheLazyTurtle33/sea-core/shared/logger"
)

const toeknPath = "/app/data/token.json"

func CheckAuth(token string) bool {

	content, err := os.ReadFile(toeknPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("auth: file doesnt exist, making new one", err)
			err = os.WriteFile(toeknPath, []byte(`{"tokens": []}`), 0644)
			if err != nil {
				logger.Error("auth: faild to make file", err)
			}
		} else {
			logger.Error("auth: faild to open file", err)
		}
		return false
	}

	var data map[string][]string
	err = json.Unmarshal(content, &data)
	if err != nil {
		logger.Error("auth: faild to Unmarshal json", err)
	}

	if data["tokens"] != nil || len(data["tokens"]) != 0 {
		tokens := data["tokens"]
		return slices.Contains(tokens, token)
	}

	return false
}
