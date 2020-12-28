package com.github.sylvek.domotik.rule;

import com.github.sylvek.domotik.DomotikRulesEngine;

public abstract class BroadcastableAction {

  protected final DomotikRulesEngine.Broadcaster broadcaster;

  public BroadcastableAction(DomotikRulesEngine.Broadcaster broadcaster) {
    this.broadcaster = broadcaster;
  }
}
