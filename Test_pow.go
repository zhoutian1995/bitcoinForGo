package main

//挖矿原理
// zhouweitian
//79629e0934a27fb866780781bfe4d811eb9fc202faf9dbd97a3a10a6f25bce5b
//zhouweitian1
//ac2afd8c560d7f20c76513f77a909a73c4b486204cbf4235ba595ac7656fd776
//												 00000000000000000

/*
func main()  {
start := time.Now()
for i:=0;i < 1000000000000;i++{
	data := sha256.Sum256([]byte(strconv.Itoa(i)))
	fmt.Printf("%10d,%x\n",i,data)
	fmt.Printf("%s\n",string(data[len(data)-2:]))//取后两位
	if string(data[len(data)-2:]) == "00"{
		usedtime := time.Since(start)
		fmt.Printf("挖矿成功 花费了%d Ms\n",usedtime)
		break
	}
}
var ID []byte
if ID == ([]byte){0}{
	fmt.Println(ID)
}


}
*/