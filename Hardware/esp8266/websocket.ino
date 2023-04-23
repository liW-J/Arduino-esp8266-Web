#include <ESP8266WiFi.h>
#include <WebSocketsClient.h>
#include <ArduinoJson.h>
#include <SoftwareSerial.h>

#define DebugBegin(baud_rate) Serial.begin(baud_rate)
#define DebugPrintln(message) Serial.println(message)
#define DebugPrint(message) Serial.print(message)

const char* AP_SSID = "your_ssid";
const char* AP_PSK = "your_psk";
const unsigned long BAUD_RATE = 57600 ;

//SoftwareSerial arduinoSerial(2, 3);
SoftwareSerial mySerial(5, 4); 

// WebSocket服务器IP地址
String remoteHost = "your_ip";
// 初始化WebSocketsClient对象
WebSocketsClient wsClient;

//int pressState;

void setup() {
  // put your setup code here, to run once:
  // Serial.begin(9600);
  mySerial.begin(57600);
  WiFi.mode(WIFI_STA);  //设置esp8266工作模式
  DebugBegin(BAUD_RATE);
  DebugPrint("Connecting......");
  WiFi.begin(AP_SSID,AP_PSK);
  WiFi.setAutoConnect(true);
  //等待WiFi连接
  while(WiFi.status() != WL_CONNECTED){
    delay(500);
    DebugPrint(".");
  }
  DebugPrintln("");
  DebugPrintln("WiFi connected!");
  DebugPrintln(WiFi.localIP());

  // 启动WebSocket客户端
  wsClient.begin(remoteHost,9090,"/user/ws");
  wsClient.sendTXT("yes!");
  // 指定事件处理函数
  wsClient.onEvent([](WStype_t type, uint8_t * payload, size_t length) {
    if (type == WStype_TEXT) {
      // 接收来自服务端的信息
      String data = (char*)payload;
      StaticJsonDocument<1024> jsonBuffer; //声明一个JsonDocument对象，长度1024
      DeserializationError error = deserializeJson(jsonBuffer, data);
      const char* str = jsonBuffer["data"];  
      if (strlen(str)>0){
        Serial.println(str);
      //给Arduino发串口信息
        mySerial.print(str);
        String jsonData= json();
        //一直循环，直到Arduino返回了数据
        while(jsonData.length()==0){
          jsonData= json();
          if(jsonData.length()>0){
            //解析回收到的数据
            Serial.println(jsonData);
            break;
          }
        }
        wsClient.sendTXT(jsonData);
        Serial.println("success");
      }
    }
  });
}

void loop() {
  // put your main code here, to run repeatedly:
  wsClient.loop();
  // 判断当前FLASH按键状态是否与之前状态一致，若不一致则按键状态发生改变
  // if (digitalRead(0) != pressState) {
  //   // 若按键状态改变，记录当前按键状态，并向服务端发送按键状态
  //   pressState = digitalRead(0);
  //   String data = String("ledStatus=") + (digitalRead(0) == HIGH ? "off": "on");
  //   Serial.println(data);
  // wsClient.sendTXT("yes!");
  // delay(5000);
  // }
  
}

String json(){
  String comdata = "";
   while(mySerial.available() > 0){//串口接收到数据
        comdata += char(mySerial.read());//获取串口接收到的数据
          delay(2);
    }
    //Serial.println(comdata);
  return comdata;
  }



