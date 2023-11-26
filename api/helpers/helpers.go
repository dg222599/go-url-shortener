package helpers

import (
	"errors"
	"math"
	"os"
	"strings"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	largestuint64 = 18446744073709551615
)


func EnforceHTTP(url string) string{

	 if url[:4]!="http"{
		 url = "http://" + url
	 }

	 return url
}

func RemoveDomainError(url string) bool{

	 if url==os.Getenv("DOMAIN"){
		return false
	 }	

	 updatedURL := strings.Replace(url,"http://","",1)
	 updatedURL = strings.Replace(updatedURL,"https://","",1)
	 updatedURL = strings.Replace(updatedURL,"www.","",1)
	 updatedURL = strings.Split(updatedURL,"/")[0]

	//  if updatedURL == os.Getenv("DOMAIN") {
	// 	return false
	//  }

	 return updatedURL!=os.Getenv("DOMAIN");

	 //return true
}

func Base62Encode(number uint64) string {
	length:= len(alphabet)

	var encodeBuilder strings.Builder

	encodeBuilder.Grow(10)

	for ;number>0;number = number/uint64(length){
		encodeBuilder.WriteByte(alphabet[(number%uint64(length))])
	}

	return encodeBuilder.String()

}

func Base62Decode(encodedString string) (uint64,error) {

	 var number uint64
	 length := len(alphabet)

	 for i,symbol := range encodedString{
		alphabeticPosition := strings.IndexRune(alphabet,symbol)
		if alphabeticPosition == -1{
			return uint64(alphabeticPosition),errors.New("cannot find the symbol in alphabet")
		}

		number+=uint64(alphabeticPosition)*uint64(math.Pow(float64(length),float64(i)))
	}

	return number,nil
}