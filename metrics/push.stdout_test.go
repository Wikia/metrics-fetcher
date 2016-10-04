package metrics

import (
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/Wikia/metrics-fetcher/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	log.SetLevel(log.ErrorLevel)
}

func TestOutputMetricsWhenFiltersEmpty(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	//when
	OutputMetrics([]models.FilteredMetrics{}, testObj)
	//then
	testObj.AssertNotCalled(t, "Write")
}

func TestShouldWriteToOutputWhenNotEmptyFilters(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources,exampleTag=exampleTagValue exampleField=\"exampleFieldValue\"\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := metric(map[string]string{"exampleTag": "exampleTagValue"}, map[string]interface{}{"exampleField": "exampleFieldValue"})
	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
}

func TestShouldErrorWhenNoFields(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	metric := metric(map[string]string{"exampleTag": "exampleTagValue"}, map[string]interface{}{})
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
	metric := metric(map[string]string{}, map[string]interface{}{"exampleField": "exampleFieldValue"})

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
	metric := metric(
		map[string]string{},
		map[string]interface{}{
			"stringField":  "example Field\\Value",
			"intField":     1,
			"floatField":   1.21,
			"booleanField": true,
		},
	)

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
	metric1 := metric(map[string]string{"exampleTag1": "exampleTagValue1"}, map[string]interface{}{"exampleField1": "exampleFieldValue1"})
	metric2 := metric(map[string]string{"exampleTag2": "exampleTagValue2"}, map[string]interface{}{"exampleField2": "exampleFieldValue2"})
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

	metric := metric(map[string]string{
		"example,T=a g": "exampleTagValue",
	}, map[string]interface{}{
		"exampleField": "exampleFieldValue",
	})
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
	metric := metric(map[string]string{
		"exampleTag": "example,T=a gValue",
	}, map[string]interface{}{
		"exampleField": "exampleFieldValue",
	})
	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
}

func TestShouldEscapeFieldKeys(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	expectedOutputLine := []byte("resources example\\,F\\=i\\ eld=\"exampleFieldValue\"\n")
	testObj.On("Write", expectedOutputLine).Return(len(expectedOutputLine), nil)
	metric := metric(map[string]string{}, map[string]interface{}{
		"example,F=i eld": "exampleFieldValue",
	})

	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
}

func TestShouldOutputGivenMetricName(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	testObj.On("Write", mock.Anything).Return(1, nil)
	metric := sampleMetricWithMeasurementName("measurementName")

	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
	calledWithMessage := string(testObj.Calls[0].Arguments.Get(0).([]byte))
	hasMeasurementNameAtBeginning := strings.HasPrefix(calledWithMessage, "measurementName")
	assert.True(t, hasMeasurementNameAtBeginning)
}

func TestShouldEscapeMeasurementName(t *testing.T) {
	//given
	testObj := new(MockedWriter)
	testObj.On("Write", mock.Anything).Return(1, nil)
	metric := sampleMetricWithMeasurementName("measurem,en tName")

	//when
	OutputMetrics([]models.FilteredMetrics{metric}, testObj)
	//then
	testObj.AssertExpectations(t)
	calledWithMessage := string(testObj.Calls[0].Arguments.Get(0).([]byte))
	hasMeasurementNameAtBeginning := strings.HasPrefix(calledWithMessage, "measurem\\,en\\ tName")
	assert.True(t, hasMeasurementNameAtBeginning)
}

type MockedWriter struct {
	mock.Mock
}

func (m *MockedWriter) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func metric(tags map[string]string, fields map[string]interface{}) models.FilteredMetrics {
	metric := models.NewFilteredMetric()
	metric.Measurement = "resources"
	metric.Tags = tags
	metric.Fields = fields
	return metric
}

func sampleMetricWithMeasurementName(measurementName string) models.FilteredMetrics {
	metric := models.NewFilteredMetric()
	metric.Measurement = measurementName
	metric.Fields = map[string]interface{}{"field": "value"}
	return metric
}
