package net.sylvek.domotik.app;

import android.content.Intent;
import android.content.pm.ApplicationInfo;
import android.content.pm.PackageManager;
import android.service.notification.NotificationListenerService;
import android.service.notification.StatusBarNotification;
import android.util.Log;
import net.sylvek.domotik.app.mqtt.MQTTservice;

/**
 * Created by sylvek on 17/12/2015.
 */
public class NotificationListener extends NotificationListenerService {

    @Override
    public void onNotificationPosted(StatusBarNotification sbn)
    {
        Log.d("NotificationListener", sbn.getNotification().toString());
        final CharSequence text = sbn.getNotification().tickerText;
        if (text != null) {
            Log.d("NotificationListener", text.toString());
            publish("triggers/lcd/text", getAppName(sbn.getPackageName()) + ";" + text.toString());
        }
    }

    private String getAppName(String packageName)
    {
        final PackageManager pm = getApplicationContext().getPackageManager();
        ApplicationInfo ai;
        try {
            ai = pm.getApplicationInfo(packageName, 0);
        } catch (final PackageManager.NameNotFoundException e) {
            ai = null;
        }
        return (String) (ai != null ? pm.getApplicationLabel(ai) : "(unknown)");
    }

    private void publish(String title, String text)
    {
        Intent service = new Intent(this, MQTTservice.class);
        service.putExtra(MQTTservice.TOPIC, title);
        service.putExtra(MQTTservice.MESSAGE, text);
        startService(service);
    }

    @Override
    public void onNotificationRemoved(StatusBarNotification sbn)
    {
        // nothing to do.
    }
}
