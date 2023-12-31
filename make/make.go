package make

var numToLetter = make(map[int]byte)
var letterToNum = make(map[byte]int)

var startSmallLts = 97
var endSmallLts = 122
var startBigLts = 65
var endBigLts = 90
var startNums = 48
var downcase = 95

const sizeUrl = 10



/*
	 Порядок следования всех символов в map:
	a-z = 0-25
	A-Z = 26-51
	0-9 = 52-61
	_ = 62
*/

func makeMaps() {
	startSym := byte('a')
	for i := 0; i < 26; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	startSym = byte('A')
	for i := 26; i < 52; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	startSym = byte('0')
	for i := 52; i < 62; i++ {
		numToLetter[i] = startSym
		startSym++
	}
	numToLetter[62] = '_'

	startSym = byte('a')
	numForMap := 0
	for startSym <= 'z' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	startSym = byte('A')
	for startSym <= 'Z' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	startSym = byte('0')
	for startSym <= '9' {
		letterToNum[startSym] = numForMap
		startSym++
		numForMap++
	}
	numToLetter['_'] = 63
}

func NextUrlString(current string) string {
	newUrl := []byte(current)
	n := len(current)
	k := n
	i := sizeUrl - 1
	for i >= 0 {
		if int(newUrl[i]) < n - 1 {
			break
		}
		i -= 1
	}
	if i == -1 {
		return ""
	}
	newUrl[i] += 1
	for j := i + 1; i < k; i++ {
		newUrl[j] = 0
	}
	result := string(newUrl)
	return result
}
