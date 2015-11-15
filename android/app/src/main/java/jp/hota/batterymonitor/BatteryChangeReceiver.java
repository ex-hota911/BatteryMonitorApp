package jp.hota.batterymonitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.os.BatteryManager;
import android.os.Build;
import android.provider.Settings;
import android.util.Log;
import android.widget.Toast;

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


public class BatteryChangeReceiver extends BroadcastReceiver {

    final static String CLIENT_ID = MainActivity.CLIENT_ID;
    public static final String URL = "https://icumn7abiu.appspot.com/";
    public static final String API_ROOT = URL + "_ah/api";

    static String PREF_ACCOUNT_NAME = "ACCOUNT_NAME";
    private static GoogleAccountCredential credential;

    public void onReceive(Context context, Intent intent) {

        int currentLevel = intent.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
        int scale = intent.getIntExtra(BatteryManager.EXTRA_SCALE, -1);

        if (currentLevel < 0 || scale <= 0) {
            return;
        }

        int level = (currentLevel * 100) / scale;

        update(context, level);
    }

    public static void update(Context context, int level) {
        credential =
                GoogleAccountCredential.usingAudience(context, CLIENT_ID);

        SharedPreferences settings = context.getSharedPreferences("BatteryMonitor", 0);
        String accountName = settings.getString(PREF_ACCOUNT_NAME, null);
        if (accountName == null) {
            // TODO: Show notification to select account.
            return;
        }

        credential.setSelectedAccountName(accountName);

        Battery.Builder builder;

        builder = new Battery.Builder(
                AndroidHttp.newCompatibleTransport(), new GsonFactory(), credential);
        builder.setRootUrl(API_ROOT);
        final Battery service = builder.build();

        Log.d(BatteryChangeReceiver.class.getName(), "" + level);
        MainActivity.level = level;
        if (MainActivity.textView != null) {
            MainActivity.textView.setText("" + level);
        }

        if (level < 0) {
            return;
        }
        History history = new History()
                .setLevel(level)
                .setTimestamp((new DateTime(new Date())));
        String androidId = Settings.Secure.getString(context.getContentResolver(), Settings.Secure.ANDROID_ID);
        final UpdateReq req = new UpdateReq()
                .setDeviceId(androidId)
                .setDeviceName(Build.MODEL)
                .setHistories(ImmutableList.of(history));
        AsyncTask<Void, Void, Void> task = new AsyncTask<Void, Void, Void>() {
            @Override
            protected Void doInBackground(Void... params) {
                try {
                    service.update(req).execute();
                    Log.d(getClass().getName(), "Success");
                } catch (IOException e) {
                    Log.e(getClass().getName(), "Failed to update.", e);
                }
                return null;
            }
        };
        task.execute();
    }
}
