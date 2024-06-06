package main

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 定义模型
type Product struct {
	ID    uint   `gorm:"primaryKey"`
	Code  string
	Price uint
}

func main() {
	// 连接数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败")
	}

	// 自动迁移模式
	db.AutoMigrate(&Product{})

	// 创建记录
	products := []Product{
		{Code: "A01", Price: 50},
		{Code: "B02", Price: 100},
		{Code: "C03", Price: 150},
		{Code: "D04", Price: 200},
	}

	for _, p := range products {
		db.Create(&p)
		fmt.Printf("创建产品: %+v\n", p)
	}

	// 读取记录
	var product Product
	db.First(&product, 1) // 根据整型主键查找
	fmt.Println("读取ID为1的产品:", product)

	// 更新记录 - 将ID为1的产品的价格更新为300
	db.Model(&product).Update("Price", 300)
	fmt.Println("更新ID为1的产品价格为300")

	// 重新读取记录以确认更新
	db.First(&product, 1)
	fmt.Println("更新后ID为1的产品:", product)

	// 读取并打印所有记录
	var allProducts []Product
	db.Find(&allProducts)
	fmt.Println("所有产品记录:")
	for _, p := range allProducts {
		fmt.Printf("%+v\n", p)
	}

	// 删除ID为1的产品
	db.Delete(&product, 1)
	fmt.Println("删除ID为1的产品")

	// 读取并打印所有记录以确认删除
	db.Find(&allProducts)
	fmt.Println("删除后所有产品记录:")
	for _, p := range allProducts {
		fmt.Printf("%+v\n", p)
	}
}
