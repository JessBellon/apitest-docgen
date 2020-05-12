package docgen

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/steinfletcher/apitest"
)

// Record represents one request-response event.
type Record struct {
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

type Request struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Headers http.Header `json:"headers,omitempty"`
	Body    string      `json:"body,omitempty"`
}
type Response struct {
	Status  string      `json:"status"`
	Headers http.Header `json:"headers,omitempty"`
	Body    string      `json:"body,omitempty"`
}

var headerWhitelist = map[string]struct{}{
	"Content-Type": struct{}{},
}

// for version 0.0.1 each doc that gets generated will be restricted to a single instance of a Reporter
// I have MANY thoughts on how to improve this, but I need a MVP first
type Reporter struct {
	Records map[string][]Record
}

func NewReporter() *Reporter {
	return &Reporter{
		Records: make(map[string][]Record),
	}
}

// Format makes a lot of assumptions but again going for a mvp

func (r *Reporter) Format(recorder *apitest.Recorder) {
	var rec Record
	for _, event := range recorder.Events {
		switch v := event.(type) {
		case apitest.HttpRequest:

			rec.Request = copyRequest(v.Value)
			rec.Title = rec.Request.Method + " " + rec.Request.Path // I dont like the default title
			rec.Subtitle = recorder.SubTitle

		case apitest.HttpResponse:
			rec.Response = copyResponse(v.Value)

			d, ok := r.Records[rec.Title]
			if !ok {
				d = make([]Record, 0, 1)
			}

			r.Records[rec.Title] = append(d, rec)
			rec = Record{}

		case apitest.MessageRequest:
			fmt.Printf("Message Request: %+v\n", v)
		case apitest.MessageResponse:
			fmt.Printf("Message Response: %+v\n", v)
		default:
			panic("received unknown event type")
		}
	}
}

func copyRequest(req *http.Request) *Request {
	var b bytes.Buffer

	req.Body = copyBody(req.Body, &b)

	return &Request{
		Method:  req.Method,
		Path:    req.URL.Path,
		Headers: req.Header,
		Body:    b.String(),
	}

}

func copyResponse(resp *http.Response) *Response {

	var b bytes.Buffer
	resp.Body = copyBody(resp.Body, &b)

	return &Response{
		Status:  resp.Status,
		Headers: resp.Header,
		Body:    b.String(),
	}

}

// CopyBody takes an io.ReadCloser, reads it in to a bytes.Buffer and returns a copy
func copyBody(in io.ReadCloser, cpy *bytes.Buffer) io.ReadCloser {
	if _, err := cpy.ReadFrom(in); err != nil {
		panic(err)
	}

	in.Close()

	return ioutil.NopCloser(bytes.NewBuffer(cpy.Bytes()))
}

func (r *Reporter) GenDocs() error {

	filename := "./DOCS.txt"

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	for key := range r.Records {
		rec := r.Records[key]
		// formatRecords(rec)
		// fmt.Println(i)
		for i := range rec {
			x := rec[i]

			fmt.Fprintf(w, "Title: %s\n", x.Title)
			fmt.Fprintf(w, "Subtitle: %s\n", x.Subtitle)
			fmt.Fprintf(w, "Request: %s\n", x.Request)
			fmt.Fprint(w, "\n")
			fmt.Fprintf(w, "Response: %s\n", x.Response)
			fmt.Fprint(w, "\n")
			fmt.Fprint(w, "-----------\n")
		}
		// fmt.Println(rec.Title)
		// fmt.Printf("\t%+v\n", rec)
		// data, err := json.Marshal(&rec)
		// if err != nil {
		// 	panic(err)
		// }
		// if _, err := w.Write(append(data, []byte("\n")...)); err != nil {
		// 	panic(err)
		// }
	}
	return nil
}

func (r *Request) String() string {
	var sb strings.Builder
	sb.WriteString("    " + r.Method + "   " + r.Path + "\n")
	if len(r.Headers) > 0 {
		for key, val := range r.Headers {
			sb.WriteString("Headers:    " + key + ":")
			for i := range val {
				sb.WriteString(val[i])
			}
			sb.WriteString("\n")
		}
	}
	if len(r.Body) > 0 {
		sb.WriteString("Body:	" + r.Body)
	}
	return sb.String()
}

func (r *Response) String() string {
	var sb strings.Builder
	sb.WriteString("    " + r.Status + "\n")
	if len(r.Headers) > 0 {
		for key, val := range r.Headers {
			if _, ok := headerWhitelist[key]; !ok {
				continue
			}
			sb.WriteString("Headers:    " + key + ":")
			for i := range val {

				sb.WriteString(val[i])
			}
			sb.WriteString("\n")
		}
	}
	if len(r.Body) > 0 {
		sb.WriteString("Body:	")
		sb.WriteString(r.Body)
	}
	return sb.String()
}
