package refresco_test

import (
	"encoding/json"
	"github.com/The-New-Fork/api-pipeline/pkg/pipeline"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/require"
	"github.com/unchainio/pkg/xconfig"
	"github.com/unchainio/pkg/xlogger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestRefrescoEndToEndSuccess(t *testing.T) {
	skipCI(t)

	// TEST SETUP
	// set up mock import-api
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		resBody, err := ioutil.ReadAll(request.Body)
		require.NoError(t, err, "request body could not be read")
		err = request.Body.Close()
		require.NoError(t, err, "request body could not be closed")

		expectedBody := `{
  "anfp": "18505100",
  "dfp": 100EP PA Apfelsaft naturtr�b NFC,
  "bnfp": "0000878937",
  "pds": 2020-01-20,
  "pde": 2020-01-21,
  "jds": 020,
  "jde": 021,
  "bbd": 2020-10-20,
  "pc": "Germany",
  "pl": "Herrath",
  "rmn": "11200100520",
  "pon": "4500885082",
  "pop": "10",
}
`
		require.Equal(t, expectedBody, string(resBody), "http request body received from pipeline on import api does not match expected body")
		render.JSON(writer, request, map[string]interface{}{
			"key": "value",
		})
	})
	go http.ListenAndServe(":8002", nil)
	log.Println("test setup complete")

	// TEST EXECUTION
	// load config & logger
	cfg := loadConfig()
	log, _ := xlogger.New(cfg.Logger)

	// create and start pipeline
	p := pipeline.New(cfg, log)

	go p.Start()
	//require.NoError(t, err)

	// send API request to trigger pipeline
	file, err := os.Open("./example.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	resp, err := http.Post("http://localhost:8888", "binary/octet-stream", file)
	require.NoError(t, err, "could not make request to pipeline")
	defer resp.Body.Close()

	// TEST ASSERTIONS
	// No error and received 200 response from pipeline
	require.NoError(t, err, "could not read returned body")
	require.Equal(t, resp.StatusCode, http.StatusOK, "response to pipeline should have status 200" )

	// check for no failed messages
	errorCount := countErrorEmails(t)
	require.Equal(t, float64(0), errorCount, "error alert email expected to be 0")

	// p.Stop()
}

func loadConfig() *pipeline.Config {
	// load config
	cfg := new(pipeline.Config)
	info := new(xconfig.Info)

	errs := xconfig.Load(
		cfg,
		xconfig.FromPathFlag("cfg", "./config.toml"),
		xconfig.FromEnv(),
		xconfig.GetInfo(info),
	)
	if errs != nil {
		log.Fatal(errs)
	}
	return cfg
}

func countErrorEmails(t *testing.T) float64 {
	resp, err := http.Get("http://localhost:8025/api/v2/messages")
	require.NoError(t, err, "could not get message")
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err, "could not read smtp server API response body")

	var objectMap map[string]interface{}
	err = json.Unmarshal(bytes, &objectMap)
	require.NoError(t, err, "could not unmarshal local smtp server API response")

	return objectMap["count"].(float64)
}
