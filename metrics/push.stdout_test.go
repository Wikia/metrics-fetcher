package metrics

import (
	"testing"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

func TestOutputMetricsWhenFiltersEmpty(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	//when
	OutputMetrics([]models.FilteredMetrics{}, testObj);
	//then
	testObj.AssertNotCalled(t, "Write")
}

func TestShouldWriteToOutputWhenNotEmptyFilters(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources,exampleTag=exampleTagValue exampleField=exampleFieldValue\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := models.NewFilteredMetric()
	metric.Tags["exampleTag"] = "exampleTagValue"
	metric.Fields["exampleField"] = "exampleFieldValue"
	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
}

func TestShouldErrorWhenNoFields(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	metric := models.NewFilteredMetric()
	metric.Tags["exampleTag"] = "exampleTagValue"
	//when
	err := OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	assert.NotNil(t, err)
}

func TestShouldOutputMultipleMetrics(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	metric1 := models.NewFilteredMetric()
	metric1.Tags["exampleTag1"] = "exampleTagValue1"
	metric1.Fields["exampleField1"] = "exampleFieldValue1"

	metric2 := models.NewFilteredMetric()
	metric2.Tags["exampleTag2"] = "exampleTagValue2"
	metric2.Fields["exampleField2"] = "exampleFieldValue2"

	expectedOutputLine1 := []byte("resources,exampleTag1=exampleTagValue1 exampleField1=exampleFieldValue1\n")
	expectedOutputLine2 := []byte("resources,exampleTag2=exampleTagValue2 exampleField2=exampleFieldValue2\n")
	testObj.On("Write", expectedOutputLine1).Return(len(expectedOutputLine1), nil).On("Write", expectedOutputLine2).Return(len(expectedOutputLine2), nil)
	//when
	OutputMetrics([]models.FilteredMetrics{metric1, metric2}, testObj)

	//then
	testObj.AssertExpectations(t)
}

type MockedWriter struct {
	mock.Mock
}

func (m *MockedWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}