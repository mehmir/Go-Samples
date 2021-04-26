package main

import (
	"fmt"
	"math"
	//"math"
	"math/rand"
	"time"
)
func decrypt(input, key string) string {
	aKey := []byte(key)
	aInput := []byte(input)
	sum := 0
	for _, v := range aKey {
		sum += int(v) >> (1 << 1 << 1 + -0 ^ 1)
	}
	for i := 0; i < len(aInput); i++ {
		aInput[i] = aInput[i] - byte(i % sum)
	}
	return string(aInput)
}

func encrypt(input, key string) string {
	aKey := []byte(key)
	aInput := []byte(input)
	sum := 0
	for _, v := range aKey {
		sum += int(v) >> (1 << 1 << 1 + -0 ^ 1)
	}
	for i := 0; i < len(aInput); i++ {
		aInput[i] = aInput[i] + byte(i % sum)
	}
	return string(aInput)
}

/*
Inside the decrypt function, characters are shifted by i % sum, i is the position of the character in input. In order to find out about sum, we should
see the first loop. In this loop, each character of key is shifted right by 5, or it is divided by 32. So
if  0 < int(v) < 31  -> 0 or
if 32 < int(v) < 63  -> 1 or
if 64 < int(v) < 95  -> 2 or
if 96 < int(v) < 127 -> 3
is added to sum. (I consider that all of the key characters are from ascii table)
according to the given plainText (Energi) and encryptedText (Eogukn), we can conclude that sum > 5.

				shif(i%sum)	i
E 69	E 69	0			0
o 111	n 110	1			1
g 103	e 101	2			2
u 117	r 114	3			3
k 107	g 103	4			4
n 110	i 105	5			5

Thus, for all of sum greater than 5, the result will be the same because we just have one plainText and one encryptedText, if we have more data, we could
more precisely find the sum. For simplicity, we consider sum = 6.
now we should find a1, a2, and a3 so that a1 + a2 + a3 = sum. Each ai show us a group of 32 characters and we can pick each of those to create our key.
there is no difference between characters of a group because they have same affect on the result and consequently on the decrypted text.

*/

func findKeyReverse(plain, encrypted string, keyLength int) string {
	maxIndex := 0
	for index, a := range plain {
		b := encrypted[index]
		if (int(b) - int(a)) == index && index > maxIndex {
			maxIndex = index
		}
	}
	// according to the decrypt function, keyLength and ascii table, sum should be in this range: maxIndex < sum < 3 * keyLength.
	// Without loss of generality, we assume that sum = maxIndex + 1
	sum := maxIndex + 1
	// we should find a1, a2, ..., ai (i=keyLength) so that a1 + a2 + ... + ai = sum.
	// similar to sum, there are multiple answer for this equation, and any of them can be the solution.
	// then we find the first answer.
	chars := make([]int, keyLength)
	find := false
	for i := range chars {
		for ; chars[i] < 2; chars[i]++ {
			if find { break }
			find = sumArray(chars) == sum - 3
		}
	}
	rand.Seed(time.Now().Unix())
	key := ""
	for _, v := range chars {
		code := rand.Intn(32) + (v+1) * 32
		key += string(code)
	}
	return key
}

func findKeyBrute(plain, encrypted string, keyLength int, concurrentTasks int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvwxyz-=[];',.!@#$%^&*()_+{}:|<>? "
	//key := make([]rune, keyLength)
	//indexes := make([]int, keyLength)
	//for i := 0; i < keyLength; i++ { key[i] = rune(charset[indexes[i]]) }
	a := len(charset)
	fraction := int(a / concurrentTasks)
	ch := make(chan string)
	for i := 0; i < concurrentTasks; i++ {
		indexes := make([]int, keyLength)
		indexes[keyLength - 1] = i * fraction
		start := i * fraction
		end := int(math.Min(float64(len(charset)), float64((i + 1) * fraction)))
		fractionLength := (end - start) * int(math.Pow(float64(len(charset)),float64(keyLength - 2)))
		go findKeyFractionBrute(keyLength,plain,encrypted,indexes,charset,fractionLength, ch)
	}
	key := ""
	for j := 0; j < concurrentTasks; j++ {
		key = <- ch
		if key != "" {
			close(ch)
			return key
		}
	}
	//for k := 0; k < int(math.Pow(float64(len(charset) - 1),float64(keyLength - 1))); k++ {
	//	for i, v := range charset {
	//		key[0] = v
	//		indexes[0] = i
	//		if decrypt(encrypted, string(key)) == plain {
	//			return string(key)
	//		}
	//	}
	//	for j := 0; j < len(indexes) - 1; j++ {
	//		if indexes[j] == len(charset) - 1 {
	//			indexes[j] = 0
	//			indexes[j+1]++
	//		}
	//		key[j] = rune(charset[indexes[j]])
	//		key[j+1] = rune(charset[indexes[j+1]])
	//	}
	//}
	return key
}

func findKeyFractionBrute(keyLength int, plain, encrypted string, indexes []int, charset string, fractionLength int, ch chan string) {
	defer func() {
		recover()
	}()
	key := make([]rune, keyLength)
	for i := 0; i < keyLength; i++ { key[i] = rune(charset[indexes[i]]) }
	for k := 0; k < fractionLength; k++ {
		for i, v := range charset {
			key[0] = v
			indexes[0] = i
			if decrypt(encrypted, string(key)) == plain {
				ch <- string(key)
			}
		}
		for j := 0; j < len(indexes) - 1; j++ {
			if indexes[j] == len(charset) - 1 {
				indexes[j] = 0
				indexes[j+1]++
			}
			key[j] = rune(charset[indexes[j]])
			key[j+1] = rune(charset[indexes[j+1]])
		}
	}
	ch <- ""
}

func sumArray(lst []int) int {
	var sum int = 0
	for _, v := range lst {
		sum += v
	}
	return sum
}

func main() {
	encrypted := "Eogukn"
	plain := "Energi"
	reverseKey := findKeyReverse(plain, encrypted, 3)
	bruteKey := findKeyBrute(plain, encrypted, 3, 4)
	fmt.Println("reverseKey:", reverseKey, " bruteKey:", bruteKey)

	decReverse := decrypt("Eogukn", reverseKey)
	decBrute := decrypt("Eogukn", bruteKey)
	fmt.Println("dec_reverse:", decReverse) // should print "Energi"
	fmt.Println("dec_brute:", decBrute)     // should print "Energi"
	time.Sleep(100*time.Second)
}

