package app.kiroku.finance;

import android.os.Bundle;
import com.getcapacitor.BridgeActivity;
import com.capacitorjs.plugins.camera.CameraPlugin;
import com.capacitorjs.plugins.clipboard.ClipboardPlugin;
import com.capacitorjs.plugins.geolocation.GeolocationPlugin;
import com.capacitorjs.plugins.localnotifications.LocalNotificationsPlugin;
import com.capacitorjs.plugins.preferences.PreferencesPlugin;
import com.capacitorjs.plugins.toast.ToastPlugin;
import io.ionic.backgroundrunner.plugin.BackgroundRunnerPlugin;

public class MainActivity extends BridgeActivity {
  @Override
  public void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    this.registerPlugin(BackgroundRunnerPlugin.class);
    this.registerPlugin(CameraPlugin.class);
    this.registerPlugin(ClipboardPlugin.class);
    this.registerPlugin(GeolocationPlugin.class);
    this.registerPlugin(LocalNotificationsPlugin.class);
    this.registerPlugin(PreferencesPlugin.class);
    this.registerPlugin(ToastPlugin.class);
  }
}
