package app.kiroku.finance;

import android.content.res.Configuration;
import android.graphics.Color;
import android.os.Bundle;
import android.view.View;
import android.view.Window;
import androidx.core.view.WindowCompat;
import androidx.core.view.WindowInsetsControllerCompat;
import com.getcapacitor.BridgeActivity;
import com.capacitorjs.plugins.camera.CameraPlugin;
import com.capacitorjs.plugins.clipboard.ClipboardPlugin;
import com.capacitorjs.plugins.geolocation.GeolocationPlugin;
import com.capacitorjs.plugins.localnotifications.LocalNotificationsPlugin;
import com.capacitorjs.plugins.preferences.PreferencesPlugin;
import com.capacitorjs.plugins.toast.ToastPlugin;

public class MainActivity extends BridgeActivity {
  @Override
  public void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);

    Window window = getWindow();

    // Content does NOT extend behind system bars — web content stops at status/nav bar edges
    WindowCompat.setDecorFitsSystemWindows(window, true);

    // System bars transparent (top + bottom), not drawing any background
    window.setStatusBarColor(Color.TRANSPARENT);
    window.setNavigationBarColor(Color.TRANSPARENT);

    // Icon appearance: dark icons on light background, light icons on dark background
    boolean isDark = (getResources().getConfiguration().uiMode
        & Configuration.UI_MODE_NIGHT_MASK) == Configuration.UI_MODE_NIGHT_YES;
    View decorView = window.getDecorView();
    WindowInsetsControllerCompat insetsController = new WindowInsetsControllerCompat(window, decorView);
    insetsController.setAppearanceLightStatusBars(!isDark);
    insetsController.setAppearanceLightNavigationBars(!isDark);

    this.registerPlugin(CameraPlugin.class);
    this.registerPlugin(ClipboardPlugin.class);
    this.registerPlugin(GeolocationPlugin.class);
    this.registerPlugin(LocalNotificationsPlugin.class);
    this.registerPlugin(PreferencesPlugin.class);
    this.registerPlugin(ToastPlugin.class);
  }
}
