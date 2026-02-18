package viewmodel



type tableSearcher struct {
	scores     []int
	selections []int
}



func EditDist(s, t string) int {

	height := len(t)+1
	width := len(s)+1

	topbuf := make([]int, width)
	buffer := make([]int, width)

	for i := range width {
		topbuf[i] = i
	}
	
	for y:=1; y<height; y++ {
		buffer[0] = y
		for x:=1; x<width; x++ {
			if t[y-1] != s[x-1] {
				del := 1 + topbuf[x]
				ins := 1 + buffer[x-1]
				cha := 1 + topbuf[x-1]
				buffer[x] = min(del, ins, cha)
			} else {
				buffer[x] = topbuf[x-1]
			}
		}
		buffer, topbuf = topbuf, buffer
	}
	return topbuf[width-1]
}
