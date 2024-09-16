package stack

import "fmt"

type Stack struct {
	top         int
	stack_slice []any
}

func initStack() *Stack {
	return &Stack{
		stack_slice: make([]any, 0),
		top:         -1,
	}
}

func (st *Stack) Push(input any) {
	st.top++
	st.stack_slice = append(st.stack_slice, input)
}

func (st *Stack) Pop() any {
	if st.IsEmpty() {
		panic("Stack Underflow: Not enough values to pop")
	}
	value := st.stack_slice[st.top]
	st.top--
	return value
}

func (st *Stack) IsEmpty() bool {
	return st.top == -1
}

func (st *Stack) Peek() any {
	if st.IsEmpty() {
		panic("Stack Underflow: Not enough values to peek")
	}
	return st.stack_slice[st.top]
}

func main() {
	mystack := initStack()
	//mystack.Pop()
	mystack.Push(150)
	mystack.Push("foo")
	fmt.Println("peeked value:", mystack.Peek())
	fmt.Println(mystack.Pop())
	fmt.Println("peeked value:", mystack.Peek())
	fmt.Println(mystack.Pop())
	fmt.Println("Stack empty?")
	fmt.Println(mystack.IsEmpty())
}
