package dbparser

import (
	"fmt"
	"time"
)
//To-Do
func RFCToTimestamp(values *[]string) {
	date := (*values)[8]
	expires := (*values)[9]
	last_modified := (*values)[10]
	layout := "Mon, 02 Jan 2006 15:04:05 MST"

	t1, err1 := time.Parse(layout, date)
	if err1 == nil {
		date = fmt.Sprintf("%d", t1.Unix())
		(*values)[8] = date
	}

	t2, err2 := time.Parse(layout, expires)
	if err2 == nil {
		expires = fmt.Sprintf("%d", t2.Unix())
		(*values)[9] = expires
	}

	t3, err3 := time.Parse(layout, last_modified)
	if err3 == nil {
		last_modified = fmt.Sprintf("%d", t3.Unix())
		(*values)[10] = last_modified
	}
}
//To-Do
func DigitToTimestamp(values *[]string) {
	date := (*values)[3]
	date = date + "000000"
	layout := "20060102150405"

	t, err := time.Parse(layout, date)
	if err != nil {
		return
	}
	location, err := time.LoadLocation("MST")
	if err != nil {
		return
	}
	t = t.In(location)
	installdate := fmt.Sprintf("%d", t.Unix())
	(*values)[3] = installdate
}
