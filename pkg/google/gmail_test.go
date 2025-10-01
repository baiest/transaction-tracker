package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func newTestGmailService(handler http.Handler) (*GmailService, error) {
	srv := httptest.NewServer(handler)

	gmailService, err := gmail.NewService(context.Background(), option.WithHTTPClient(srv.Client()), option.WithEndpoint(srv.URL))
	if err != nil {
		return nil, err
	}

	return &GmailService{Client: gmailService}, nil
}

func TestGetAttachmentData(t *testing.T) {
	c := require.New(t)

	mockAttachmentID := "test_attachment_id"
	mockMessageID := "test_message_id"
	mockData := "test-data"
	mockDataEncoded := base64.URLEncoding.EncodeToString([]byte(mockData))

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"data": "%s"}`, mockDataEncoded)
	})

	g, err := newTestGmailService(h)
	c.NoError(err)

	part := &gmail.MessagePart{
		Body: &gmail.MessagePartBody{
			AttachmentId: mockAttachmentID,
		},
	}

	data, err := g.getAttachmentData(context.Background(), mockMessageID, part)
	c.NoError(err)
	c.Equal([]byte(mockData), data)
}

func TestSaveAttachment(t *testing.T) {
	c := require.New(t)

	t.Run("saves a file successfully", func(t *testing.T) {
		accountID := "test_account"
		fileName := "test_file.txt"
		fileContent := []byte("test content")

		g, err := newTestGmailService(nil)
		c.NoError(err)

		filePath, err := g.saveAttachment(accountID, fileName, fileContent, 0600)
		c.NoError(err)

		// Verify file exists
		_, err = os.Stat(filePath)
		c.NoError(err)

		// Verify content
		content, err := os.ReadFile(filePath)
		c.NoError(err)
		c.Equal(fileContent, content)

		// Clean up
		c.NoError(os.RemoveAll(filepath.Dir(filepath.Dir(filePath))))
	})

	t.Run("it sanitizes the account id", func(t *testing.T) {
		accountID := "../../test_account"
		fileName := "test_file.txt"
		fileContent := []byte("test content")

		g, err := newTestGmailService(nil)
		c.NoError(err)

		filePath, err := g.saveAttachment(accountID, fileName, fileContent, 0600)
		c.NoError(err)

		c.NotContains(filePath, "..")

		// Clean up
		c.NoError(os.RemoveAll(filepath.Dir(filepath.Dir(filePath))))
	})
}

func TestDownloadAttachments(t *testing.T) {
	c := require.New(t)

	mockAttachmentID := "test_attachment_id"
	mockMessageID := "test_message_id"
	mockData := "test-data"
	mockDataEncoded := base64.URLEncoding.EncodeToString([]byte(mockData))
	mockFilename := "test_file.pdf"

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if strings.Contains(r.URL.Path, fmt.Sprintf("messages/%s/attachments/%s", mockMessageID, mockAttachmentID)) {
			fmt.Fprintf(w, `{"data": "%s"}`, mockDataEncoded)
			return
		}

		if strings.Contains(r.URL.Path, fmt.Sprintf("messages/%s", mockMessageID)) {
			fmt.Fprintf(w, `{
				"id": "%s",
				"internalDate": "%d",
				"payload": {
					"parts": [
						{
							"filename": "%s",
							"body": {
								"attachmentId": "%s"
							}
						}
					]
				}
			}`, mockMessageID, time.Now().UnixMilli(), mockFilename, mockAttachmentID)
			return
		}
	})

	g, err := newTestGmailService(h)
	c.NoError(err)

	month, year, path, err := g.DownloadAttachments(context.Background(), "test_account", mockMessageID)

	c.NoError(err)
	c.NotEmpty(path)
	c.Equal(time.Now().Month(), month)
	c.Equal(time.Now().Year(), year)

	// Clean up
	c.NoError(os.RemoveAll(filepath.Dir(filepath.Dir(path))))
}
