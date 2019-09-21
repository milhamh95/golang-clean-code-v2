package time_test

import (
	"errors"
	"testing"
	"time"

	nTime "github.com/milhamhidayat/golang-clean-code-v2/pkg/time"

	"github.com/stretchr/testify/require"
)

func TestGetLocalTime(t *testing.T) {
	get, err := nTime.GetLocalTime()
	require.NoError(t, err)
	require.NotNil(t, get)
	t.Log(get)
}

func TestGetUTCTime(t *testing.T) {
	get, err := nTime.GetUTCTime()
	require.NoError(t, err)
	require.NotNil(t, get)
	t.Log(get)
}

func TestConvertToUTCTime(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		jakartaTimezone, err := time.LoadLocation("Asia/Jakarta")
		require.NoError(t, err)
		localTime := time.Date(2019, 9, 21, 7, 0, 0, 0, jakartaTimezone)
		newLocalTime, err := nTime.ConvertTimeWithTimeStamp(localTime)
		require.NoError(t, err)

		utcTimezone, err := time.LoadLocation("UTC")
		require.NoError(t, err)
		expectedTime := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, utcTimezone)
		newExpectedTime, err := nTime.ConvertTimeWithTimeStamp(expectedTime)
		require.NoError(t, err)

		result, err := nTime.ConvertToUTCTime(newLocalTime)
		require.NoError(t, err)

		require.Equal(t, newExpectedTime, result)
		t.Log(newExpectedTime)
		t.Log(result)
	})

	t.Run("failed parsing time", func(t *testing.T) {
		loc := time.Location{}
		localTime := time.Date(-2019, 9, 21, 7, 0, 0, 0, &loc)
		expectedErr := errors.New(`parsing time "-2019-09-21T07:00:00Z" as "2006-01-02T15:04:05Z07:00": cannot parse "-2019-09-21T07:00:00Z" as "2006"`)

		_, err := nTime.ConvertToUTCTime(localTime)
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestConvertTimeToDifferentTimezone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		loc, err := time.LoadLocation("UTC")
		require.NoError(t, err)

		jakartaTimezone, err := time.LoadLocation("Asia/Jakarta")
		require.NoError(t, err)
		localTime := time.Date(2019, 9, 21, 7, 0, 0, 0, jakartaTimezone)

		expectedTime := time.Date(localTime.Year(), localTime.Month(), localTime.Day(), 0, 0, 0, 0, loc)
		newExpectedTime, err := nTime.ConvertTimeWithTimeStamp(expectedTime)
		require.NoError(t, err)

		result, err := nTime.ConvertTimeToDifferentTimezone(localTime, loc)
		require.NoError(t, err)
		require.Equal(t, newExpectedTime, result)
		t.Log(result)
		t.Log(newExpectedTime)
	})

	t.Run("error parsing time", func(t *testing.T) {
		loc := time.Location{}

		expectedErr := errors.New("parsing time \"-2019-09-21T07:00:00Z\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"-2019-09-21T07:00:00Z\" as \"2006\"")

		localTime := time.Date(-2019, 9, 21, 7, 0, 0, 0, &loc)
		_, err := nTime.ConvertTimeToDifferentTimezone(localTime, &loc)
		require.EqualError(t, err, expectedErr.Error())
	})
}

func TestConvertTimeWithTimeStamp(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		jakartaTimezone, err := time.LoadLocation("Asia/Jakarta")
		require.NoError(t, err)
		localTime := time.Date(2019, 9, 21, 7, 0, 0, 0, jakartaTimezone)

		tmpLocalTime := localTime.Format(time.RFC3339)
		expectedLocalTime, err := time.Parse(time.RFC3339, tmpLocalTime)
		require.NoError(t, err)

		result, err := nTime.ConvertTimeWithTimeStamp(localTime)
		require.NoError(t, err)
		require.Equal(t, expectedLocalTime, result)
	})

	t.Run("failed parsing time", func(t *testing.T) {
		loc := time.Location{}
		localTime := time.Date(-2019, 9, 21, 7, 0, 0, 0, &loc)
		expectedErr := errors.New(`parsing time "-2019-09-21T07:00:00Z" as "2006-01-02T15:04:05Z07:00": cannot parse "-2019-09-21T07:00:00Z" as "2006"`)

		_, err := nTime.ConvertTimeWithTimeStamp(localTime)
		require.EqualError(t, err, expectedErr.Error())
	})
}
