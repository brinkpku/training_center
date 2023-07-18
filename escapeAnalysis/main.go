package main

type demo struct {
	Msg string
}

func example() *demo {
	d := &demo{}
	return d
}

func generate8191() {
	nums := make([]int, 8191) // < 64KB
	for i := 0; i < 8191; i++ {
		nums[i] = i
	}
}

func generate8192() {
	nums := make([]int, 8192) // = 64KB
	for i := 0; i < 8192; i++ {
		nums[i] = i
	}
}

func generate(n int) {
	nums := make([]int, n) // 不确定大小
	for i := 0; i < n; i++ {
		nums[i] = i
	}
}

// escap analysis
// go tool compile -l -m -m main.go
// go tool compile -l -m main.go
func main() {
	example()
	generate8191()
	generate8192()
	generate(1)
}
