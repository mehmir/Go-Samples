# ReverseDecrypt
It is reverse engineering of a decrypt function and also a brut force solution to break a decrypt function using Golang


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
