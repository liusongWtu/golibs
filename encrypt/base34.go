/*
 * base34用于生成指定位数的邀请码。
 * 比如我们要生成6位的邀请码，格式是：0-9十个数字，加上24个大写字母（除去O、I连个易混淆字母）的组合。
 * Base34(200441052,8)---004DZRX2
 * Base34ToNum([]byte("004DZRX2"))---200441052
 */
package encrypt

import (
	"container/list"
	"errors"
)

var baseStr = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
var base = []byte(baseStr)
var baseMap map[byte]int //用于base34解码

func init() {
	baseMap = make(map[byte]int)
	for i, v := range base {
		baseMap[v] = i
	}
}

//将uint64类型数字转换成base34编码
func Base34(n uint64, len int) []byte {
	quotient := n
	mod := uint64(0)
	l := list.New()
	for quotient != 0 {
		mod = quotient % 34
		quotient = quotient / 34
		l.PushFront(base[int(mod)])
	}
	listLen := l.Len()

	if listLen >= len {
		res := make([]byte, 0, listLen)
		for i := l.Front(); i != nil; i = i.Next() {

			res = append(res, i.Value.(byte))
		}
		return res
	} else {
		res := make([]byte, 0, len)
		for i := 0; i < len; i++ {
			if i < len-listLen {
				res = append(res, base[0])
			} else {
				res = append(res, l.Front().Value.(byte))
				l.Remove(l.Front())
			}

		}
		return res
	}
}

//将base34字符串解码成uint64
func Base34ToNum(str []byte) (uint64, error) {
	if str == nil || len(str) == 0 {
		return 0, errors.New("parameter is nil or empty")
	}
	var res uint64 = 0
	var r uint64 = 0
	for i := len(str) - 1; i >= 0; i-- {
		v, ok := baseMap[str[i]]
		if !ok {
			return 0, errors.New("character is not base")
		}
		var b uint64 = 1
		for j := uint64(0); j < r; j++ {
			b *= 34
		}
		res += b * uint64(v)
		r++
	}
	return res, nil
}
