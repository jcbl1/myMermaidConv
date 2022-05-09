package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var Layout string

type L2 struct {
	Content string
	Id      int
}

type L1 struct {
	Name         string
	Id           int
	Subordinates []L2
}

type L0 struct {
	Name         string
	Subordinates []L1
}

func main() {
	fmt.Printf("欢迎使用Notion转Mermaid思维导图自动生成程序\n版本：v0.1.0\n使用中由任何问题请联系开发者：x@rincons.cc\n请将待转换文件放置在用户文件夹下的.mmc文件夹中\n在参数中填入需要的布局（LR, TB等）\n")
	Layout = os.Args[1]
	converter(fileEstab())
	fmt.Println("转换完成！")
}

func fileEstab() (*bufio.Scanner, *bufio.Writer) {
	home := os.Getenv("HOME")
	f, err := os.Open(home + "/.mmc/original.txt")
	if err != nil {
		panic(err.Error())
	}
	f2, err := os.Create(home + "/.mmc/converted.txt")
	if err != nil {
		panic(err.Error())
	}
	return bufio.NewScanner(f), bufio.NewWriter(f2)
}

func converter(scn *bufio.Scanner, wtr *bufio.Writer) {
	//设置布局
	wtr.WriteString("flowchart " + Layout + "\n")
	//定义class
	wtr.WriteString("classDef l0 fill:#AD4B3D,color:#fff\nclassDef l1 fill:#5FB89A,color:#fff\nclassDef l2 fill:#f2f2f2,color:#000\n\n")
	scn.Scan()
	line := scn.Text()
	id0 := L0{Name: string(line)}
	//debug
	//wtr.WriteString(l0.Name + "\n" + "dd")
	//wtr.Flush()
	wtr.WriteString("id0(" + id0.Name + "):::l0\n")
	wtr.Flush()
	sharpCounter := 0
	nowIn := -1
	id := 0
	for scn.Scan() {
		line = scn.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == uint8(35) {
			sharpCounter++
			if len(line) == 1 {
				continue
			}
			if line[1] == uint8(35) {
				sharpCounter++
				if line[2] == uint8(35) {
					sharpCounter++
				}
			}
		}
		switch sharpCounter {
		case 1:
			nowIn++
			id = 0
			tmp := L1{Name: string(line[2:]), Id: nowIn + 1}
			id0.Subordinates = append(id0.Subordinates, tmp)
			sharpCounter = 0
		case 2:
			id++
			tmp := L2{Content: string(line[3:]), Id: id}
			id0.Subordinates[nowIn].Subordinates = append(id0.Subordinates[nowIn].Subordinates, tmp)
			sharpCounter = 0
		default:
			sharpCounter = 0
		}
	}
	for _, u := range id0.Subordinates {
		wtr.WriteString("id" + strconv.Itoa(u.Id) + "(" + u.Name + "):::l1\n")
		for _, v := range u.Subordinates {
			wtr.WriteString("id" + strconv.Itoa(u.Id) + "." + strconv.Itoa(v.Id) + "(" + v.Content + "):::l2\n")
		}
	}
	//写入各节点的联系
	wtr.WriteString("\nid0-->")
	for i, u := range id0.Subordinates {
		if i != 0 {
			wtr.WriteString(" & ")
		}
		wtr.WriteString("id" + strconv.Itoa(u.Id))
	}
	wtr.WriteString("\n")
	for _, u := range id0.Subordinates {
		wtr.WriteString("id" + strconv.Itoa(u.Id) + "-->")
		for j, v := range u.Subordinates {
			if j != 0 {
				wtr.WriteString(" & ")
			}
			wtr.WriteString("id" + strconv.Itoa(u.Id) + "." + strconv.Itoa(v.Id))
		}
		wtr.WriteString("\n")
	}
	wtr.Flush()
	printId0(id0)
}

func printId0(id0 L0) {
	fmt.Println(id0.Name)
	for _, u := range id0.Subordinates {
		fmt.Printf("------")
		fmt.Println(u.Name)
		for _, v := range u.Subordinates {
			fmt.Printf("------------" + v.Content + "\n")
		}
	}
}
