package jp.hota.batterymonitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.os.BatteryManager;
import android.os.Build;
import android.provider.Settings;
import android.support.annotation.Nullable;
import android.util.Log;

import com.appspot.icumn7abiu.batteryservice.model.Battery;
import com.appspot.icumn7abiu.batteryservice.Batteryservice;
import com.appspot.icumn7abiu.batteryservice.model.Device;
import com.appspot.icumn7abiu.batteryservice.model.UpdateReq;
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
        update(context, level, goAsync());
    }

    public static void update(Context context, int level, @Nullable final PendingResult result) {
        credential =
                GoogleAccountCredential.usingAudience(context, CLIENT_ID);

        SharedPreferences settings = context.getSharedPreferences("BatteryMonitor", 0);
        String accountName = settings.getString(PREF_ACCOUNT_NAME, null);
        if (accountName == null) {
            // TODO: Show notification to select account.
            return;
        }

        credential.setSelectedAccountName(accountName);

        Batteryservice.Builder builder = new Batteryservice.Builder(
                AndroidHttp.newCompatibleTransport(), new GsonFactory(), credential);
        builder.setRootUrl(API_ROOT).setApplicationName("AndroidApp");
        final Batteryservice service = builder.build();

        Log.d(BatteryChangeReceiver.class.getName(), "" + level);
        MainActivity.level = level;
        if (MainActivity.textView != null) {
            MainActivity.textView.setText("" + level);
        }

        if (level < 0) {
            return;
        }

        Battery battery = new Battery()
                .setBattery(level)
                .setTime(new DateTime(new Date()));
        String androidId = Settings.Secure.getString(context.getContentResolver(), Settings.Secure.ANDROID_ID);
        Device device = new Device()
                .setDeviceId(androidId)
                .setDeviceName(Build.MODEL)
                .setBatteries(ImmutableList.of(battery));
        final UpdateReq req = new UpdateReq()
                .setDevice(device);

        AsyncTask<Void, Void, Void> task = new AsyncTask<Void, Void, Void>() {
            @Override
            protected Void doInBackground(Void... params) {
                try {
                    service.update(req).execute();
                    Log.d(getClass().getName(), "Success");
                } catch (IOException e) {
                    Log.e(getClass().getName(), "Failed to update.", e);
                }
                if (result != null) {
                    result.finish();
                }
                return null;
            }
        };
        task.execute();
    }
}
