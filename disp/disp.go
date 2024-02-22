package disp

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"strings"
)

func Start(){
	reader := transform.NewReader(bufio.NewReader(os.Stdin), simplifiedchinese.GBK.NewEncoder())
	b,_:=ioutil.ReadAll(reader)
	//b,_=GbkToUtf8(b)
	lines:=strings.Split(string(b),"\n")
	fmt.Println(len(lines))
	//fmt.Println(string(b))
}
func GbkToUtf8(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewEncoder())
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return
}