package requests

// Find 使用多级参数查找一个对象中的元素
// args 中的元素可以是 string 或 int, string 代表查询的是字典， int 代表查询的是数组
// 如 JSON 对象 {"abc":[123,456,789]} 中的 789 可以使用参数 ["abc",2] 查出
// 再如 JSON 对象 ["abc",[123,456,789,{"def":111}]] 中的 111 可以使用参数 [1,3,"def"] 查出
// 再如 JSON 对象 ["abc",[123,456,789,{"def":111}]] 中的 {"def":111} 可以使用参数 [1,3] 查出
// 返回值分别为查找结果和是否找到
func Find(obj interface{}, args ...interface{}) (interface{}, bool) {
	if len(args) == 0 {
		return obj, true
	}
	if obj == nil {
		return nil, false
	}
	arg := args[0]
	switch t := arg.(type) {
	case string:
		r, ok := obj.(map[string]interface{})
		if !ok {
			r, ok = obj.(Dict)
			if !ok {
				return nil, false
			}
		}
		r1, ok := r[t]
		if !ok {
			return nil, false
		}
		return Find(r1, args[1:]...)
	case int:
		r, ok := obj.([]interface{})
		if !ok {
			r, ok = obj.(List)
			if !ok {
				return nil, false
			}
		}
		if t >= len(r) {
			return nil, false
		}
		return Find(r[t], args[1:]...)
	}
	return nil, false
}

func FindInt(obj interface{}, args ...interface{}) (int, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(int); ok {
		return r1, true
	}
	return 0, false
}

func FindInt64(obj interface{}, args ...interface{}) (int64, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(int64); ok {
		return r1, true
	}
	return 0, false
}

func FindUint8(obj interface{}, args ...interface{}) (uint8, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(uint8); ok {
		return r1, true
	}
	return 0, false
}

func FindUint16(obj interface{}, args ...interface{}) (uint16, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(uint16); ok {
		return r1, true
	}
	return 0, false
}

func FindUint32(obj interface{}, args ...interface{}) (uint32, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(uint32); ok {
		return r1, true
	}
	return 0, false
}

func FindUint64(obj interface{}, args ...interface{}) (uint64, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(uint64); ok {
		return r1, true
	}
	return 0, false
}

func FindFloat32(obj interface{}, args ...interface{}) (float32, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(float32); ok {
		return r1, true
	}
	return 0, false
}

func FindFloat64(obj interface{}, args ...interface{}) (float64, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return 0, false
	}
	if r1, ok := r.(float64); ok {
		return r1, true
	}
	return 0, false
}

func FindString(obj interface{}, args ...interface{}) (string, bool) {
	r, ok := Find(obj, args...)
	if !ok {
		return "", false
	}
	if r1, ok := r.(string); ok {
		return r1, true
	}
	return "", false
}
