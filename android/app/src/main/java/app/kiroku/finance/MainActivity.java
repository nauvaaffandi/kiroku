package app.kiroku.finance;

import android.os.Bundle;
import com.getcapacitor.BridgeActivity;
import io.ionic.backgroundrunner.plugin.BackgroundRunnerPlugin;

public class MainActivity extends BridgeActivity {
  @Override
  public void onCreate(Bundle savedInstanceState) {
    super.onCreate(savedInstanceState);
    this.registerPlugin(BackgroundRunnerPlugin.class);
  }
}
