package net.sylvek.domotik.app.mqtt;


import org.eclipse.paho.client.mqttv3.IMqttDeliveryToken;
import org.eclipse.paho.client.mqttv3.MqttCallback;
import org.eclipse.paho.client.mqttv3.MqttClient;
import org.eclipse.paho.client.mqttv3.MqttConnectOptions;
import org.eclipse.paho.client.mqttv3.MqttException;
import org.eclipse.paho.client.mqttv3.MqttMessage;

import java.util.Iterator;
import java.util.Vector;

import android.app.Service;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.IBinder;
import android.os.Message;
import android.os.Messenger;
import android.os.RemoteException;
import android.util.Log;


public class MQTTservice extends Service {
    private static boolean serviceRunning = false;
    private static int mid = 0;
    private static MQTTConnection connection = null;
    private final Messenger clientMessenger = new Messenger(new ClientHandler());

    @Override
    public void onCreate()
    {
        super.onCreate();
        connection = new MQTTConnection();
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId)
    {
        boolean publishSomething = intent != null && intent.getStringExtra(TOPIC) != null;

        if (!publishSomething) {
            if (isRunning()) {
                return START_STICKY;
            }

            super.onStartCommand(intent, flags, startId);
        /*
         * Start the MQTT Thread.
		 */
            connection.start();
        } else {
            Message msg = Message.obtain(null, MQTTservice.PUBLISH);
            msg.setData(intent.getExtras());
            connection.makeRequest(msg);
        }

        return START_STICKY;
    }

    @Override
    public void onDestroy()
    {
        connection.end();
    }

    @Override
    public IBinder onBind(Intent intent)
    {
        /*
		 * Return a reference to our client handler.
		 */
        return clientMessenger.getBinder();
    }

    private synchronized static boolean isRunning()
    {
		 /*
		  * Only run one instance of the service.
		  */
        if (serviceRunning == false) {
            serviceRunning = true;
            return false;
        } else {
            return true;
        }
    }

    /*
     * These are the supported messages from bound clients
     */
    public static final int REGISTER = 0;
    public static final int SUBSCRIBE = 1;
    public static final int PUBLISH = 2;

    /*
     * Fixed strings for the supported messages.
     */
    public static final String TOPIC = "topic";
    public static final String MESSAGE = "message";
    public static final String STATUS = "status";
    public static final String CLASSNAME = "classname";
    public static final String INTENTNAME = "intentname";

    /*
     * This class handles messages sent to the service by
     * bound clients.
     */
    class ClientHandler extends Handler {
        @Override
        public void handleMessage(Message msg)
        {
            boolean status = false;

            switch (msg.what) {
                case SUBSCRIBE:
                case PUBLISH:
           		 	/*
           		 	 * These two requests should be handled by
           		 	 * the connection thread, call makeRequest
           		 	 */
                    connection.makeRequest(msg);
                    break;
                case REGISTER: {
                    Bundle b = msg.getData();
                    if (b != null) {
                        Object target = b.getSerializable(CLASSNAME);
                        if (target != null) {
        				 /*
        				  * This request can be handled in-line
        				  * call the API
        				  */
                            connection.setPushCallback((Class<?>) target);
                            status = true;
                        }
                        CharSequence cs = b.getCharSequence(INTENTNAME);
                        if (cs != null) {
                            String name = cs.toString().trim();
                            if (name.isEmpty() == false) {
            				 /*
            				  * This request can be handled in-line
            				  * call the API
            				  */
                                connection.setIntentName(name);
                                status = true;
                            }
                        }
                    }
                    ReplytoClient(msg.replyTo, msg.what, status);
                    break;
                }
            }
        }
    }

    private void ReplytoClient(Messenger responseMessenger, int type, boolean status)
    {
		 /*
		  * A response can be sent back to a requester when
		  * the replyTo field is set in a Message, passed to this
		  * method as the first parameter.
		  */
        if (responseMessenger != null) {
            Bundle data = new Bundle();
            data.putBoolean(STATUS, status);
            Message reply = Message.obtain(null, type);
            reply.setData(data);

            try {
                responseMessenger.send(reply);
            } catch (RemoteException e) {
                // TODO Auto-generated catch block
                e.printStackTrace();
            }
        }
    }

    enum CONNECT_STATE {
        DISCONNECTED,
        CONNECTING,
        CONNECTED
    }

    private class MQTTConnection extends Thread {
        private Class<?> launchActivity = null;
        private String intentName = null;
        private MsgHandler msgHandler = null;
        private static final int STOP = PUBLISH + 1;
        private static final int CONNECT = PUBLISH + 2;
        private static final int RESETTIMER = PUBLISH + 3;
        private CONNECT_STATE connState = CONNECT_STATE.DISCONNECTED;

        MQTTConnection()
        {
            msgHandler = new MsgHandler();
            msgHandler.sendMessage(Message.obtain(null, CONNECT));
        }

        public void end()
        {
            msgHandler.sendMessage(Message.obtain(null, STOP));
        }

        public void makeRequest(Message msg)
        {
			/*
			 * It is expected that the caller only invokes
			 * this method with valid msg.what.
			 */
            msgHandler.sendMessage(Message.obtain(msg));
        }

        public void setPushCallback(Class<?> activityClass)
        {
            launchActivity = activityClass;
        }

        public void setIntentName(String name)
        {
            intentName = name;
        }


        private class MsgHandler extends Handler implements MqttCallback {
            private final String HOST = "192.168.0.4";
            private final int PORT = 1883;
            private final String uri = "tcp://" + HOST + ":" + PORT;
            private final int MINTIMEOUT = 2000;
            private final int MAXTIMEOUT = 32000;
            private int timeout = MINTIMEOUT;
            private MqttClient client = null;
            private MqttConnectOptions options = new MqttConnectOptions();
            private Vector<String> topics = new Vector<String>();


            MsgHandler()
            {
                options.setCleanSession(true);
                try {
                    client = new MqttClient(uri, MqttClient.generateClientId(), null);
                    client.setCallback(this);
                } catch (MqttException e1) {
                    // TODO Auto-generated catch block
                    e1.printStackTrace();
                }
            }

            @Override
            public void handleMessage(Message msg)
            {
                switch (msg.what) {
                    case STOP: {
					/*
					 * Clean up, and terminate.
					 */
                        client.setCallback(null);
                        if (client.isConnected()) {
                            try {
                                client.disconnect();
                                client.close();
                            } catch (MqttException e) {
                                // TODO Auto-generated catch block
                                e.printStackTrace();
                            }
                        }
                        getLooper().quit();
                        break;
                    }
                    case CONNECT: {
                        if (connState != CONNECT_STATE.CONNECTED) {
                            try {
                                client.connect(options);
                                connState = CONNECT_STATE.CONNECTED;
                                Log.d(getClass().getCanonicalName(), "Connected");
                                timeout = MINTIMEOUT;
                            } catch (MqttException e) {
                                Log.d(getClass().getCanonicalName(), "Connection attemp failed with reason code = " + e.getReasonCode() + e.getCause());
                                if (timeout < MAXTIMEOUT) {
                                    timeout *= 2;
                                }
                                this.sendMessageDelayed(Message.obtain(null, CONNECT), timeout);
                                return;
                            }
					    
					    /*
					     * Re-subscribe to previously subscribed topics 
					     */
                            Iterator<String> i = topics.iterator();
                            while (i.hasNext()) {
                                subscribe(i.next());
                            }
                        }
                        break;
                    }
                    case RESETTIMER: {
                        timeout = MINTIMEOUT;
                        break;
                    }
                    case SUBSCRIBE: {
                        boolean status = false;
                        Bundle b = msg.getData();
                        if (b != null) {
                            CharSequence cs = b.getCharSequence(TOPIC);
                            if (cs != null) {
                                String topic = cs.toString().trim();
                                if (topic.isEmpty() == false) {
                                    status = subscribe(topic);
	        					/*
	        					 * Save this topic for re-subscription if needed.
	        					 */
                                    if (status) {
                                        topics.add(topic);
                                    }
                                }
                            }
                        }
                        ReplytoClient(msg.replyTo, msg.what, status);
                        break;
                    }
                    case PUBLISH: {
                        boolean status = false;
                        Bundle b = msg.getData();
                        if (b != null) {
                            CharSequence cs = b.getCharSequence(TOPIC);
                            if (cs != null) {
                                String topic = cs.toString().trim();
                                if (topic.isEmpty() == false) {
                                    cs = b.getCharSequence(MESSAGE);
                                    if (cs != null) {
                                        String message = cs.toString().trim();
                                        if (message.isEmpty() == false) {
                                            status = publish(topic, message);
                                        }
                                    }
                                }
                            }
                        }
                        ReplytoClient(msg.replyTo, msg.what, status);
                        break;
                    }
                }
            }

            private boolean subscribe(String topic)
            {
                try {
                    client.subscribe(topic);
                } catch (MqttException e) {
                    Log.d(getClass().getCanonicalName(), "Subscribe failed with reason code = " + e.getReasonCode());
                    return false;
                }
                return true;
            }

            private boolean publish(String topic, String msg)
            {
                try {
                    MqttMessage message = new MqttMessage();
                    message.setPayload(msg.getBytes());
                    client.publish(topic, message);
                } catch (MqttException e) {
                    Log.d(getClass().getCanonicalName(), "Publish failed with reason code = " + e.getReasonCode());
                    return false;
                }
                return true;
            }

            @Override
            public void connectionLost(Throwable arg0)
            {
                Log.d(getClass().getCanonicalName(), "connectionLost");
                connState = CONNECT_STATE.DISCONNECTED;
                sendMessageDelayed(Message.obtain(null, CONNECT), timeout);
            }

            @Override
            public void deliveryComplete(IMqttDeliveryToken arg0)
            {
            }

            @Override
            public void messageArrived(String topic, MqttMessage message) throws Exception
            {
                Log.d(getClass().getCanonicalName(), topic + ":" + message.toString());

                if (intentName != null) {
                    Intent intent = new Intent();
                    intent.setAction(intentName);
                    intent.putExtra(TOPIC, topic);
                    intent.putExtra(MESSAGE, message.toString());
                    sendBroadcast(intent);
                }
            }
        }
    }
}
