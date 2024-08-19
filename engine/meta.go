package engine

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gruntwork-io/terragrunt-engine-go/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/structpb"
)

// Meta extract request.Meta to go map[string]interface{} struct.
func Meta(request *proto.RunRequest) (map[string]interface{}, error) {
	protoMeta := request.Meta
	meta := make(map[string]interface{})
	for key, anyValue := range protoMeta {
		var value structpb.Value
		if err := anyValue.UnmarshalTo(&value); err != nil {
			return nil, fmt.Errorf("error unmarshaling any value to structpb.Value: %w", err)
		}
		jsonData, err := value.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling structpb.Value to JSON: %w", err)
		}
		var result interface{}
		if err := json.Unmarshal(jsonData, &result); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON to interface{}: %w", err)
		}
		meta[key] = result
	}

	return meta, nil
}

// MetaString returns the value of a meta key from the RunRequest. If the key does not exist, an empty string is returned.
func MetaString(request *proto.RunRequest, key string) (string, error) {
	meta, err := Meta(request)
	log.Printf("Meta: %v", meta)
	if err != nil {
		return "", err
	}
	value, ok := meta[key]
	log.Printf("value: %v", value)
	if !ok {
		return "", nil
	}
	return strconv.Unquote(fmt.Sprintf("%v", value))
}
