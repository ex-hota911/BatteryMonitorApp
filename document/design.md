# Battery Monitor

## Tables

*User*

+ (PK) user_id: string - _email address_

*Device*

+ (PK) device_id: string
+ device name: string
+ alert threshold: int - _0 to 100_
+ (optional) history: []Battery

*Battery*

+ (PK) time: time
+ battery: int

## API

*Update*

+ Request
  + devices: []Device
+ Response

*List*

+ Request
  + (optional) device_id: string
  + (optional) start_date: time
  + (optional) end_date: time
+ Response
  + devices: []Device



