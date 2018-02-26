package gquery

import (
	"fmt"
)

func convertByteToString(b byte) string {
	return fmt.Sprintf("%c", b)
}

func ReStrCmp(matchStr, reStr string) bool {
	i, j := 0, 0
	for i < len(matchStr) && j < len(reStr) {
		if matchStr[i] == reStr[j] {
			//fmt.Println("-", convertByteToString(matchStr[i]), convertByteToString(reStr[j]))
			i++
			j++
			continue
		} else {
			if reStr[j] == byte('*') {
				except := byte(0)
				if j+1 < len(reStr) {
					except = reStr[j+1]
				}
				j++
				for i < len(matchStr) {
					//fmt.Println("!", convertByteToString(matchStr[i]), convertByteToString(reStr[j]))
					if matchStr[i] == except {
						i++
						j++
						break
					} else {
						i++
					}
				}
			} else {
				break
				//fmt.Println(reStr[j], '*', convertByteToString(reStr[j]), convertByteToString('*'))
			}
		}
	}
	if i < len(matchStr) || j < len(reStr) {
		return false
	}
	return true
}
