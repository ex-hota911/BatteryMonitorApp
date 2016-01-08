# Battery Monitor

## Objective

Alert users before their devices batteries die.

## Overview

Android app and Chrome app periodically reports device battery level to the server. When battery level is less than threshold, an alert is sent to other devices.

## Detailed Design

### Tables

+ User
  + Device
    + Battery

*User*

+ (PK) UserId: string
  + email address

*Device*

+ (PK) UserId: string
+ (PK) DeviceId: string
+ DeviceName: string
+ AlertThreshold: int
  + 0 to 100

*Battery*

+ (PK) UserId: string
+ (PK) DeviceId: string
+ (PK) Time: time
+ battery: int

### API

*Update*

Updates or register device and its battery history.

+ Request
  + device: Device
    + (required) device_id
    + (optional) device_name
    + (optional) alert_threshold
  	+ history: []Battery
	    + (required) battery
  	  + (optional) time
+ Response

*List*

Lists devices and battery history.

+ Request
  + (optional) device_id: string
  + (optional) start_date: time
  + (optional) end_date: time
+ Response
  + devices: []Device
