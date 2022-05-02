# Controlling smart home systems through Telegram bot

<p align="center">
    <img src="/img/p2.jpeg" width="600"/>
</p>

This system allows users to control smart home through telegram bot. Users can open door, turn on/off light, turn on/off secure mode, get temperature data, and receive important notifications (e.g. robbery is happening in the house) anywhere with only telegram application.

## Introduction
The aim of this project is to show the interaction between smart home and telegram bot which is designed to send requests for certain functions of the home and get the status of the functions being performed. The Telegram bot will allow the user to use the following features in the home: control of all light in the room (on/off), request statistics on the air temperature in the room, house entrance gate control (opening/closing), and the most important feature is the notification system, which will send a message in a telegram that there are strangers inside and sound a siren in the house if the secure mode is on.

## Materials
1.  Arduino Mega 2560 and breadboard.
2.  Wemos D1 mini (WiFi development board based on the ESP8266).
3.  Wires.
4.  Human Body Infrared Sensor (PIR).
5.  1 servo motor for automatic door opening and closing.
6.  LEDs.
7.  Active Buzzer.
8.  Resistors.
9.  Temperature Sensor DS18B20.
10. 2 x 7 Segment LED Displays.
11. 2 x BCD to 7 Segment Decoders.
12. Cardboard for the house and door.


## Installation and setup
1.  Setup and run backend service
```
cd microcontrollers && make build
make run
```
2.  Run telegram bot
```
cd tg-home-bot 
go run cmd/main.go
```
3. Disconnect tx0 and rx0 from arduino mega upload mega firmware program
4. While tx0 and rx0 ports are disconnected upload esp firmware to wemos d1 mini
5. Push reset button on arduino mega

## Schema
<p align="center">
    <img src="/img/schema.jpeg" width="600"/>
</p>

## Additional photos
<p align="center">
  <img src="/img/p1.jpeg" width="400"/>
  <img src="/img/p3.jpeg" width="400"/>
</p>