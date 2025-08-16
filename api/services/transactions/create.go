package transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"transaction-tracker/api/models"
	"transaction-tracker/shared"
)

var (
	createTransactionURL = os.Getenv("CREATE_TRANSACTION_URL")
)

func Create(tr *models.TransactionRequest) error {
	bodyBytes, err := json.Marshal(tr)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(string(models.POST), createTransactionURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := shared.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var responseBody bytes.Buffer

	_, err = responseBody.ReadFrom(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("Error unexpected status: %d request: %s, message: %s", res.StatusCode, string(bodyBytes), responseBody.String())
	}

	return nil
}
