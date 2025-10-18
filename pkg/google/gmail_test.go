package google

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

// mockRoundTripper intercepta requests y devuelve respuestas falsas
type mockRoundTripper struct {
	fn func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.fn(req), nil
}

// helper para crear un cliente falso
func newMockClient(fn func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{Transport: &mockRoundTripper{fn: fn}}
}

func TestDeleteWatch_Success(t *testing.T) {
	client := newMockClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("{}")),
			Header:     make(http.Header),
		}
	})

	g, _ := NewGmailClient(context.Background(), client)
	err := g.DeleteWatch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetMessageByID_Success(t *testing.T) {
	client := newMockClient(func(req *http.Request) *http.Response {
		body := `{"id":"abc123","internalDate":"1735689600000","payload":{"parts":[]}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})

	g, _ := NewGmailClient(context.Background(), client)
	msg, err := g.GetMessageByID(context.Background(), "abc123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Id != "abc123" {
		t.Errorf("expected abc123, got %s", msg.Id)
	}
}

func TestGetExtractMessages_Success(t *testing.T) {
	client := newMockClient(func(req *http.Request) *http.Response {
		body := `{"messages":[{"id":"m1"},{"id":"m2"}]}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})

	g, _ := NewGmailClient(context.Background(), client)
	res, err := g.GetExtractMessages(context.Background(), "davivienda")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(res.Messages))
	}
}

func TestGetMessageAttachment_Success(t *testing.T) {
	data := base64.URLEncoding.EncodeToString([]byte("hello world"))

	client := newMockClient(func(req *http.Request) *http.Response {
		body := `{"data":"` + data + `"}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})

	g, _ := NewGmailClient(context.Background(), client)
	att, err := g.GetMessageAttachment(context.Background(), "msg1", "att1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, _ := base64.URLEncoding.DecodeString(att.Data)
	if string(decoded) != "hello world" {
		t.Errorf("expected 'hello world', got %s", decoded)
	}
}

func TestDownloadAttachments_Success(t *testing.T) {
	// Attachment body
	data := base64.URLEncoding.EncodeToString([]byte("pdfdata"))

	// Step 1: intercept calls to Messages.Get
	calls := 0
	client := newMockClient(func(req *http.Request) *http.Response {
		if strings.Contains(req.URL.Path, "/messages/") && !strings.Contains(req.URL.Path, "/attachments/") {
			// Return message with attachment
			body := `{
				"id":"msg123",
				"internalDate":"1735689600000",
				"payload":{
					"parts":[
						{"filename":"extracto.pdf","body":{"attachmentId":"att123"}}
					]
				}
			}`
			calls++
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		}

		if strings.Contains(req.URL.Path, "/attachments/") {
			// Return attachment
			body := `{"data":"` + data + `"}`
			calls++
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}
		}

		return &http.Response{StatusCode: 404, Body: io.NopCloser(strings.NewReader(`{}`))}
	})

	g, _ := NewGmailClient(context.Background(), client)

	month, year, path, err := g.DownloadAttachments(context.Background(), "acct1", "msg123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if month == 0 || year == 0 {
		t.Errorf("expected valid month/year, got %d %d", month, year)
	}
	if !strings.HasSuffix(path, "extracto.pdf") {
		t.Errorf("expected file extracto.pdf, got %s", path)
	}

	// cleanup
	_ = os.Remove(path)
}

func TestDownloadAttachments_NoAttachment(t *testing.T) {
	client := newMockClient(func(req *http.Request) *http.Response {
		body := `{"id":"msg123","internalDate":"1735689600000","payload":{"parts":[]}}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}
	})

	g, _ := NewGmailClient(context.Background(), client)
	_, _, _, err := g.DownloadAttachments(context.Background(), "acct1", "msg123")

	if err == nil || !strings.Contains(err.Error(), "no attachment found") {
		t.Errorf("expected no attachment error, got %v", err)
	}
}
