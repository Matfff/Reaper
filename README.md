# Reaper

#### 0x00 简介
Reaper是一款基于go语言开发的指纹扫描器，可以对响应数据中的headers、body、mmh3-icon、md5-icon、js、title进行扫描识别，擅长于对大量数据进行批量扫描，经过测试可以很好的完成百万级别数据量的指纹扫描。reaper.json中有近30000条指纹数据。

#### 0x01 用法
```shell
reaper -h
  -l string
        扫描的目标URL/主机列表的文件路径(一行一个)
  -o string
        输出结果(csv格式) (default "result.csv")
  -t int
        并发数 (default 100)
```

![image](https://github.com/xxxxfang/Reaper/assets/86756456/bd37d09f-88d7-472a-b2cd-c28f06f18332)  
![image](https://github.com/xxxxfang/Reaper/assets/86756456/ae555aab-2c99-47ce-9404-72601bba5733)


#### 0x02 reaper.json数据格式说明
```shell
{
    "fp": "fingerprint-name",  // 匹配指纹名
    "headers": [],             // 匹配响应头中的数据
    "body": [],                // 匹配响应体中的数据
    "icon": [],                // 匹配mmh3-icon指纹或md5-icon指纹，例："icon": ["-123456789", "a794712345601f2247921cf4c2b99a78"], 
    "js": [],                  // 匹配响应数据中的js文件名，例："js": ["jquery.js"],   
    "title": [],               // 匹配响应页面的标题
    "regexp": "true"           // 是否是用正则匹配
}
```

icon的匹配逻辑是"或"，其余项的匹配逻辑是"与"，整体的匹配逻辑是"与" 

例：
```shell
{
    "fp": "fp-name",
    "headers": [],
    "body": [],
    "icon": ["-123456789", "a794712345601f2247921cf4c2b99a78"],
    "js": [],
    "title": [
        "abcd",
    ],
    "regexp": "true"
}
```
icon只需满足其中之一则icon为true，若再次满足title条件，则title为true，最终结果为：icon && title --> fp-name

#### 0x03 前人栽树、后人乘凉
本项目最初是用的python进行开发，后由于效率问题，在学习了一点go语法后对项目进行了快速重构。  
其中多线程的写法使用了作者 @shihuang 的开源项目Ehole中的写法，指纹数据提取于各大开源指纹库中的数据，目前已提取近30000条指纹数据，后续将会继续追加。  
如有疑问或对代码优化建议，希望不吝赐教，欢迎提交issues

#### 0x04 特别鸣谢
[EHole(棱洞)3.0 重构版](https://github.com/EdgeSecurityTeam/EHole)

本项目指纹数据源于以下开源项目：  
[dismap](https://github.com/zhzyker/dismap)  
[cmsprint](https://github.com/Lucifer1993/cmsprint)  
[FingerprintHub](https://github.com/0x727/FingerprintHub)  
[Goby](https://github.com/gobysec/GobyVuls)  
[ObserverWard](https://github.com/0x727/ObserverWard)  
[wappalyzergo](https://github.com/projectdiscovery/wappalyzergo)  
[whatscan](https://github.com/killmonday/whatscan)  

作者 @r0eXpeR 收集了各种工具中的指纹数据  
[fingerprint](https://github.com/r0eXpeR/fingerprint)  
