package datavolts

import (
	"log"
	"strconv"
	"time"
)

const DataVoltsTableName string = "data-volts"

type DataVolts struct {
	RTensions      []float64 `json:"rTensions"`
	STensions      []float64 `json:"sTensions"`
	TTensions      []float64 `json:"tTensions"`
	RCurrents      []float64 `json:"rCurrents"`
	SCurrents      []float64 `json:"sCurrents"`
	TCurrents      []float64 `json:"tCurrents"`
	RealTimestamp  time.Time `json:"realTimestamp"`
	QueueTimestamp time.Time `json:"queueTimestamp"`
	Timestamp      time.Time `json:"timestamp"`
	MessageID      string    `json:"messageId"`
}

func New(
	realTimestamp string,
	queueTimestamp string,
	messageId string,
) *DataVolts {
	now := time.Now().UTC()
	i, err := strconv.ParseInt(realTimestamp, 10, 64)

	if err != nil {
		log.Fatalln("datavolts err: ", err)
		return &DataVolts{}
	}

	realTimestampDate := time.Unix(i, 0)

	i, err = strconv.ParseInt(queueTimestamp, 10, 64)

	if err != nil {
		log.Fatalln("datavolts err: ", err)
		return &DataVolts{}
	}

	queueTimestampDate := time.Unix(0, i*int64(time.Millisecond))

	return &DataVolts{
		RealTimestamp:  realTimestampDate,
		QueueTimestamp: queueTimestampDate,
		Timestamp:      now,
		MessageID:      messageId,
	}
}

func (dataVolts *DataVolts) AddTensions(values []string, phase string) {
	tensions := make([]float64, len(values))

	for i := range values {
		tensions[i], _ = strconv.ParseFloat(values[i], 64)
	}

	switch phase {
	case "R":
		dataVolts.RTensions = tensions
	case "S":
		dataVolts.STensions = tensions
	case "T":
		dataVolts.TTensions = tensions
	}
}

func (dataVolts *DataVolts) AddCurrents(values []string, phase string) {
	currents := make([]float64, len(values))

	for i := range values {
		currents[i], _ = strconv.ParseFloat(values[i], 64)
	}

	switch phase {
	case "R":
		dataVolts.RTensions = currents
	case "S":
		dataVolts.STensions = currents
	case "T":
		dataVolts.TTensions = currents
	}
}
