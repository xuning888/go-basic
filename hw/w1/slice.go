package w1

import "errors"

var ErrIndexOutOfRange = errors.New("下标超出范围")
var ErrNullPoint = errors.New("空指针")

// DeleteAt 删除指定位置的元素
// 如果下标不是合法的下标，返回 ErrIndexOutOfRange
func DeleteAt(s []int, idx int) ([]int, error) {

	if s == nil {
		return nil, ErrNullPoint
	}

	// 下标校验
	err := rangeCheck(s, idx)
	if err != nil {
		return nil, err
	}

	temp := make([]int, len(s)-1)
	index := 0
	for i, val := range s {
		if i == idx {
			continue
		}
		temp[index] = val
		index++
	}
	return temp, nil
}

func DeleteAt1(src []int, idx int) ([]int, error) {

	if src == nil {
		return nil, ErrNullPoint
	}

	err := rangeCheck(src, idx)
	if err != nil {
		return nil, err
	}
	// 先切出idx前面的元素, 再追加idx 后面的元素
	// 使用 [::] 的语法，防止后面 append 元素时，修改了源数据
	// 这种方式会导致 res 触发扩容
	res := src[0:idx:idx]
	res = append(res, src[idx+1:]...)
	return res, nil
}

// DeleteAt2 这种方式修改了源切片
func DeleteAt2(src []int, idx int) ([]int, error) {
	if src == nil {
		return nil, ErrNullPoint
	}

	err := rangeCheck(src, idx)
	if err != nil {
		return nil, err
	}

	copy(src[idx:], src[idx+1:])

	res := src[: len(src)-1 : len(src)-1]

	return res, nil
}

func rangeCheck(src []int, idx int) error {
	if idx < 0 || idx >= len(src) {
		return ErrIndexOutOfRange
	}
	return nil
}

func DeleteAt3[T any](src []T, idx int) ([]T, error) {
	if src == nil {
		return nil, ErrNullPoint
	}
	if idx < 0 || idx >= len(src) {
		return nil, ErrIndexOutOfRange
	}
	res := make([]T, len(src)-1)
	copy(res, src[0:idx])
	copy(res[idx:], src[idx+1:])
	return res, nil
}
