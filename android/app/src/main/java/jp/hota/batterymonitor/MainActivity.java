package jp.hota.batterymonitor;

import android.accounts.AccountManager;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.os.Bundle;
import android.provider.Settings;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.appspot.icumn7abiu.battery.Battery;
import com.appspot.icumn7abiu.battery.model.History;
import com.appspot.icumn7abiu.battery.model.UpdateReq;
import com.google.api.client.extensions.android.http.AndroidHttp;
import com.google.api.client.googleapis.extensions.android.gms.auth.GoogleAccountCredential;
import com.google.api.client.json.gson.GsonFactory;
import com.google.api.client.util.DateTime;
import com.google.common.collect.ImmutableList;

import java.io.IOException;
import java.util.Date;


public class MainActivity extends AppCompatActivity {

    final static int REQUEST_PICK_ACCOUNT = 1;


    final static String PREF_ACCOUNT_NAME = "ACCOUNT_NAME";
    final static String CLIENT_ID = "server:client_id:"
            + "546634630324-mkannoor781g7scn86vodbhol9qss1ev.apps.googleusercontent.com";

    static int level = -1;
    static TextView textView;
    private Battery service;
    private Button button;
    private Button accountButton;
    private SharedPreferences settings;
    private GoogleAccountCredential credential;
    private String accountName;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        textView = (TextView) findViewById(R.id.text);
        accountButton = (Button) findViewById(R.id.button2);

        credential = GoogleAccountCredential.usingAudience(this, CLIENT_ID);

        startService(new Intent(this, BatteryLogger.class));
        Battery.Builder builder;
        builder = new Battery.Builder(
                AndroidHttp.newCompatibleTransport(), new GsonFactory(), credential);
        builder.setRootUrl("https://icumn7abiu.appspot.com/_ah/api");
        service = builder.build();

        // Inside your Activity class onCreate method
        settings = getSharedPreferences("BatteryMonitor", 0);

        final String androidId = Settings.Secure.getString(MainActivity.this.getContentResolver(), Settings.Secure.ANDROID_ID);

        button = (Button) findViewById(R.id.button);
        button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                AsyncTask<Void, Void, Void> task = new AsyncTask<Void, Void, Void>() {
                    @Override
                    protected Void doInBackground(Void... params) {
                        if (level < 0) {
                            return null;
                        }
                        History history = new History()
                                .setLevel(level)
                                .setTimestamp((new DateTime(new Date())));
                        UpdateReq req = new UpdateReq()
                                .setDeviceId(androidId                               )
                                .setHistories(ImmutableList.of(history));
                        try {
                            service.battery().update(req).execute();
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

        settings.getString(PREF_ACCOUNT_NAME, null);

  //        if (accountName == null) {
        chooseAccount();
        //      } else {
        //        credential.setSelectedAccountName(accountName);
        this.accountName = accountName;
//        }
    }

    private void chooseAccount() {
        startActivityForResult(credential.newChooseAccountIntent(), REQUEST_PICK_ACCOUNT);
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
                            data.getExtras().getString(AccountManager.KEY_ACCOUNT_NAME);
                    if (accountName != null) {
                        setSelectedAccountName(accountName);
                    }
                }
                break;
        }
    }
}
