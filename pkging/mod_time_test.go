package pkging

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ModTime(t *testing.T) {
	r := require.New(t)

	now := time.Now()

	mt := ModTime(now)

	b, err := mt.MarshalJSON()
	r.NoError(err)

	var mt2 ModTime
	r.NoError(json.Unmarshal(b, &mt2))

	at := time.Time(mt).Format(time.RFC3339)
	bt := time.Time(mt2).Format(time.RFC3339)
	r.Equal(at, bt)
}
