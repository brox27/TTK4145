package repeater

import (
	"reflect"
)

func Repeater(ch_in interface{}, chs_out ...interface{}) {
	for {
		v, _ := reflect.ValueOf(ch_in).Recv()
		for _, c := range chs_out {
			reflect.ValueOf(c).Send(v)
		}
	}
}
