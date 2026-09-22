package main

import (
	"fmt"
	"net/http"
	//"strings"
	//"net/url"
	"io"
	"os"
	"bytes"
	"encoding/json"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ------ 结构体定义 ------
// 返回的项目详情
type ProDetailResq struct {
	Id                     int    `gorm:"primaryKey" json:"id"`
	ProgramCode            string `gorm:"size:64" json:"programCode"`
	ProgramName            string `gorm:"size:128" json:"programName"`
	OrgName                string `gorm:"size:64" json:"orgName"`
	Difficulty             string `gorm:"size:64" json:"difficulty"`
	ProgrammingLanguageTag string `gorm:"size:128" json:"programmingLanguageTag"`
}

// 返回的项目列表
type ProListResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Rows []struct {
		ProgramName string `json:"programName"`
		ProgramCode string `json:"programCode"`
		OrgName     string `json:"orgName"`
		ProId       int    `json:"proId"`
	} `json:"rows"`
}

// 项目详情请求表头
type ProDetailRequest struct {
	ProgramId string `json:"programId"`
	Type      string `json:"type"`
}

// 项目列表请求表头
type ProListRequest struct {
	PageNum  string `json:"pageNum"`
	PageSize string `json:"pageSize"`
}

// PDF 请求表头
type PdfRequest struct{
	ProId string `json:"proId"`
}

// ------ 全局变量 ------
var cookie string = "tgt=1790039450.963.7352.139514|2125b4d7ab6829ae59e8509cd0133f66; UM_distinctid=1a0bf30d3cc5dc-097a8c20bd7fbc-4c657b58-1fa400-1a0bf30d3cd1a04; cna=1a0c97183e74489ea41977c1dd7ceb4f; tgt=1790039451.532.7352.70303|180c5bd6c6a9187530d255592766b85d; CNZZDATA1281243141=1096149695-1789914043-%7C1790039463"
var DB *gorm.DB

// 查询项目列表
var URLlist = "https://summer.ospp.ac.cn/2025/api/getProList"

// 查询项目详情
var URLdetail = "https://summer.ospp.ac.cn/2025/api/getProDetail"

// 爬取总页数
var n int = 12

// 初始化数据库
func InitDB() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/programlist?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("连接数据库失败：%v", err)
	}
	DB.AutoMigrate(&ProDetailResq{})
}

// 爬取PDF
func getProgramPpf(proid int) {
	v := PdfRequest{ProId: fmt.Sprintf("%d", proid)}
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("序列化请求体失败%v", err)
	}
	req, err := http.NewRequest("POST", "https://summer.ospp.ac.cn/2025/api/publicApplication", bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("创建请求失败%v", err)
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36 Edg/153.0.0.0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin","https://summer.ospp.ac.cn")
	req.Header.Set("Referer","https://summer.ospp.ac.cn")

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败%v", err)
	}
	body,err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败%v", err)
	}

	if err := os.WriteFile(fmt.Sprintf("%s.pdf",v.ProId), body, 0o644); err != nil {
		fmt.Printf("写入文件失败%v", err)
	}
}

// 获取项目详细信息写入库并爬取 PDF
func getProDetail(id string, orgname string, ch1 chan int) {
	v := ProDetailRequest{ProgramId: id, Type: "org"}
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("序列化请求体失败%v", err)
	}

	req, err := http.NewRequest("POST", URLdetail, bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("创建请求失败%v", err)
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36 Edg/153.0.0.0")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败%v", err)
	}

	var res ProDetailResq
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Printf("解析响应失败:%v", err)
	}
	res.OrgName = orgname
	DB.Create(&res)
	ch1 <- 1
}

// 获取项目列表信息
func getProList(n int, ch chan int) {
	v := ProListRequest{PageNum: fmt.Sprintf("%d", n), PageSize: "50"}

	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("序列化请求体失败%v", err)
	}

	req, err := http.NewRequest("POST", URLlist, bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("创建请求失败%v", err)
	}

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/153.0.0.0 Safari/537.36 Edg/153.0.0.0")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("发送请求失败%v", err)
	}

	var res ProListResp
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Printf("解析响应失败:%v", err)
	}
	ch1 := make(chan int)
	length := len(res.Rows)
	for i := 0; i < length; i++ {
		go getProDetail(res.Rows[i].ProgramCode, res.Rows[i].OrgName, ch1)
		getProgramPpf(res.Rows[i].ProId)
	}
	fmt.Printf("第%d页", n)
	for i := 0; i < length; i++ {
		<-ch1
	}

	ch <- n
}



func main() {

	// 无 PDF
	// 完全单程  1m40.1237254s
	// 列表单程 详细信息最多 50 并行 50.3590556s
	// 列表 12 并行 详细信息 50 并行 5.7801475s
	// 加速比 ≈ 1:17.3

	// 有 PDF
	// 列表 12 并行 详细信息 50 并行 PDF不并行 1m5.1900843s

	InitDB()
	start := time.Now()
	ch := make(chan int)
	for i := 0; i < n; i++ {
		go getProList(i+1, ch)
		fmt.Printf("执行时间: %v\n", time.Since(start))
	}
	for i := 0; i < n; i++ {
		<-ch
	}

	fmt.Printf("最终执行时间: %v\n", time.Since(start))

	

}
