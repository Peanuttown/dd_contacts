package contacts_callback

import (
	"context"
	"encoding/json"
)

type EventHandler func(ctx context.Context, req map[string]interface{}) error

func unMarshalFromMap(source map[string]interface{}, toObj interface{}) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, toObj)
	return err
}


