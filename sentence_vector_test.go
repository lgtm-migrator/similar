package similar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSentenceVector_EncodeDecode(t *testing.T) {
	assert := assert.New(t)
	s := SentenceVector{12, 3, 424}
	assert.EqualValues("HP-BAgEBDlNlbnRlbmNlVmVjdG9yAf-CAAEEAAAJ_4IAAxgG_gNQ", s.ToString())
	s2, err := NewSentenceVecFromBase64(s.ToString())
	assert.Nil(err)
	assert.EqualValues(s, s2)
}
