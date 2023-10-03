package currencyapi

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

func ConvertCurrencies(from string, to string) (float64, error) {
	apiKey := viper.GetString("API_KEY")

	requestURL := fmt.Sprintf("https://api.freecurrencyapi.com/v1/latest?apikey=%s&currencies=%s&base_currency=%s",
		apiKey,
		to,
		from,
	)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return 0, fmt.Errorf("error while trying to create currency api request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("error while trying to get currency response: %v", err)
	}

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error while trying read response body: %v", err)
	}

	var responseMap map[string]map[string]float64

	err = json.Unmarshal(result, &responseMap)
	if err != nil {
		return 0, fmt.Errorf("error while trying to unmarshal response body: %v", err)
	}

	return responseMap["data"][to], nil
}
