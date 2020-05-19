package helpers

type TData struct {
	Timestamp int64
	Data      string
}

func TDataQuickSort(td []*TData) {
	size := len(td)
	if size < 2 {
		return
	}

	var i int
	tmp := &TData{}
	pivot := (*td[size/2]).Timestamp

	for j := (size - 1); ; i, j = i+1, j-1 {
		for ; (*td[i]).Timestamp > pivot; i = i + 1 {
		}
		for ; (*td[j]).Timestamp < pivot; j = j - 1 {
		}

		if i >= j {
			break
		}

		tmp = td[i]
		td[i] = td[j]
		td[j] = tmp
	}

	TDataQuickSort(td[:i])
	TDataQuickSort(td[i:])

	return
}
