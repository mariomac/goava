package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEventually_Manual_Fail(t *testing.T) {
	t.Skip("this test will always fail. Aimed for some manual checks")
	i := 0
	Eventually(t, 5*time.Second, func(t require.TestingT) {
		i++
		fmt.Println("try", i)
		if i == 1 {
			t.FailNow()
		} else {
			t.Errorf("ooooh")
		}
	}, Interval(time.Second))

	fmt.Println("this won't be printed")
}

func TestEventually_Manual_AlwaysFatal(t *testing.T) {
	t.Skip("this test will always fail. Aimed for some manual checks")
	Eventually(t, 5*time.Second, func(t require.TestingT) {
		fmt.Println("try")
		t.FailNow()
	}, Interval(time.Second))

	fmt.Println("this won't be printed")
}

func TestEventually_Manual_Error(t *testing.T) {
	t.Skip("this test will always fail. Aimed for some manual checks")
	Eventually(t, 5*time.Second, func(t require.TestingT) {
		fmt.Println("trying")
		t.Errorf("ooooh")
	}, Interval(time.Second))

	fmt.Println("this will be be printed")
}
