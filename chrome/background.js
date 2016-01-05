// Copyright (c) 2011 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

var ROOT = 'https://icumn7abiu.appspot.com/_ah/api'

// Called when the user clicks on the browser action.
chrome.browserAction.onClicked.addListener(function(tab) {
  // No tabs or host permissions needed!
  console.log('Turning ' + tab.url + ' red!');
  navigator.getBattery()
	.then(function(battery) {
	  if (battery) {
		console.log(battery.level * 100);
		return battery.level * 100;
	  }
	});

  gapi.client.load('battery', 'v1', function() {
	console.log('ok');
  }, ROOT);

  chrome.tabs.executeScript({
    code: 'document.body.style.backgroundColor="red"'
  });
});

chrome.alarms.onAlarm.addListener(function( alarm ) {
  console.log("Got an alarm!", alarm);
});

chrome.runtime.onInstalled.addListener(function(details){
  chrome.alarms.create("checkBattery", {delayInMinutes: 1, periodInMinutes: 1});
});

function callback2(jsonResp, rawResp) {
  console.log("callback2");
  console.log(jsonResp);
  console.log(rawResp);
}

function callback() {
  console.log("callback");
  gapi.client.request({
	root: ROOT,
	path: '/battery/v1/battery.hello',
	callback: callback2,
  });
}

function gapiIsLoaded() {
  var params = { 'immediate': false };
  
  gapi.auth.authorize(params, function(accessToken) {
	console.log("auth:");
	console.log(accessToken);
    if (accessToken) {
      callback();
    }
  });
}
