package similar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var text100 = "美国联邦调查局13日证实，黑客当天侵入联邦调查局电子邮件服务器，发送了大量虚假信息。"
var text101 = "美国联邦调查局13日，黑客当天入侵联邦调查局邮件服务器，发送了大量虚假信息。"
var text102 = "美国联邦调查局13日证实黑客当天侵入联邦调查局电子邮件服务器发送了大量虚假信息"

func TestSimilar(t *testing.T) {
	assert := assert.New(t)

	similar := NewSimilar(1000)

	r := similar.FindMostSimilar(text100)

	assert.Nil(r.vec)
	assert.EqualValues(0, r.similarity)

	similar.Remember(text100)

	r = similar.FindMostSimilar(text100)
	assert.NotNil(r.vec)
	assert.EqualValues(1, r.similarity)

	r = similar.FindMostSimilar(text101)
	assert.NotNil(r.GetVector())
	assert.EqualValues(0.9759000729485332, r.GetSimilarity())
	assert.EqualValues(
		text102,
		VecToSentence(*r.vec, similar.dict),
	)
	assert.EqualValues(
		text102,
		r.ToOriginalSentence(),
	)
}

var text200 = "The weather today is pretty good"
var text201 = "Today's weather is pretty good"
var text202 = "Today's weather is pretty fine"
var text203 = "Today's weather is really bad"
var text204 = "Tommorry's weather will be really bad"

func Test_FindAllSimilar(t *testing.T) {

	assert := assert.New(t)
	s := NewSimilar(15)
	s.Remember(text100)
	s.Remember(text101)
	s.Remember(text102)
	s.Remember(text200)
	s.Remember(text201)

	r := s.FindSimilar(text202, 0.1)
	assert.EqualValues(2, len(r))
	r = s.FindSimilar(text202, 0.6)
	assert.EqualValues(1, len(r))

	s.Remember(text203)

	r = s.FindSimilar(text202, 0.1)
	assert.EqualValues(3, len(r))
	r = s.FindSimilar(text202, 0.6)
	assert.EqualValues(2, len(r))

	r = s.FindSimilar(text204, 0.9)
	assert.EqualValues(0, len(r))

}
