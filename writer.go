package logger

import "io"
import "fmt"
import "encoding/json"

type jsonToTextWriter struct {
	w io.Writer
}

// Write will convert a JSON message from slog to AWS Lambda default text format
func (e jsonToTextWriter) Write(p []byte) (int, error) {
	data := make(map[string]interface{})

	err := json.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}

	// requestID might be emtpy
	requestID := "-"

	if rID, has := data["requestId"]; has {
		requestID = rID.(string)
	}

	msg := fmt.Sprintf("%s\t%s\t%s\t%s", data["timestamp"], requestID, data["level"], data["message"])

	// remove data that has been added to the message...
	delete(data, "timestamp")
	delete(data, "level")
	delete(data, "requestId")
	delete(data, "message")

	// ...and add remaining data as JSON formatted string to the message
	if len(data) > 0 {
		remainingJSON, err := json.Marshal(data)
		if err != nil {
			return 0, err
		}

		msg = msg + " " + string(remainingJSON)
	}

	msg = msg + "\n"

	return e.w.Write([]byte(msg))
}
