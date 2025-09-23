package main

func main() {
	x := 10
	{
		y := 20
		x = x + y
	}
	return x
}