// Load the Visualization API and the piechart package.
google.load('visualization', '1.0', {'packages':['corechart']});

// Set a callback to run when the Google Visualization API is loaded.
google.setOnLoadCallback(drawCharts);

function drawCharts() {
	for (var i = 0; i < battery.length; i++) {
		// TODO(hota): How to pass battery?
	    drawChart(battery[i], "chart_" + i);
    }
}

function drawChart(d, id) {
	var data = new google.visualization.DataTable();
    data.addColumn('datetime', 'Time');
    data.addColumn('number', 'Charging');
    data.addColumn('number', 'On Battery');

    for(var i = 0; i < d.length; i++) {
	  if (i > 0 && d[i-1].charging != d[i].charging) {
		// The state is changed. Adding data points for both in order to connect the lines
		data.addRow([new Date(d[i].time), d[i].battery, d[i].battery]);
	  } else if (d[i].charging) {
		data.addRow([new Date(d[i].time), d[i].battery, null]);
	  } else {
		data.addRow([new Date(d[i].time), null, d[i].battery]);
	  }
    }

    var options = {
      title: 'Battery (%)',
      legend: { position: 'bottom' },
	  series: {
        0: { color: 'red' },
        1: { color: 'blue' },
	  },
	  vAxis: {
		viewWindowMode: 'explicit',
		viewWindow: {
		  min: 0,
		  max: 100,
		}
	  },
    };
    var chart = new google.visualization.AreaChart(document.getElementById(id));
    chart.draw(data, options);
}
