package docgen_test

import (
	"os"
	"testing"

	docgen "github.com/JessBellon/apitest-docgen"
	"github.com/steinfletcher/apitest"
)

var r *docgen.Reporter

func TestMain(m *testing.M) {

	r = docgen.NewReporter()

	code := m.Run()

	r.GenDocs()
	os.Exit(code)
}
func TestOne(t *testing.T) {
	apitest.New("Test Subtitle 1").
		Report(r).
		EnableNetworking().
		Post("https://postman-echo.com/post").
		Body(`{"howdy": "yall"}`).
		Header("Content-Type", "application/json").
		Expect(t).
		End()
}

func TestTwo(t *testing.T) {
	r := docgen.NewReporter()
	defer r.GenDocs()
	apitest.New("Test Subtitle 2").
		Report(r).
		EnableNetworking().
		Post("https://postman-echo.com/post").
		Query("x", "y").
		Header("Content-Type", "application/json").
		Expect(t).
		End()
}

func TestThree(t *testing.T) {
	apitest.New("Test Subtitle 3").
		Report(r).
		EnableNetworking().
		Get("https://postman-echo.com/post").
		Header("Content-Type", "application/json").
		Expect(t).
		End()
}

// func TestGenDocs(t *testing.T) {
// 	GenDocs()
// }
