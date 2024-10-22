<p align="center">
  <a href="https://github.com/xxxxfang/Reaper/blob/main/README_EN.md">English</a> •
  <a href="https://github.com/xxxxfang/Reaper/blob/main/README.md">中文</a> 
</p>


---



# Reaper

#### 0x00 Introduction

Reaper is a fingerprint scanner developed based on the Go language. It can scan and identify headers, body, mmh3-icon, md5-icon, js, and title in the response data. It is good at batch scanning of large amounts of data. After testing, it can Very good at completing fingerprint scanning of millions of data volumes. There are 32000+ fingerprint data in Aper.json.

Reaper can recognize two encryption types: mmh3-icon and md5-icon. By setting the value of regexp to start whether to use regular rules for rule matching, it can be well compatible with fingerprint data in other fingerprint libraries.

#### 0x01 How it works

Download the latest version of the executable file: [Reaper](https://github.com/xxxxfang/Reaper/releases)  

```shell
reaper -h
  -l string
        File path to the scanned target URL/host list (one per line)
  -o string
        Output results (csv format) (default "result.csv")
  -t int
        Concurrent threads (default 100)
```

![image](https://github.com/xxxxfang/Reaper/assets/86756456/bd37d09f-88d7-472a-b2cd-c28f06f18332)

![image](https://github.com/xxxxfang/Reaper/assets/86756456/ae555aab-2c99-47ce-9404-72601bba5733)



#### 0x02 reaper.json description

```shell
{
    "fp": "fingerprint-name",  // Match fingerprint name
    "headers": [],             // Match data in response headers
    "body": [],                // Match data in response body
    "icon": [],                // Match mmh3-icon fingerprint or md5-icon fingerprint, for example: "icon": ["-123456789", "a794712345601f2247921cf4c2b99a78"], 
    "js": [],                  // Match the js file name in the response data, for example: "js": ["jquery.js"],   
    "title": [],               // Match the title of the response page
    "regexp": "true"           // Whether to use regular matching
}
```

The matching logic of icon is "OR", the matching logic of other items is "AND", and the overall matching logic is "AND" 

example:

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

If the icon only needs to satisfy one of them, the icon will be true. If the title condition is met again, the title will be true. The final result is: icon && title --> fp-name



#### 0x03 How to use ProxyPool

![image](https://github.com/xxxxfang/Reaper/assets/86756456/5c4578bd-cc62-45aa-846c-b167e95cf148)


![image](https://github.com/xxxxfang/Reaper/assets/86756456/75da86fd-e6a3-48d0-a231-74a6450735e6)




#### 0x05 Reference

Fingerprint data is extracted from data in major open source fingerprint databases. Currently, nearly 30,000 pieces of fingerprint data have been extracted, and more will be added in the future.

[EHole(棱洞)3.0 重构版](https://github.com/EdgeSecurityTeam/EHole)

The fingerprint data of this project comes from the following open source projects:  
[dismap](https://github.com/zhzyker/dismap)  
[cmsprint](https://github.com/Lucifer1993/cmsprint)  
[FingerprintHub](https://github.com/0x727/FingerprintHub)  
[Goby](https://github.com/gobysec/GobyVuls)  
[ObserverWard](https://github.com/0x727/ObserverWard)  
[wappalyzergo](https://github.com/projectdiscovery/wappalyzergo)  
[whatscan](https://github.com/killmonday/whatscan)  
[EHole](https://github.com/EdgeSecurityTeam/EHole)  

The author @r0eXpeR collected fingerprint data from various tools  
[fingerprint](https://github.com/r0eXpeR/fingerprint)  
