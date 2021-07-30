package similar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimilar(t *testing.T) {
	assert := assert.New(t)

	similar := NewSimilar(1000)
	assert.EqualValues(
		0,
		similar.Compare("刚刚通报！厦门新增2例确诊病例+2例无症状感染者 行程轨迹公布"),
	)

	assert.EqualValues(
		0.9984286328136085,
		similar.Compare("刚刚通报！厦门新增2例确诊病例+2例无症状感染者 行程轨迹公布?"),
	)

	assert.EqualValues(
		0.7315138028918062,
		similar.Compare("刚刚通报！厦门新增2例确诊病例+2例无症状感染者，为一家人，行程轨迹公布"),
	)
}
