package jp.hota.batterymonitor;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.os.BatteryManager;
import android.util.Log;
import android.widget.Toast;

public class BatteryChangeReceiver extends BroadcastReceiver {
    public void onReceive(Context context, Intent intent) {
        context.unregisterReceiver(this);

        int currentLevel = intent.getIntExtra(BatteryManager.EXTRA_LEVEL, -1);
        int scale = intent.getIntExtra(BatteryManager.EXTRA_SCALE, -1);
        int level = -1;

        if (currentLevel >= 0 && scale > 0) {
            level = (currentLevel * 100) / scale;
        }
        Log.d("BatteryChangeReceiver", "" + level);
        if (MainActivity.textView != null) {
            MainActivity.textView.setText("" + level);
        }
        Toast.makeText(context, "" + level, Toast.LENGTH_SHORT).show();
    }
}
