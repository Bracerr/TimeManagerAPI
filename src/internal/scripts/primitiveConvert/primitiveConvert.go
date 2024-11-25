package primitiveConvert

import (
	"TimeManagerAuth/src/pkg/customErrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func StringToPrimitiveDate(myTime string) (primitive.DateTime, error) {
	const layout = "2006-01-02T15:04:05Z"
	timeError := time.Now()
	parseTime, err := time.Parse(layout, myTime)
	if err != nil {
		return primitive.NewDateTimeFromTime(timeError), customErrors.NewAppError(http.StatusBadRequest, "неверный формат времени")
	}
	resultTime := primitive.NewDateTimeFromTime(parseTime)
	return resultTime, nil
}
