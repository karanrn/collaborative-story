package story

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"CollaborativeStory/colab/models"

	"github.com/stretchr/testify/assert"
)

func TestPostStory(t *testing.T) {

	mockDB := StoryDBMock{}
	mockDB.On("GetLatestStory").Return(models.Story{ID: 0}, nil)
	mockDB.On("AddStory", 1, "Hello", true).Return(nil)

	csMock := ColabStory{Database: &mockDB}

	var req *http.Request
	var err error
	var rr *httptest.ResponseRecorder

	testData := []struct {
		request      []byte
		expectedCode int
		expectedData string
	}{
		{[]byte(`{ "word": "Hello" }`), http.StatusOK, `"{\"id\":1,\"title\":\"Hello\",\"current_sentence\":\"\"}"` + "\n"},
		{[]byte(`{ "word": "Hello World" }`), http.StatusBadRequest, `"{'error': 'multiple words sent'}"` + "\n"},
	}

	handler := http.HandlerFunc(csMock.PostStory())

	for _, tt := range testData {
		req, err = http.NewRequest("POST", "/add", bytes.NewBuffer(tt.request))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		rr = httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, tt.expectedCode)
		assert.Equal(t, rr.Body.String(), tt.expectedData)
	}
}
