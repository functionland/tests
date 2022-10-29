package com.fmtest; // replace com.your-app-name with your app’s name
import com.facebook.react.bridge.NativeModule;
import com.facebook.react.bridge.ReactApplicationContext;
import com.facebook.react.bridge.ReactContext;
import com.facebook.react.bridge.ReactContextBaseJavaModule;
import com.facebook.react.bridge.ReactMethod;
import com.facebook.react.bridge.Promise;
import java.util.Map;
import java.util.HashMap;
import fulaMobile.FulaMobile;
import fulaMobile.Fula;
import java.io.File;

public class FulaModule extends ReactContextBaseJavaModule {
    String appDir;
    String storeDirPath;
    Fula fula;
   public FulaModule(ReactApplicationContext context) {
       super(context);
       appDir = context.getFilesDir().toString();
        storeDirPath = appDir + "/fula/";
        File storeDir = new File(storeDirPath);
        boolean success = true;
        if (!storeDir.exists()) {
        storeDir.mkdirs();
        }

        try{
        this.fula = FulaMobile.newFula(storeDirPath);
        }
        catch (Exception e) {
            this.fula = null;
        }
   }

   @Override
    public String getName() {
        return "FulaModule";
    }

    @ReactMethod
    public void testGraphExchange(String pid, String ma, String link, Promise promise) throws Exception {
        try {
            String res = this.fula.testGraphExchange(pid, ma, link);
            // String res = pid + ma + link;
            promise.resolve(res);
        } catch (Exception e) {
            // TODO: handle exception
            promise.reject(e);
        }
    }
}