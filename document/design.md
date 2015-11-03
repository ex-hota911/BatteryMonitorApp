# Battery Monitor

## Tables

*User*

+ (PK) user_id: string
  + email address

*Device*

+ (PK) device_id: string
+ device name: string
+ alert threshold: int
  + 0 to 100
+ (optional) history: []Battery
  + Used only for API

*Battery*

+ (PK) time: time
+ battery: int

## API

*Update*

+ Request
  + devices: []Device
    + (required) device_id
	+ history: []Battery
	  + (required) battery
	  + (optional) time
+ Response

*List*

+ Request
  + (optional) device_id: string
  + (optional) start_date: time
  + (optional) end_date: time
+ Response
  + devices: []Device



