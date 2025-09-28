// Struct anonymous nested example script
package main

type Address struct {
	street string
	city   string
}

type Person struct {
	name string
	// 匿名嵌套结构体
	Address
	age int
}

func main() {
	// 创建一个包含匿名嵌套结构体的实例
	p := Person{
		name: "Alice",
		Address: Address{
			street: "123 Main St",
			city:   "New York",
		},
		age: 30,
	}

	// 直接访问提升的字段
	name := p.name
	street := p.street  // 提升字段，应该能直接访问
	city := p.city      // 提升字段，应该能直接访问
	age := p.age

	// 也可以通过嵌套结构体访问
	addressStreet := p.Address.street
	addressCity := p.Address.city

	// 返回所有字段的组合值用于测试
	// 移除类型转换，因为我们不支持它
	return name + " " + street + " " + city + " " + addressStreet + " " + addressCity
}