package jp.hota.batterymonitor;

import android.accounts.AccountManager;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.support.v7.app.ActionBarActivity;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.provider.Settings;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

import com.google.android.gms.gcm.GoogleCloudMessaging;
import com.google.android.gms.iid.InstanceID;
import com.google.api.client.googleapis.extensions.android.gms.auth.GoogleAccountCredential;

import java.io.IOException;


public class MainActivity extends AppCompatActivity {

    final static int REQUEST_PICK_ACCOUNT = 1;

    final static String PREF_ACCOUNT_NAME = "ACCOUNT_NAME";
    final static String CLIENT_ID = "server:client_id:"
            + "546634630324-mkannoor781g7scn86vodbhol9qss1ev.apps.googleusercontent.com";

    static int level = -1;
    static boolean charging = false;
    static TextView textView;
    private Button button;
    private Button accountButton;
    private SharedPreferences settings;
    private String accountName;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        Log.d(getClass().getName(), "onCreate");
        setContentView(R.layout.activity_main);
        textView = (TextView) findViewById(R.id.text);
        accountButton = (Button) findViewById(R.id.button2);

        startService(new Intent(this, BatteryService.class));

        // Inside your Activity class onCreate method
        settings = getSharedPreferences("BatteryMonitor", 0);

        final String androidId = Settings.Secure.getString(MainActivity.this.getContentResolver(), Settings.Secure.ANDROID_ID);

        button = (Button) findViewById(R.id.button);
        button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                if (level < 0) {
                    return;
                }
                BatteryChangeReceiver.update(v.getContext(), level, charging, null);
            }
        });

        accountButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                chooseAccount();
            }
        });

        findViewById(R.id.button3).setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                startActivity(new Intent(Intent.ACTION_VIEW, Uri.parse(BatteryChangeReceiver.URL)));
            }
        });

        accountName = settings.getString(PREF_ACCOUNT_NAME, null);
        if (accountName != null) {
            accountButton.setText(accountName);
        } else {
            chooseAccount();
        }
        Log.d(this.getClass().getName(), Build.MODEL);


        Log.d(this.getClass().getName(), "Start registering GCM");
        // Start IntentService to register this application with GCM.
        Intent intent = new Intent(this, RegistrationIntentService.class);
        startService(intent);
    }

    private void chooseAccount() {
        startActivityForResult(GoogleAccountCredential.usingAudience(this, CLIENT_ID).newChooseAccountIntent(),
                REQUEST_PICK_ACCOUNT);
    }



    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        super.onActivityResult(requestCode, resultCode, data);
        switch (requestCode) {
            case REQUEST_PICK_ACCOUNT:
                if (data != null && data.getExtras() != null) {
                    accountName = data.getExtras().getString(AccountManager.KEY_ACCOUNT_NAME);
                    if (accountName != null) {
                        SharedPreferences.Editor editor = settings.edit();
                        editor.putString(PREF_ACCOUNT_NAME, accountName);
                        editor.commit();
                        accountButton.setText(accountName);
                    }
                }
                break;
        }
    }
}
