# 基于Arduino+ESP8266的身份识别测温系统

一个简单的物联网项目

后端+硬件：[@Jeanne Willis](https://github.com/liW-J)

前端：[@pan-jy](https://github.com/pan-jy)

### 最终展示形式：Web+硬件

#### Web端

[用户](http://175.178.162.207:9999/#/index)
[管理员](http://175.178.162.207:9999/#/admin/login)

#### 硬件设备

Arduino（Arduino UNO R3）：负责控制传感器

esp8266（NodeMCU 1.0）：负责与服务端进行websocket通信

## 技术栈

- **前端实现**：Vue3 + TypeScript + Vite
- **后端实现**：GoLang + gin + gorm +mysql
- **硬件实现**：Aeduino IDE、C++
- **通信协议**：Websocket、http

