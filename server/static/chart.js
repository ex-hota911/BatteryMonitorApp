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
    data.addColumn('number', 'Battery');
    for(var i = 0; i < d.length; i++) {
		data.addRow([new Date(d[i].time), d[i].battery])
    }
    var options = {
        title: 'Battery (%)',
        curveType: 'function',
        legend: { position: 'bottom' }
    };
    var chart = new google.visualization.LineChart(document.getElementById(id));
    chart.draw(data, options);
}
