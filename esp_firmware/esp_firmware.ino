#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <ArduinoJson.h>

char* ssid = "NU";              // Wi-fi name         
char* password = "1234512345";  // Wi-fi password

const char* ID = "GKEzceLKQOss";     // your Home ID
const char* CLIENT_ID = "327146290"; // your Chat ID in telegram bot


String serverName = "http://10.101.52.98:8000/api";  // server ip where backend service is running


void setup() {
  WiFi.begin(ssid, password); // connect to wi-fi network with given data
  Serial.begin(115200);       // initialize serial communication with data transmission speed of 115200 for tx0 and rx0 ports
  
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);  // loop until microcontroller is connected to wi-fi network
  }
  
  Serial.println("");
  Serial.print("IP Address: "); 
  Serial.println(WiFi.localIP());  // print IP address given to the microcontroller

  if (WiFi.status() == WL_CONNECTED) {
    WiFiClient client;    // initialization of wi-fi client
    HTTPClient http;      // initializtion of http client to send data to backend service
  
    String serverPath = serverName + "/home/client";  // http endpoint URL
    
    http.begin(client, serverPath.c_str()); // declaring to http client the API we will work with and client through which to send data to the network
    // the idea of this request is to register our smart house in the Backend Service 
    // to say we are ready to accept data
    String json;
    
    DynamicJsonDocument doc(256); // we will use JSON format for data exchange between services
    doc["client_id"] = CLIENT_ID;
    doc["home_id"] = ID;
    serializeJson(doc, json);     // this is data payload for our 

    http.addHeader("Content-Type", "application/json");  // we say that our payload is json formatted
    
    http.POST(json);   // post request with json payload to the network through wi-fi client

    http.end();  // end of the http client 
  }

  delay(3000); // wait until start
}

void loop() {
  bool messageReady = false;
  String message = "";

  while (!messageReady) {  // this loop waits until message is received from Arduino Mega via serial communication
    if (Serial.available()) {
      message = Serial.readString();
      messageReady = true;
    }
  }

  if (WiFi.status() == WL_CONNECTED) { // we check if we are still connected to the wi-fi network
    WiFiClient client;    // initialization of wi-fi client
    HTTPClient http;      // initializtion of http client to send data to backend service
    
    String serverPath = serverName + "/home/client/" + ID; // http endpoint URL for sending the data we received from Arduino Mega
    
    http.begin(client, serverPath); // declaring to http client the API we will work with and client through which to send data to the network
    
    int httpResponseCode = http.PUT(message); // sending request to the backend service with data received from arduino mega
    
    if (httpResponseCode == 200) {
      Serial.println(http.getString()); // if request was successful we send returned data from backend service to Arduino Mega
    }

    http.end(); // end of http client
  }
  
  delay(500); // wait 0.5 seconds after every iteration, so arduino mega can proccess the message we just sent to serial
}
