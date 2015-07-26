package jp.hota.batterymonitor;

import android.content.Intent;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.TextView;

import com.appspot.icumn7abiu.battery.Battery;
import com.google.api.client.extensions.android.http.AndroidHttp;
import com.google.api.client.json.gson.GsonFactory;


public class MainActivity extends ActionBarActivity {

    static TextView textView;
    private Battery service;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        textView = (TextView) findViewById(R.id.text);
        startService(new Intent(this, BatteryLogger.class));
        Battery.Builder builder;
        builder = new Battery.Builder(
                AndroidHttp.newCompatibleTransport(), new GsonFactory(), null);
        service = builder.build();
    }
}
