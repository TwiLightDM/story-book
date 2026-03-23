package helperservice

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const layout = "2006-01-02"

func Validate[T any](t *T) T {
	if t != nil {
		return *t
	}

	var zero T
	return zero
}

func FromBytesToString(b []byte, mime string) string {
	if b == nil {
		return ""
	}
	if mime == "" {
		mime = "image/png"
	}
	return "data:" + mime + ";base64," +
		base64.StdEncoding.EncodeToString(b)
}

func FromStringToBytes(str string) ([]byte, string, error) {
	if str == "" {
		return nil, "", nil
	}

	const prefix = "data:"
	if !strings.HasPrefix(str, prefix) {
		b, err := base64.StdEncoding.DecodeString(str)
		return b, "", err
	}

	parts := strings.SplitN(str, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid data url")
	}

	meta := parts[0]
	data := parts[1]

	mime := ""
	if strings.Contains(meta, ";") {
		mime = strings.TrimPrefix(strings.Split(meta, ";")[0], "data:")
	}

	b, err := base64.StdEncoding.DecodeString(data)
	return b, mime, err
}

func GetLimitAndOffset(c echo.Context) (int, int, error) {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 1000
	offset := 0

	var err error

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			return 0, 0, ErrInvalidLimit
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, ErrInvalidOffset
		}
	}

	return limit, offset, nil
}

func GetStartAndEndDate(c echo.Context) (time.Time, time.Time, error) {
	startDateStr := c.QueryParam("start_date")
	endDateStr := c.QueryParam("end_date")

	var startDate time.Time
	var endDate time.Time
	var err error

	if startDateStr == "" {
		startDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC) // минимальная дата
	} else {
		startDate, err = time.Parse(layout, startDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD")
		}
	}

	if endDateStr == "" {
		endDate = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC) // конец дня
	} else {
		endDate, err = time.Parse(layout, endDateStr)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD")
		}
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	}

	return startDate, endDate, nil
}
