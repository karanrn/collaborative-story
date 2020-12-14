package story

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"CollaborativeStory/colab/models"
)

func TestGetStories(t *testing.T) {
	mockDB := StoryDBMock{}

	testStories := []models.Story{
		{
			ID:        1,
			Title:     "Hello World!",
			CreatedAt: "2020-12-08T12:13:42Z",
			UpdatedAt: "2020-12-08T13:13:42Z",
		},
	}

	mockDB.On("FetchStories", "created_at", "asc", int64(0), int64(100)).Return(testStories, nil)
	mockDB.On("FetchStories", "created_at", "asc", int64(0), int64(10)).Return(testStories, nil)

	csMock := ColabStory{Database: &mockDB}

	var req *http.Request
	var err error
	var rr *httptest.ResponseRecorder

	// Table driven testing
	testData := []struct {
		url          string
		expectedCode int
		expectedData string
	}{
		{"/stories", http.StatusOK, `"{'limit': 100, 'offset': 0, 'count': 1, 'results': [{\"ID\":1,\"Title\":\"Hello World!\",\"created_at\":\"2020-12-08T12:13:42Z\",\"updated_at\":\"2020-12-08T13:13:42Z\"}] }"` + "\n"},
		{"/stories?limit=str", http.StatusBadRequest, `"{'error': 'limit is not an integer'}"` + "\n"},
		{"/stories?limit=10", http.StatusOK, `"{'limit': 10, 'offset': 0, 'count': 1, 'results': [{\"ID\":1,\"Title\":\"Hello World!\",\"created_at\":\"2020-12-08T12:13:42Z\",\"updated_at\":\"2020-12-08T13:13:42Z\"}] }"` + "\n"},
		{"/stories?offset=str", http.StatusBadRequest, `"{'error': 'offset is not an integer'}"` + "\n"},
		{"/stories?sort=hello", http.StatusBadRequest, `"{'error': 'sort should be among these values [title created_at updated_at]'}"` + "\n"},
		{"/stories?order=kkk", http.StatusBadRequest, `"{'error': 'order should be among these values [asc desc]'}"` + "\n"},
	}

	handler := http.HandlerFunc(csMock.GetStories())

	for _, tt := range testData {
		req, err = http.NewRequest("GET", tt.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, tt.expectedCode)
		assert.Equal(t, rr.Body.String(), tt.expectedData)
	}
}
