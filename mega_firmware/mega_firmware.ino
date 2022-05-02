#include <ArduinoJson.h>
#include <OneWire.h>
#include <Servo.h> 

OneWire ds(8); // initialization of temperature sensor DS18B20 with OneWire library it is needed to access 1-wire temperature sensors
Servo servo1;  // initialization of servo motor with Servo library

int temperature = 0; // global variable for storing the temperature from DS18B20 sensor

long lastUpdateTime = 0;          // variable for storing the last update time of reading from temp sensor
const int TEMP_UPDATE_TIME = 100; // tempreature sensor update time frequency

String msg = "";
bool msgReady = false;

bool secureMode = false;        // variable to store secure mode
bool currentDoorState = false;  // variable to store door state whether it is closed or open
bool currentLEDState = false;   // variable to store led state whether it is turned on or off

// this byte data is needed for bcd encoder
byte segValue[10][4] = {
   {0,0,0,0}, // 0
   {0,0,0,1}, // 1
   {0,0,1,0}, // 2
   {0,0,1,1}, // 3
   {0,1,0,0}, // 4
   {0,1,0,1}, // 5
   {0,1,1,0}, // 6
   {0,1,1,1}, // 7
   {1,0,0,0}, // 8
   {1,0,0,1}  // 9  
};

// ports for first 7 segment display
int id0 = 46; 
int id1 = 48; 
int id2 = 50; 
int id3 = 52; 
// ports for second 7 segment display
int id5 = 47;
int id6 = 49;
int id7 = 51;
int id8 = 53;

// function that displays decimal from 0-9 on first 7 segment display using our byte array declared previously
void display_N1(int num) {
  digitalWrite(id3, segValue[num][0]);
  digitalWrite(id2, segValue[num][1]);
  digitalWrite(id1, segValue[num][2]);
  digitalWrite(id0, segValue[num][3]);
}

// function that displays decimal from 0-9 on second 7 segment display using our byte array declared previously
void display_N2(int num) {
  digitalWrite(id8, segValue[num][0]);
  digitalWrite(id7, segValue[num][1]);
  digitalWrite(id6, segValue[num][2]);
  digitalWrite(id5, segValue[num][3]);
}


void setup() {
  pinMode(9, OUTPUT);   // output for active buzzer
  pinMode(10, OUTPUT);  // output for led to turn on or off
  pinMode(12, INPUT);   // input  for pir sensor
  pinMode(13, OUTPUT);  // output for our debug led that blinks every end of the main loop iteration
  
  Serial.begin(115200); // initialize serial communication with data transmission speed of 115200 for tx0 and rx0 ports
  
  digitalWrite(9, !secureMode); // disable active buzzer
  
  servo1.attach(11); // attaching servo motor to port 11
  
  delay(3000);  // wait until start
  
  servo1.write(0);      // save 0 degree position on servo motor
  detectTemperature();  // call detect temperature function before main loop
}

int detectTemperature(){
  byte data[2];
  ds.reset();
  ds.write(0xCC);
  ds.write(0x44);

  if (millis() - lastUpdateTime > TEMP_UPDATE_TIME) {
    lastUpdateTime = millis();
    ds.reset();
    ds.write(0xCC);
    ds.write(0xBE);
    data[0] = ds.read();
    data[1] = ds.read();

    // here we make bit operations with data we received from DS18B20 temperature sensor to obtain celcius degree like number
    temperature = (data[1] << 8) + data[0]; 
    temperature = temperature >> 4;
    
    display_N1(temperature % 10);  // display last digit (smaller by significance) on the first 7 segment display
    display_N2(temperature / 10);  // display first digit (higher by significance) on the second 7 segment display
  }
}

// function to update the door state 
void doorHandler(bool openDoor) {
  if (currentDoorState != openDoor) {
    if (currentDoorState) {
      // close door rotate motor positive side
      servo1.write(-90);
    } else {
      // open door rotate motor negative side
      servo1.write(90);
    }
    currentDoorState = openDoor;
  }
}

// function to update led state
void ledHandler(bool turnLED) {
  if (currentLEDState != turnLED) {
    currentLEDState = turnLED;
  }
}

// function to update secure mode
void smHandler(bool secureM) {
  if (secureMode != secureM) {
    secureMode = secureM;
  }
}

// function which detects is there robbery
bool isRobbery() {
  if (secureMode) {
    bool pirSensor = digitalRead(12);
    
    if (pirSensor) {
      return true;
    }
  }

  return false;
}

// function that sends json data to serial
void transmitData() {
  DynamicJsonDocument doc(512);             // initialization of key-value pair like variable for JSON
  doc["temperature"] = String(temperature); // writing current temperature 
  doc["is_robbery"] = isRobbery();          // writing is there robbery
  serializeJson(doc, Serial);               // sending to the serial tx0, so wemos d1 mini can read data
}
  
void loop() {
  transmitData();   // we send data from sensors to the wemos d1 mini
  
  if (secureMode) {
    bool pirDetected = digitalRead(12);
    
    if (pirDetected) {
      digitalWrite(9, LOW);   // if pir sensor detected motion then enable active buzzer
    } else {
      digitalWrite(9, HIGH);  // otherwise disable active buzzer
    }
  } else {
    digitalWrite(9, HIGH);    // disable active buzzer
  }
  
  digitalWrite(10, currentLEDState); // turn on/off the led according to the updated data got from wemos d1 mini
  
  bool messageReady = false;
  String message = "";

  while (!messageReady) {  // wait for message to be received from wemos d1 mini from rx0 port (serial)
    if (Serial.available()) {
      message = Serial.readString();
      messageReady = true;
    }
  }

  if (messageReady) {
    DynamicJsonDocument doc(512);
    DeserializationError err = deserializeJson(doc, message); // parse the received json data from wemos d1 mini 

    if (err) {
      msgReady = false; // if we could not parse json because of data loss we simply return 
      return;
    }
    // otherwise, we handle the received data and update our state
    doorHandler(doc["is_gate_opened"]);  // update door state
    ledHandler(doc["is_led_turned"]);    // update led  state
    smHandler(doc["secure_mode"]);       // update secure mode
  }
  detectTemperature();  // update temperature

  // we blink with debug led to show every end of main loop
  digitalWrite(13, HIGH);  // turn on led
  delay(500);              // wait 0.5 seconds after every iteration, so arduino mega can proccess the message we just sent to serial
  digitalWrite(13, LOW);   // turn off led
}
