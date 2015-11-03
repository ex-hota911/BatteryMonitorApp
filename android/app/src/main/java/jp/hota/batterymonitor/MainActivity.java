package jp.hota.batterymonitor;

import android.accounts.AccountManager;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.support.v7.app.ActionBarActivity;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import com.appspot.icumn7abiu.battery.Battery;
import com.appspot.icumn7abiu.battery.model.History;
import com.appspot.icumn7abiu.battery.model.UpdateReq;
import com.google.android.gms.auth.GoogleAuthUtil;
import com.google.api.client.extensions.android.http.AndroidHttp;
import com.google.api.client.googleapis.extensions.android.gms.auth.GoogleAccountCredential;
import com.google.api.client.json.gson.GsonFactory;
import com.google.api.client.util.DateTime;
import com.google.android.gms.common.AccountPicker;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;


public class MainActivity extends AppCompatActivity {

    final static int REQUEST_PICK_ACCOUNT = 1;

    static String PREF_ACCOUNT_NAME = "ACCOUNT_NAME";

    static TextView textView;
    private Battery service;
    private Button button;
    private SharedPreferences settings;
    private GoogleAccountCredential credential;
    private String accountName;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        textView = (TextView) findViewById(R.id.text);


        credential = GoogleAccountCredential.usingAudience(this,
                "server:client_id:546634630324-mkannoor781g7scn86vodbhol9qss1ev.apps.googleusercontent.com");

        startService(new Intent(this, BatteryLogger.class));
        Battery.Builder builder;
        builder = new Battery.Builder(
                AndroidHttp.newCompatibleTransport(), new GsonFactory(), credential);
        builder.setRootUrl("http://192.168.11.4:8080/_ah/api");
        service = builder.build();

        // Inside your Activity class onCreate method
        settings = getSharedPreferences("BatteryMonitor", 0);

        button = (Button) findViewById(R.id.button);
        button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AsyncTask<Void,Void,Void> task = new AsyncTask<Void, Void, Void>() {
                    @Override
                    protected Void doInBackground(Void... params) {
                        List<History> histories = new ArrayList<>();
                        histories.add(new History().setLevel(10).setTimestamp(new DateTime(new Date())));
                        try {
                            service.battery().update(new UpdateReq().setDeviceId("1").setHistories(histories)).execute();
                            Log.d("Battery", "Success");
                        } catch (IOException e) {
                            Log.e("Battery", "Failed to update.", e);
                        }
                        return null;
                    }
                };
                task.execute();
                return;
            }
        });

        setSelectedAccountName(settings.getString(PREF_ACCOUNT_NAME, null));

            chooseAccount();
    }

    private void chooseAccount() {
        startActivityForResult(credential.newChooseAccountIntent(),
                REQUEST_PICK_ACCOUNT);
    }


    private void setSelectedAccountName(String accountName) {
        SharedPreferences.Editor editor = settings.edit();
        editor.putString(PREF_ACCOUNT_NAME, accountName);
        editor.commit();
        credential.setSelectedAccountName(accountName);
        this.accountName = accountName;
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {

        super.onActivityResult(requestCode, resultCode, data);
        switch (requestCode) {
            case REQUEST_PICK_ACCOUNT:
                if (data != null && data.getExtras() != null) {
                    String accountName =
                            data.getExtras().getString(
                                    AccountManager.KEY_ACCOUNT_NAME);
                    if (accountName != null) {
                        setSelectedAccountName(accountName);
                        SharedPreferences.Editor editor = settings.edit();
                        editor.putString(PREF_ACCOUNT_NAME, accountName);
                        editor.commit();
                        // User is authorized.
                    }
                }
                break;
        }
    }
}
