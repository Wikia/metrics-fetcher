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
	expectedOutputLine := []byte("resources,exampleTag=exampleTagValue exampleField=\"exampleFieldValue\"\n")
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

func TestShouldHandleNoTags(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources exampleField=\"exampleFieldValue\"\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := models.NewFilteredMetric()
	metric.Fields["exampleField"] = "exampleFieldValue"

	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)

	//then
	testObj.AssertExpectations(t)
}

func TestShouldHandleDifferentFieldTypes(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources stringField=\"exampleFieldValue\",intField=1,floatField=1.21,booleanField=true\n")
	testObj.On("Write", mock.Anything).Return(len(expectedOutputLine), nil)
	metric := models.NewFilteredMetric()
	metric.Fields["stringField"] = "example Field\\Value"
	metric.Fields["intField"] = 1
	metric.Fields["floatField"] = 1.21
	metric.Fields["booleanField"] = true

	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)

	//then
	testObj.AssertExpectations(t)
	calledWithMessage := string(testObj.Calls[0].Arguments.Get(0).([]byte))
	assert.Contains(t, calledWithMessage, "stringField=\"example Field\\\\Value\"")
	assert.Contains(t, calledWithMessage, "intField=1")
	assert.Contains(t, calledWithMessage, "floatField=1.21")
	assert.Contains(t, calledWithMessage, "booleanField=true")
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

	expectedOutputLine1 := []byte("resources,exampleTag1=exampleTagValue1 exampleField1=\"exampleFieldValue1\"\n")
	expectedOutputLine2 := []byte("resources,exampleTag2=exampleTagValue2 exampleField2=\"exampleFieldValue2\"\n")
	testObj.On("Write", expectedOutputLine1).Return(len(expectedOutputLine1), nil).On("Write", expectedOutputLine2).Return(len(expectedOutputLine2), nil)
	//when
	OutputMetrics([]models.FilteredMetrics{metric1, metric2}, testObj)

	//then
	testObj.AssertExpectations(t)
}

func TestShouldEscapeTagKeys(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources,example\\,T\\=a\\ g=exampleTagValue exampleField=\"exampleFieldValue\"\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := models.NewFilteredMetric()
	metric.Tags["example,T=a g"] = "exampleTagValue"
	metric.Fields["exampleField"] = "exampleFieldValue"
	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
}

func TestShouldEscapeTagValues(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources,exampleTag=example\\,T\\=a\\ gValue exampleField=\"exampleFieldValue\"\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := models.NewFilteredMetric()
	metric.Tags["exampleTag"] = "example,T=a gValue"
	metric.Fields["exampleField"] = "exampleFieldValue"
	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
}

func TestShouldEscapeFieldKeys(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources,exampleTag=exampleTagValue example\\,F\\=i\\ eld=\"exampleFieldValue\"\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := models.NewFilteredMetric()
	metric.Tags["exampleTag"] = "exampleTagValue"
	metric.Fields["example,F=i eld"] = "exampleFieldValue"
	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
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