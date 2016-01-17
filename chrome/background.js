//var ROOT = 'https://icumn7abiu.appspot.com'
var ROOT = 'localhost:8080'

// Called when the user clicks on the browser action.
chrome.browserAction.onClicked.addListener(function(tab) {
});

chrome.alarms.onAlarm.addListener(function( alarm ) {
});

chrome.runtime.onInstalled.addListener(function(details){
  chrome.alarms.create("checkBattery", {delayInMinutes: 1, periodInMinutes: 1});
});

/**
 * Updates the battery status.
 */
function updateBatteryStatus(battery) {
  console.log(battery);
  var level = battery.level * 100;
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
  formData.append("device_id", 'chrome_extension');
  formData.append("battery", level);

  xhr.open("POST", ROOT + "/battery");
  xhr.send(formData);
}

// Register listener.
navigator.getBattery().then(function(battery) {
  console.log(battery);
  battery.addEventListener('levelchange', function() {
	updateBatteryStatus(battery);
  });
  battery.addEventListener('chargingchange', function() {
	updateBatteryStatus(battery);
  });
});
