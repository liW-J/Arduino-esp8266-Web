/**
 * 接收esp8266发送的json数据，控制led灯
 */
#include <ArduinoJson.h>
#include <SoftwareSerial.h>
#include <OneWire.h>
#include <DallasTemperature.h>
#include <Adafruit_Fingerprint.h>
// Data wire is plugged into port 2 on the Arduino
#define ONE_WIRE_BUS 6
// Setup a oneWire instance to communicate with any OneWire devices (not just Maxim/Dallas temperature ICs)
OneWire oneWire(ONE_WIRE_BUS);
// Pass our oneWire reference to Dallas Temperature. 
DallasTemperature sensors(&oneWire);

#if (defined(__AVR__) || defined(ESP8266)) && !defined(__AVR_ATmega2560__)
// For UNO and others without hardware serial, we must use software serial...
// pin #2 is IN from sensor (GREEN wire)
// pin #3 is OUT from arduino  (WHITE wire)
// Set up the serial port to use softwareserial..
//配置指纹模块串口
SoftwareSerial fingerSerial(2, 3);

#else
// On Leonardo/M0/etc, others with hardware serial, use hardware serial!
// #0 is green wire, #1 is white
#define fingerSerial Serial1

#endif


Adafruit_Fingerprint finger = Adafruit_Fingerprint(&fingerSerial);
StaticJsonDocument<1024> jsonSender;   //声明一个JsonDocument对象，长度200
uint8_t id;

SoftwareSerial ESP_Serial(4, 5); //
void setup() {
  Serial.begin(57600);
  ESP_Serial.begin(57600);
  delay(2000);
  Serial.println(F("Hello ESP8266."));// "初始化输出信息 hello world"
  finger.begin(57600);
  
  id = 1;
  // Start up the library
  sensors.begin();
  
}
 
void loop() {
  
ESP_Serial.listen();
String jsonData= json();
      if(jsonData.length()>0){
       if (jsonData == "Record Temperature......"){
        
          Serial.println("yes!");
          sensors.requestTemperatures(); // Send the command to get temperatures
          Serial.println(sensors.getTempCByIndex(0)); 
          //ESP_Serial.println(sensors.getTempCByIndex(0));
          //String content = String(sensors.getTempCByIndex(0));
          char content[25];
        //Serial.println(content);
          dtostrf(sensors.getTempCByIndex(0)+10,2,2,content);
          jsonSender["type"] = 3;      // 获取温度
          jsonSender["content"] = content;
          String sendString;
          //将指纹认证的结果发送给ESP8266
          serializeJson(jsonSender, ESP_Serial); 
       }else if (jsonData == "Verify fingerprint......"){
          fingerSerial.listen();
          if (finger.verifyPassword()) {
            Serial.println(F("Found fingerprint sensor!"));
          } else {
            Serial.println(F("Did not find fingerprint sensor :("));
            while (1) { delay(1); }
          }
        int getId;
        int count = 0;
        //计数，不能等待超过15次
        while(getId == -1 || getId ==0){
          getId  = getFingerprintIDez();
          count ++;
          Serial.println(count);
          if(getId){
            break;
          }else if(count > 50){
            break;
          }
        }
        char content[25];
        //Serial.println(content);
        itoa(getId, content,10);
        ESP_Serial.listen();
        jsonSender["type"] = 1;      // 验证指纹
        jsonSender["content"] = content;
        //将指纹认证的结果发送给ESP8266
        serializeJson(jsonSender, ESP_Serial); 
          
       }else if(jsonData == "Add fingerprint......" || jsonData == "Update fingerprint......" ){
          
          fingerSerial.listen();
          int newId ;
          int count = 0;
        //计数，不能等待超过15次
        while(newId == -1 || newId ==0){
          newId  = getFingerprintEnroll();
          count ++;
          Serial.println(count);
          if(newId > 0){
            break;
          }else if(count > 50){
            break;
          }
        }
        char content[25];
        //Serial.println(content);
        itoa(newId-1, content,10);
          //string content = String(newId);
          ESP_Serial.listen();
          jsonSender["type"] = 2;      // 添加指纹
          jsonSender["content"] = content;
          //将指纹认证的结果发送给ESP8266
          serializeJson(jsonSender, ESP_Serial); 

       }
       delay(300);
      }else{
          //不做操作
          
           delay(1000);        
   }
    
    
}
//获取sep8266上发送过来的数据
String json(){
  String comdata = "";
   while(ESP_Serial.available() > 0){//串口接收到数据
        comdata += char(ESP_Serial.read());//获取串口接收到的数据
        delay(2);
    }
    Serial.print(comdata);
  return comdata;
  }



// returns -1 if failed, otherwise returns ID #
int getFingerprintIDez() {
  uint8_t p = finger.getImage();
  if (p == FINGERPRINT_NOFINGER){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;

  p = finger.image2Tz();
  if (p == FINGERPRINT_FEATUREFAIL){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;

  p = finger.fingerFastSearch();
  if (p == FINGERPRINT_NOTFOUND){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;

  // found a match!
  Serial.print("Found ID #"); Serial.print(finger.fingerID);
  Serial.print(" with confidence of "); Serial.println(finger.confidence);
  return finger.fingerID;
}

int getFingerprintEnroll() {

  uint8_t p = finger.getImage();
  if (p == FINGERPRINT_NOFINGER){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;

  p = finger.image2Tz(1);
  if (p == FINGERPRINT_FEATUREFAIL){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;


  Serial.println("Remove finger");
  delay(2000);
  
  p = finger.getImage();
  if (p == FINGERPRINT_NOFINGER){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;

  p = finger.image2Tz(2);
  if (p == FINGERPRINT_FEATUREFAIL){
    Serial.println(".");
  }else if (p != FINGERPRINT_OK)  return -1;


  // OK converted!
  Serial.print("Creating model for #");  Serial.println(id);

  p = finger.createModel();
  if (p == FINGERPRINT_OK) {
    Serial.println("Prints matched!");
  } else if (p != FINGERPRINT_OK) return -1;

  Serial.print("ID "); Serial.println(id);
  p = finger.storeModel(id);
  if (p == FINGERPRINT_OK) {
    Serial.println("Stored!");
  } else if (p != FINGERPRINT_OK) return -1;
  
  if (id ==0) return -1;
  id++;
  return id;
}