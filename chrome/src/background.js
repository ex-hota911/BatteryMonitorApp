var ROOT = 'https://icumn7abiu.appspot.com'
//var ROOT = 'localhost:8080'

// Called when the user clicks on the browser action.
chrome.browserAction.onClicked.addListener(function(tab) {
});

chrome.alarms.onAlarm.addListener(function( alarm ) {
});

chrome.runtime.onInstalled.addListener(function(details){
  chrome.alarms.create("checkBattery", {delayInMinutes: 1, periodInMinutes: 1});
});

var deviceIdKey = "deviceId"

/**
 * Updates the battery status.
 */
function updateBatteryStatus(battery) {
  var deviceId = localStorage.getItem(deviceIdKey);
  if (deviceId == null) {
	console.log("Device is not registered.");
	register();
	return;
  }

  // Make sure level is an int.
  var level = Math.floor(battery.level * 100);
  var charging = (battery.charging)? "charging" : "";
  console.log(battery);
  console.log(level);

  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
	if (xhr.readyState != 4) {
	  return;
	}
	console.log(xhr.status);
	console.log(xhr.responseText);
  };

  var formData = new FormData();
  // TODO: Set device ID.
  formData.append("device_id", deviceId);
  formData.append("battery", level);
  formData.append("charging", charging);

  xhr.open("POST", ROOT + "/api/v1/battery");
  xhr.send(formData);
}

// Register the device
function register() {
  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
	if (xhr.readyState != 4) {
	  return;
	}
	console.log(xhr.status);
	console.log(xhr.responseText);
	if (xhr.status == 200) {
	  var resp = JSON.parse(xhr.responseText);
	  if (resp.id == undefined) {
		console.log("ID is not set in response!");
		return;
	  }
	  localStorage.setItem(deviceIdKey, resp.id);
	} else {
	  //
	}
  };

  xhr.open("POST", ROOT + "/api/v1/register");
  xhr.send();
}

// Register listener.
navigator.getBattery().then(function(battery) {
  console.log(battery);
  updateBatteryStatus(battery)

  battery.addEventListener('levelchange', function() {
	updateBatteryStatus(battery);
  });
  battery.addEventListener('chargingchange', function() {
	updateBatteryStatus(battery);
  });
});
