package hetzner_dns

import (
	"encoding/json"
	"time"
)

var (
	HETZNER_TIME_FORMATS = []string{
		"2006-01-02 15:04:05.000 -0700 MST",
		"2006-01-02 15:04:05.00 -0700 MST",
		"2006-01-02 15:04:05.0 -0700 MST",
		"2006-01-02 15:04:05 -0700 MST",
		time.RFC3339,
		time.RFC3339Nano,
	}
)

// HetznerTime is a specialization of time.Time handling Hetzner JSON format for time.
type HetznerTime time.Time

func (hzTime *HetznerTime) String() string {
	if hzTime == nil {
		return "null"
	}
	return time.Time(*hzTime).String()
}

func (hzTime *HetznerTime) MarshalJSON() ([]byte, error) {
	if hzTime == nil {
		return []byte("null"), nil
	}
	return json.Marshal((*time.Time)(hzTime))
}

func (hzTime *HetznerTime) UnmarshalJSON(b []byte) error {
	t := time.Time{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		if (string(b) == `""`) || (string(b) == "") {
			*hzTime = HetznerTime(time.Time{})
			return nil
		}

		for _, timeLayout := range HETZNER_TIME_FORMATS {
			t, err = time.Parse(`"`+timeLayout+`"`, string(b))
			if err == nil {
				*hzTime = HetznerTime(t)
				return nil
			}
		}
		return err
	}
	*hzTime = HetznerTime(t)
	return nil
}
