# Similar Sentence Check

[![build](https://github.com/Soontao/similar/actions/workflows/go.yml/badge.svg)](https://github.com/Soontao/similar/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/Soontao/similar/branch/main/graph/badge.svg?token=hfsoMNTwZW)](https://codecov.io/gh/Soontao/similar)



## Get Started

```go

var text100 = "美国联邦调查局13日证实，黑客当天侵入联邦调查局电子邮件服务器，发送了大量虚假信息。"
var text101 = "美国联邦调查局13日，黑客当天入侵联邦调查局邮件服务器，发送了大量虚假信息。"
var text102 = "美国联邦调查局13日证实黑客当天侵入联邦调查局电子邮件服务器发送了大量虚假信息"

func TestSimilar(t *testing.T) {
	assert := assert.New(t)

	similar := NewSimilar(1000)

	r := similar.FindMostSimilar(text100)

	assert.Nil(r.Vec)
	assert.EqualValues(0, r.Similarity)

	similar.Remember(text100)

	r = similar.FindMostSimilar(text100)
	assert.NotNil(r.Vec)
	assert.EqualValues(1, r.Similarity)

	r = similar.FindMostSimilar(text101)
	assert.NotNil(r.GetVector())
	assert.EqualValues(0.9759000729485332, r.GetSimilarity())
	assert.EqualValues(
		text102,
		VecToSentence(*r.Vec, similar.dict),
	)
	assert.EqualValues(
		text102,
		r.ToOriginalSentence(),
	)
}

```