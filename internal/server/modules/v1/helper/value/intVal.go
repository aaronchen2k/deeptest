package valueGen

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
)

func GenerateIntItems(start int64, end int64, step interface{}, rand bool, repeat int, repeatTag string) []interface{} {
	if !rand {
		return generateIntItemsByStep(start, end, step.(int), repeat, repeatTag)
	} else {
		return generateIntItemsRand(start, end, step.(int), repeat, repeatTag)
	}
}

func generateIntItemsByStep(start int64, end int64, step int, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	total := 0
	if repeatTag == "" {
		for i := 0; true; {
			val := start + int64(i*step)
			if (val > end && step > 0) || (val < end && step < 0) {
				break
			}

			for round := 0; round < repeat; round++ {
				arr = append(arr, val)

				total++
				if total > consts.MaxNum {
					break
				}
			}

			if total >= consts.MaxNum {
				break
			}
			i++
		}
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := 0; true; {
				val := start + int64(i*step)
				if (val > end && step > 0) || (val < end && step < 0) {
					break
				}

				arr = append(arr, val)

				if total >= consts.MaxNum {
					break
				}
				i++
			}

			if total >= consts.MaxNum {
				break
			}
		}
	}

	return arr
}

func generateIntItemsRand(start int64, end int64, step int, repeat int, repeatTag string) []interface{} {
	arr := make([]interface{}, 0)

	countInRound := (end - start) / int64(step)
	total := 0

	if repeatTag == "" {
		for i := int64(0); i < countInRound; {
			rand := RandNum64(countInRound)
			if step < 0 {
				rand = rand * -1
			}

			val := start + rand
			for round := 0; round < repeat; round++ {
				arr = append(arr, val)

				total++

				if total > consts.MaxNum {
					break
				}
			}

			if total > consts.MaxNum {
				break
			}
			i++
		}
	} else if repeatTag == "!" {
		for round := 0; round < repeat; round++ {
			for i := int64(0); i < countInRound; {
				rand := RandNum64(countInRound)
				if step < 0 {
					rand = rand * -1
				}

				val := start + rand
				arr = append(arr, val)

				if total > consts.MaxNum {
					break
				}
				i++
			}

			if total > consts.MaxNum {
				break
			}
		}
	}

	return arr
}
