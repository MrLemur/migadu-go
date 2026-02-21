package migadu

import (
	"encoding/json"
	"io"
)

// readAndUnmarshal reads the response body and unmarshals it into the target.
func readAndUnmarshal(reader io.Reader, target interface{}) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
