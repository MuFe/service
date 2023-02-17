package utils

import (
	"reflect"
	"testing"
)

func TestAtoi(t *testing.T) {
	got,err := Atoi([]byte("8"))         // 程序输出的结果
	want := 8   // 期望的结果
	if err!=nil||!reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

func TestParseInt(t *testing.T) {
	got,err := ParseInt([]byte("8"),10,64)         // 程序输出的结果
	want := int64(8)   // 期望的结果
	if err!=nil||!reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

func TestParseUint(t *testing.T) {
	got,err := ParseUint([]byte("8"),10,64)         // 程序输出的结果
	want := int64(8)   // 期望的结果
	if err!=nil||!reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}

func TestParseFloat(t *testing.T) {
	got,err := ParseFloat([]byte("8"),10)         // 程序输出的结果
	want := float64(8) // 期望的结果
	if err!=nil||!reflect.DeepEqual(want, got) { // 因为slice不能比较直接，借助反射包中的方法比较
		t.Errorf("expected:%v, got:%v", want, got) // 测试失败输出错误提示
	}
}
