# 基于Arduino+ESP8266的身份识别测温系统

   <img alt="MIT" src="https://img.shields.io/badge/license-MIT-green">
  <img alt="docker build" src="https://img.shields.io/badge/docker%20build-passing-brightgreen">


一个简单的物联网项目:cyclone:

用户通过指纹传感器进行身份认证 身份认证通过后使用温度传感器进行体温测量

管理员端可以查看所有设备使用记录

>后端+硬件：[@Jeanne Willis](https://github.com/liW-J)
>
>前端：[@pan-jy](https://github.com/pan-jy)

### 最终展示形式：Web+硬件

#### Web端

[用户](http://175.178.162.207:9999/#/index)
[管理员](http://175.178.162.207:9999/#/admin/login)

#### 硬件设备

Arduino（Arduino UNO R3）：负责控制传感器

esp8266（NodeMCU 1.0）：负责与服务端进行websocket通信

## 技术栈

- **前端实现**：Vue3 + TypeScript + Vite
- **后端实现**：Golang + gin + gorm +mysql
- **硬件实现**：Arduino IDE、C++
- **通信协议**：Websocket、http

## 使用指南

1. Arduino IDE烧录硬件代码
2. docker部署前后端服务
3. 打开WIFI，esp8266显示连接成功
4. [可以开始测温打卡啦](http://175.178.162.207:9999/#/index)

>**[这里有更详细的介绍](https://slidev.mp333player.com)**

## License
基于MIT开源协议，供学习参考。详细请阅读(LICENSE)
