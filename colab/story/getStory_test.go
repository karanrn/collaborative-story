package story

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"CollaborativeStory/colab/models"
)

type StoryDBMock struct {
	mock.Mock
}

func (d *StoryDBMock) FetchStory(storyID string) (models.DetailedStory, error) {
	args := d.Called(storyID)
	return args.Get(0).(models.DetailedStory), args.Error(1)
}

func (d *StoryDBMock) FetchParagraphs(start int, end int, isComplete bool) ([]models.Paragraph, error) {
	args := d.Called(start, end, isComplete)
	return args.Get(0).([]models.Paragraph), args.Error(1)
}

func (d *StoryDBMock) FetchSentences(start int, end int, isComplete bool) ([]string, error) {
	args := d.Called(start, end, isComplete)
	return args.Get(0).([]string), args.Error(1)
}

func (d *StoryDBMock) AddStory(storyID int, word string, isNew bool) error {
	args := d.Called(storyID, word, isNew)
	return args.Error(0)
}

func (d *StoryDBMock) AddToSentence(word string) (sentenceID int, err error) {
	args := d.Called(word)
	return args.Int(0), args.Error(1)
}

func (d *StoryDBMock) AddToParagraph(sentenceID int) (paragraphID int, err error) {
	args := d.Called(sentenceID)
	return args.Int(0), args.Error(1)
}

func (d *StoryDBMock) FetchStories(sort string, order string, offset int64, limit int64) ([]models.Story, error) {
	args := d.Called(sort, order, offset, limit)
	return args.Get(0).([]models.Story), args.Error(1)
}

func (d *StoryDBMock) GetLatestStory() (models.Story, error) {
	args := d.Called()
	return args.Get(0).(models.Story), args.Error(1)
}

func (d *StoryDBMock) UpdateStoryParagraph(storyID int, paragraphID int, isEnd bool) error {
	args := d.Called(storyID, paragraphID, isEnd)
	return args.Error(0)
}

func (d *StoryDBMock) UpdateStoryTimestamp(storyID int) error {
	args := d.Called(storyID)
	return args.Error(0)
}

// Test

func TestGetStory(t *testing.T) {
	mockDB := StoryDBMock{}
	sentences := [][]string{{"hello", "world!", "welcome", "john"},
		{"world", "is", "beautiful,", "john"}}

	paragraphs := []models.Paragraph{{
		ID:            1,
		StartSentence: 1,
		EndSentence:   10,
		Sentences:     sentences[0],
	},
		{
			ID:            2,
			StartSentence: 11,
			EndSentence:   13,
			Sentences:     sentences[1],
		},
	}
	testStory := models.DetailedStory{
		ID:             1,
		Title:          "hello world",
		StartParagraph: 1,
		EndParagraph:   7,
		CreatedAt:      "2020-12-08T12:13:42Z",
		UpdatedAt:      "2020-12-08T13:13:42Z",
	}

	// Mock FetchSentences
	mockDB.On("FetchSentences", 1, 10, true).Return(sentences[0], nil)
	mockDB.On("FetchSentences", 11, 13, true).Return(sentences[1], nil)

	// Mock FetchParagraph
	mockDB.On("FetchParagraphs", 1, 7, true).Return(paragraphs, nil)

	// Mock FetchStory
	mockDB.On("FetchStory", "1").Return(testStory, nil)
	mockDB.On("FetchStory", "2").Return(models.DetailedStory{}, nil)

	csMock := ColabStory{Database: &mockDB}

	var req *http.Request
	var err error
	var rr *httptest.ResponseRecorder

	r := mux.NewRouter()
	r.HandleFunc("/stories/{id:[0-9]+}", csMock.GetStory())

	testData := []struct {
		url          string
		expectedCode int
		expectedData string
	}{
		{"/stories/1", http.StatusOK, `"{\"id\":1,\"title\":\"hello world\",\"created_at\":\"2020-12-08T12:13:42Z\",\"updated_at\":\"2020-12-08T13:13:42Z\",\"paragraphs\":[{\"sentences\":[\"hello\",\"world!\",\"welcome\",\"john\"]},{\"sentences\":[\"world\",\"is\",\"beautiful,\",\"john\"]}]}"` + "\n"},
		{"/stories/2", http.StatusNotFound, `"{'error': 'story does not exist'}"` + "\n"},
	}

	for _, tt := range testData {
		req, err = http.NewRequest("GET", tt.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr = httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, tt.expectedCode)
		assert.Equal(t, rr.Body.String(), tt.expectedData)
	}
}
