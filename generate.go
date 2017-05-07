package generator

import (
	"bytes"

	"github.com/summed/gopermutations"
)

type regexCombination []interface{}
type regexPermutation []interface{}
type regexOptional []interface{}
type regexMandatory []interface{}
type regexTerm string
type regexSpace []interface{}

var regs = []regexCombination{
	regexCombination{
		regexSpace{
			regexOptional{
				regexPermutation{
					regexOptional{regexTerm("[[:space:]]*Grand[[:space:]]*")},
					regexOptional{regexTerm("[[:space:]]*1er[[:space:]]*")},
					regexOptional{regexTerm("[[:space:]]*Premiere?[[:space:]]*")},
				},
			},
			regexMandatory{regexTerm("cru")},
			regexOptional{regexTerm("classe?")},
		},
	},
}

func generateRegexp(s interface{}) string {
	var (
		b bytes.Buffer
	)

	switch s.(type) {
	case []interface{}:
		for _, r := range s.([]interface{}) {
			b.WriteString(generateRegexp(r))
		}
	case []regexCombination:
		for _, r := range s.([]regexCombination) {
			b.WriteString(generateRegexp(r))
		}

	case regexCombination:
		for _, r := range s.(regexCombination) {
			b.WriteString(generateRegexp(r))
		}
	case regexSpace:
		for i, r := range s.(regexSpace) {
			b.WriteString(generateRegexp(r))
			if i < len(s.(regexSpace))-1 {
				b.WriteString("[[:space:]]*")
			}
		}
	case regexPermutation:
		c := permutator.GetPermutationChannel(s.(regexPermutation))
		for p := range c {
			if len(p) > 0 {
				for _, i := range p {
					b.WriteString("(?:")
					b.WriteString(
						generateRegexp(i),
					)

					b.WriteRune(')')
				}
			}
		}
	case regexTerm:
		return string(s.(regexTerm))
	case regexMandatory:
		for _, r := range s.(regexMandatory) {
			b.WriteString("(?:" + generateRegexp(r) + ")")
		}
	case regexOptional:
		b.WriteString("(?:")
		for _, r := range s.(regexOptional) {
			b.WriteString(generateRegexp(r))
		}
		b.WriteString(")?")
	}
	return b.String()
}
