package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeText(t *testing.T) {
	result := normalizeText([]byte("Воркалось, ХливКие шарЬки пырялись на НОВЕ, - и хрюкотАли зелюКИ.. КАК мюмзики в МАВе"))
	assert.Equal(t, result, "воркалось хливкие шарьки пырялись на нове  и хрюкотали зелюки как мюмзики в маве", "Strings should be equal")

	result = normalizeText([]byte("ॐ भूर्भुवः स्वः तत्सवितुर्वरेण्यं भर्गो देवस्यः धीमहि धियो यो नः प्रचोदयात्"))
	assert.Equal(t, result, "ॐ भरभव सव ततसवतरवरणय भरग दवसय धमह धय य न परचदयत", "Strings should be equal")
}

func TestTop10(t *testing.T) {
	result := top10("ло ло ло ло лооол ла ла ла ла лалла ла лалла ла лааа")
	assert.Equal(t, result, []string{"ла", "ло", "лалла", "лааа", "лооол"}, "Slices should be equal")
}
