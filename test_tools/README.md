# Test Tools
  
  
### udpclient.go
- 用于并发测试UDP服务器
- 每个连接发送完数据，收到服务器返回数据后就关闭
- 统计成功和失败的次数
---
### macstat.go
- 用于统计mac的数量
- 统计总数和去重后的数量

example:

```go
  package main
  
  import (
    "fmt"
    "io/ioutil"
    "path"
  )

  func main() {
    ms := &MacStat{}
    files, _ := ioutil.ReadDir("./")
    for _, f := range files {
      fileSuffix := path.Ext(f.Name())
      if fileSuffix == ".log" {
        ms.StatFile(f.Name())
        fmt.Println("Parse: ", f.Name())
      }
    }
    ms.Show()
    fmt.Println("finish.")
  }
```
   
