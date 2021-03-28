package testdata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/labstack/echo/v4"
)

// basepath is the root directory of this package
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path return the absolute path to a file
// if path is already absolute from root, then it will return immediately
// but if not, it will concat base path and relative path from this root project folder
func Path(relPath string) string {
	if filepath.IsAbs(relPath) {
		return relPath
	}

	return filepath.Join(basepath, relPath)
}

// GetGolden is a function to get golden file
func GetGolden(t *testing.T, filename string) []byte {
	t.Helper()

	b, err := ioutil.ReadFile(Path(filename + ".golden"))
	if err != nil {
		t.Fatal(err)
	}

	return b
}

// UnmarshallGoldenToJSON is a function to unmarshall golden file to json
func UnmarshallGoldenToJSON(t *testing.T, filename string, input interface{}) {
	err := json.Unmarshal(GetGolden(t, filename), &input)
	if err != nil {
		t.Fatal(err)
	}
}

// HTTPCall is a form of request to mock http server with expected response
type HTTPCall struct {
	Header       map[string]string
	Method       string
	Status       int
	ExpectedResp []byte
}

// FuncCall is a form of function call with input and expected output
type FuncCall struct {
	Called bool
	Input  []interface{}
	Output []interface{}
}

func GetJSON(t *testing.T, filename string) []byte {
	t.Helper()

	b, err := ioutil.ReadFile(Path(filename + ".json"))
	if err != nil {
		t.Fatal(err)
	}

	return b
}

// MockServer is a mock server for HTTP Call testing
func MockServer(t *testing.T, reqs map[string]HTTPCall) (*httptest.Server, func()) {
	t.Helper()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		choosenRoute := fmt.Sprintf("%s %s", r.Method, r.RequestURI)
		req, ok := reqs[choosenRoute]
		if ok {
			for k, v := range req.Header {
				w.Header().Set(k, v)
			}
			w.WriteHeader(req.Status)

			_, err := w.Write(req.ExpectedResp)
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	return server, func() { server.Close() }
}

func GetEchoServer() *echo.Echo {
	e := echo.New()
	return e
}
