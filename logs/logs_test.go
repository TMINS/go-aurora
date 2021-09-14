package logs

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	fmt.Printf("%s ==> info:\n",LogNowTime())
}
