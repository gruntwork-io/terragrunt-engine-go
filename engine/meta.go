package engine

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gruntwork-io/terragrunt-engine-go/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// Meta extracts request.Meta to go map[string]any struct.
func Meta(request *proto.RunRequest) (map[string]any, error) {
	protoMeta := request.GetMeta()
	meta := make(map[string]any)

	for key, anyValue := range protoMeta {
		var value structpb.Value
		if err := anyValue.UnmarshalTo(&value); err != nil {
			return nil, fmt.Errorf("error unmarshaling any value to structpb.Value: %w", err)
		}

		jsonData, err := value.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("error marshaling structpb.Value to JSON: %w", err)
		}

		var result any
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
	if err != nil {
		return "", err
	}

	value, ok := meta[key]
	if !ok {
		return "", nil
	}

	return strconv.Unquote(fmt.Sprintf("%v", value))
}
