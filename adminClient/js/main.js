// === Admin Console — Eras do Brasil ===

const MAX_LOG_ENTRIES = 500;

// --- DOM refs ---
const dom = {
  status: document.getElementById("connection-status"),
  // Game Loop panel
  glTick: document.getElementById("gl-tick"),
  glDuration: document.getElementById("gl-duration"),
  glUptime: document.getElementById("gl-uptime"),
  glRunning: document.getElementById("gl-running"),
  // World panel
  wsGametime: document.getElementById("ws-gametime"),
  wsPeriod: document.getElementById("ws-period"),
  wsTick: document.getElementById("ws-tick"),
  // Connections panel
  connCount: document.getElementById("conn-count"),
  connLog: document.getElementById("conn-log"),
  // Event log panel
  eventLog: document.getElementById("event-log"),
  // Console
  cmdInput: document.getElementById("cmd-input"),
};

// --- WebSocket ---
let ws = null;
let reconnectTimer = null;

function connect() {
  setStatus("connecting");

  const protocol = location.protocol === "https:" ? "wss:" : "ws:";
  ws = new WebSocket(`${protocol}//${location.host}/ws/admin`);

  ws.onopen = () => {
    setStatus("connected");
    if (reconnectTimer) {
      clearTimeout(reconnectTimer);
      reconnectTimer = null;
    }
  };

  ws.onclose = () => {
    setStatus("disconnected");
    scheduleReconnect();
  };

  ws.onerror = () => {
    ws.close();
  };

  ws.onmessage = (evt) => {
    try {
      const event = JSON.parse(evt.data);
      handleEvent(event);
    } catch (e) {
      console.error("Failed to parse admin event:", e);
    }
  };
}

function scheduleReconnect() {
  if (!reconnectTimer) {
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null;
      connect();
    }, 3000);
  }
}

function setStatus(state) {
  dom.status.className = `status-${state}`;
  dom.status.textContent =
    state === "connected" ? "Conectado" :
    state === "connecting" ? "Conectando..." : "Desconectado";
}

// --- Event routing ---
function handleEvent(event) {
  const { category, type, data, ts } = event;

  switch (category) {
    case "gameloop":
      updateGameLoop(data);
      break;

    case "lobby":
      handleLobbyEvent(type, data, ts);
      break;

    case "system":
      handleSystemEvent(type, data, ts);
      break;

    case "command":
      handleCommandResponse(type, data, ts);
      break;
  }

  // Tudo vai para o event log (exceto ticks — muito ruidoso).
  if (type !== "tick") {
    appendLog(dom.eventLog, ts, type, formatData(type, data));
  }
}

// --- Game Loop panel ---
let lastGameLoopStatus = null;

function updateGameLoop(data) {
  dom.wsTick.textContent = data.tick;

  if (data.game_time) {
    const gt = new Date(data.game_time);
    const h = String(gt.getUTCHours()).padStart(2, "0");
    const m = String(gt.getUTCMinutes()).padStart(2, "0");
    const y = gt.getUTCFullYear();
    const mo = String(gt.getUTCMonth() + 1).padStart(2, "0");
    const d = String(gt.getUTCDate()).padStart(2, "0");
    dom.wsGametime.textContent = `${y}-${mo}-${d} ${h}:${m}`;
  }

  if (data.period) {
    dom.wsPeriod.textContent = data.period;
    dom.wsPeriod.className = `value period-${data.period}`;
  }

  // Atualiza tick counter no game loop panel.
  dom.glTick.textContent = data.tick;
}

function updateGameLoopStatus(status) {
  lastGameLoopStatus = status;
  dom.glRunning.textContent = status.running ? "Yes" : "No";
  dom.glTick.textContent = status.tick_count || status.tick || "—";
  dom.glUptime.textContent = formatDuration(status.up_time);
  dom.glDuration.textContent = formatNanoDuration(status.last_tick_duration);
}

// --- Lobby events (connections panel) ---
function handleLobbyEvent(type, data, ts) {
  if (type === "player_joined" || type === "player_left") {
    if (data.online !== undefined) {
      dom.connCount.textContent = data.online;
    }
  }

  appendLog(dom.connLog, ts, type, formatData(type, data));
}

// --- System events (welcome) ---
function handleSystemEvent(type, data, ts) {
  if (type === "welcome" && data) {
    if (data.game_loop) updateGameLoopStatus(data.game_loop);
    if (data.world) {
      dom.wsTick.textContent = data.world.tick;
      dom.wsPeriod.textContent = data.world.period;
      dom.wsPeriod.className = `value period-${data.world.period}`;
      if (data.world.game_time) {
        const gt = new Date(data.world.game_time);
        const h = String(gt.getUTCHours()).padStart(2, "0");
        const m = String(gt.getUTCMinutes()).padStart(2, "0");
        dom.wsGametime.textContent = `${gt.getUTCFullYear()}-${String(gt.getUTCMonth()+1).padStart(2,"0")}-${String(gt.getUTCDate()).padStart(2,"0")} ${h}:${m}`;
      }
    }
  }
}

// --- Command responses ---
function handleCommandResponse(type, data, ts) {
  if (type === "status" && data) {
    updateGameLoopStatus(data);
  }

  if (type === "world" && data) {
    dom.wsTick.textContent = data.tick;
    dom.wsPeriod.textContent = data.period;
    if (data.game_time) {
      const gt = new Date(data.game_time);
      const h = String(gt.getUTCHours()).padStart(2, "0");
      const m = String(gt.getUTCMinutes()).padStart(2, "0");
      dom.wsGametime.textContent = `${gt.getUTCFullYear()}-${String(gt.getUTCMonth()+1).padStart(2,"0")}-${String(gt.getUTCDate()).padStart(2,"0")} ${h}:${m}`;
    }
  }

  if (type === "players" && data) {
    dom.connCount.textContent = data.online;
  }

  // Command output goes to event log (already handled above).
}

// --- Console input ---
dom.cmdInput.addEventListener("keydown", (e) => {
  if (e.key !== "Enter") return;
  const raw = dom.cmdInput.value.trim();
  if (!raw) return;

  dom.cmdInput.value = "";

  if (!ws || ws.readyState !== WebSocket.OPEN) {
    appendLog(dom.eventLog, new Date().toISOString(), "error", "Nao conectado ao servidor");
    return;
  }

  ws.send(JSON.stringify({ command: raw }));
  appendLog(dom.eventLog, new Date().toISOString(), "cmd", raw);
});

// --- Helpers ---
function appendLog(container, ts, type, text) {
  const entry = document.createElement("div");
  entry.className = "log-entry";

  const tsStr = ts ? formatTimestamp(ts) : "";

  entry.innerHTML =
    `<span class="log-ts">${tsStr}</span>` +
    `<span class="log-type log-type-${type}">[${type}]</span>` +
    `<span class="log-text">${escapeHtml(text)}</span>`;

  container.appendChild(entry);

  // Limita entradas para não estourar memória.
  while (container.children.length > MAX_LOG_ENTRIES) {
    container.removeChild(container.firstChild);
  }

  // Auto-scroll para baixo.
  container.scrollTop = container.scrollHeight;
}

function formatData(type, data) {
  if (!data) return "";
  if (typeof data === "string") return data;

  switch (type) {
    case "player_joined":
      return `${data.player} entrou (online: ${data.online})`;
    case "player_left":
      return `${data.player} saiu (online: ${data.online})`;
    case "player_renamed":
      return `${data.from} → ${data.to}`;
    case "chat":
      return `${data.from}: ${data.body}`;
    case "welcome":
      return data.message || "Admin conectado";
    case "help":
      if (data.commands) {
        return data.commands.map(c => `${c.name} — ${c.desc}`).join(" | ");
      }
      return JSON.stringify(data);
    case "status":
      return `running=${data.running} tick=${data.tick_count} uptime=${formatDuration(data.up_time)} last=${formatNanoDuration(data.last_tick_duration)}`;
    case "world":
      return `tick=${data.tick} time=${data.game_time} period=${data.period}`;
    case "players":
      return `online: ${data.online}`;
    case "error":
      return data.error || JSON.stringify(data);
    default:
      return JSON.stringify(data);
  }
}

function formatTimestamp(ts) {
  try {
    const d = new Date(ts);
    const h = String(d.getHours()).padStart(2, "0");
    const m = String(d.getMinutes()).padStart(2, "0");
    const s = String(d.getSeconds()).padStart(2, "0");
    return `${h}:${m}:${s}`;
  } catch {
    return ts;
  }
}

function formatDuration(nanoseconds) {
  if (!nanoseconds) return "—";
  const seconds = nanoseconds / 1e9;
  if (seconds < 60) return `${seconds.toFixed(0)}s`;
  const minutes = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  if (minutes < 60) return `${minutes}m${secs}s`;
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return `${hours}h${mins}m`;
}

function formatNanoDuration(nanoseconds) {
  if (!nanoseconds) return "—";
  const ms = nanoseconds / 1e6;
  if (ms < 1) return `${(nanoseconds / 1e3).toFixed(0)}us`;
  if (ms < 1000) return `${ms.toFixed(1)}ms`;
  return `${(ms / 1000).toFixed(2)}s`;
}

function escapeHtml(str) {
  const div = document.createElement("div");
  div.textContent = str;
  return div.innerHTML;
}

// --- Start ---
connect();
