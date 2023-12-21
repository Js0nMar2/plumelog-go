package main

import (
	"reflect"
)

type Greeter interface {
	Greet() string
}

type English struct{}

func (e English) Greet() string {
	return "Hello!"
}

type Spanish struct{}

func (s Spanish) Greet() string {
	return "¡Hola!"
}

type Container struct {
	dependencies map[string]reflect.Type
}

func (c *Container) Provide(name string, dependency interface{}) {
	c.dependencies[name] = reflect.TypeOf(dependency)
}

func (c *Container) Resolve(name string) interface{} {
	dependencyType := c.dependencies[name]
	dependencyValue := reflect.New(dependencyType).Elem()
	dependencyInterface := dependencyValue.Interface()
	return dependencyInterface
}

//func main() {
//container := Container{
//	dependencies: make(map[string]reflect.Type),
//}
//
//container.Provide("english", English{})
//container.Provide("spanish", Spanish{})
//
//englishGreeter := container.Resolve("english").(Greeter)
//spanishGreeter := container.Resolve("spanish").(Greeter)
//
//fmt.Println(englishGreeter.Greet())
//fmt.Println(spanishGreeter.Greet())

//}
