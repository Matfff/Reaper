# Reaper

#### 0x01. 简介
Reaper是一款基于go语言开发的指纹扫描器，可以对响应数据中的headers、body、mmh3-icon、md5-icon、js、title进行扫描识别，擅长于对大量数据进行批量扫描，经过测试可以很好的完成百万级别数据量的指纹扫描。指纹库reaper.json中有近30000条指纹数据，是从github上其他指纹库提取、收集出来的。指纹收集目前只完成了一部分，后续将会继续追加指纹库数据。

#### 0x02 用法
```shell
reaper -h
  -l string
        扫描的目标URL/主机列表的文件路径(一行一个)
  -o string
        输出结果(csv格式) (default "result.csv")
  -t int
        并发数 (default 100)
```

