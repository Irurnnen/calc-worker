package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Irurnnen/calc-worker/internal/models"
)

func Worker(orchestratorURL string) {
	// Send http request to orchestrator
	// TODO: replace hardcode to normal variable
	resp, err := FetchGET(context.Background(), orchestratorURL+"/internal/task", time.Second*3)
	if err != nil {
		log.Printf("Error while fetching orchestrator API")
	}

	// Process of status code
	switch resp.StatusCode {
	case http.StatusNotFound:
		// TODO: Add normal comments
		// TODO: Add normal sleep
		// TODO: replace hardcode to normal variable
		time.Sleep(time.Second * 5)
		break
	case http.StatusOK:
		// unmarshal json
		var task models.Task
		if err := json.Unmarshal([]byte(resp.Data), &task); err != nil {
			log.Printf("Error while parsing json: %s", err.Error())
			break
		}

		//  solve that request
		answer := Solver(task.FirstArgument, task.SecondArgument, task.Operation)

		// time.sleep for emulation of work
		time.Sleep(time.Second * time.Duration(task.OperationTime))

		// Send answer to orchestrator
		body, err := StructToReader(models.Answer{
			ID:     task.ID,
			Result: answer,
		})
		if err != nil {
			log.Printf("Error while converting struct to reader: %s", err.Error())
		}
		respPost, err := FetchPOST(context.Background(), orchestratorURL+"/internal/task", body, time.Second*3)
		if err != nil {
			log.Printf("Error while sending answer to orchestrator API")
		}

		switch respPost.StatusCode {
		case http.StatusOK:
			break
		case http.StatusNotFound:
			log.Printf("Not found task with id %s", task.ID)
			break
		case http.StatusUnprocessableEntity:
			log.Printf("Unprocessable Entity in task with id %s", task.ID)
			break
		case http.StatusInternalServerError:
			log.Printf("Unknown error while sending answer to task with id %s", task.ID)
		}
		break
	case http.StatusInternalServerError:
	default:
		log.Printf("Error while fetching orchestrator API: get %d code", resp.StatusCode)
		break
	}
}

func StructToReader(data any) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonData), nil
}

func FetchAPI(ctx context.Context, method, url string, body io.Reader, timeout time.Duration) (*APIResponse, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctxWithTimeout, method, url, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, context.DeadlineExceeded
		}
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &APIResponse{
		Data:       string(respBody),
		StatusCode: resp.StatusCode,
	}, nil
}

func FetchGET(ctx context.Context, url string, timeout time.Duration) (*APIResponse, error) {
	return FetchAPI(ctx, http.MethodGet, url, nil, timeout)
}

func FetchPOST(ctx context.Context, url string, body io.Reader, timeout time.Duration) (*APIResponse, error) {
	return FetchAPI(ctx, http.MethodPost, url, body, timeout)
}

type APIResponse struct {
	Data       string // тело ответа
	StatusCode int    // код ответа
}

func Solver(arg1, arg2 float64, operation string) float64 {
	switch operation {
	case "+":
		return arg1 + arg2
	case "-":
		return arg1 - arg2
	case "*":
		return arg1 * arg2
	case "/":
		return arg1 / arg2
	default:
		return 0
	}
}
