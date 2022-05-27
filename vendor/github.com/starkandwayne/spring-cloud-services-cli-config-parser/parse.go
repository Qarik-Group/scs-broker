package scsccparser

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type EnvironmentSetup struct {
	Config map[string]string
}

const baseConfig string = "SPRING_CLOUD_CONFIG_SERVER_"

func (config *EnvironmentSetup) ParseEnvironmentFromString(configline string) (map[string]string, error) {

	log.Print("Received Raw Config: " + configline)

	var result map[string]interface{}

	err := json.Unmarshal([]byte(configline), &result)

	if err != nil {
		return nil, err
	}

	str, err := buildKey(true, baseConfig, result)
	if err != nil {
		return nil, err
	}

	return str, nil
}

func (config *EnvironmentSetup) ParseEnvironmentFromRaw(configraw json.RawMessage) (map[string]string, error) {

	var result map[string]interface{}

	err := json.Unmarshal(configraw, &result)

	if err != nil {
		return nil, err
	}

	str, err := buildKey(true, baseConfig, result)
	if err != nil {
		return nil, err
	}

	return str, nil
}

func buildKey(root bool, base string, keyint map[string]interface{}) (map[string]string, error) {

	envmap := make(map[string]string)

	var sb strings.Builder

	keys := keyint

	for key, value := range keys {

		if !root {
			sb.WriteString(base + key)
		} else {
			sb.Reset()
			sb.WriteString(baseConfig + key)
		}

		switch t := value.(type) {
		default:
			log.Panicf("Invalid Type! %T", t)
		case map[string]interface{}:
			sb.WriteString("_")
			str, err := buildKey(false, sb.String(), value.(map[string]interface{}))
			if err != nil {
				log.Println(err)
			}
			for k, v := range str {
				envmap[k] = v
			}
		case bool:
			s := strconv.FormatBool(value.(bool))
			envmap[strings.ToUpper(sb.String())] = s
			sb.Reset()
		case float64:
			s := fmt.Sprintf("%g", value.(float64))
			envmap[strings.ToUpper(sb.String())] = s
			sb.Reset()
		case string:
			envmap[strings.ToUpper(sb.String())] = value.(string)
			sb.Reset()
		}

	}
	sb.Reset()
	return envmap, nil
}
