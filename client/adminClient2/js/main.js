const wsUrl = "ws://localhost:8080/ws/admin";
let socket;

// Elementos da DOM
const statStatus = document.getElementById('stat-status');
const statTick = document.getElementById('stat-tick');
const statTime = document.getElementById('stat-time');
const worldData = document.getElementById('world-data');
const logContainer = document.getElementById('log-container');

function connect() {
    socket = new WebSocket(wsUrl);

    socket.onopen = () => {
        statStatus.textContent = "ONLINE";
        statStatus.className = "status-online";
        appendLog("Conexão com o Hub estabelecida.", "sys");
    };

    socket.onclose = (event) => {
        statStatus.textContent = "OFFLINE";
        statStatus.className = "status-offline";
        appendLog(`Conexão perdida. Motivo: ${event.reason || 'Desconhecido'}`, "err");
        setTimeout(connect, 3000); // Tenta reconectar auto
    };

    socket.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            routeMessage(msg);
        } catch (e) {
            appendLog(`Erro ao ler pacote: ${event.data}`, "err");
        }
    };
}

// O Roteador do Frontend
function routeMessage(msg) {
    // Modo Deus: O motor mandou a fotografia do tick
    if (msg.category === "system" && msg.type === "snapshot") {
        updateWorldView(msg.data);
        return;
    }

    // Erros e avisos
    if (msg.type === "error") {
        appendLog(`ERRO: ${JSON.stringify(msg.data)}`, "err");
        return;
    }

    // Qualquer outra coisa (Welcome, Broadcasts, Respostas de consultas)
    appendLog(`[${msg.type.toUpperCase()}] ${JSON.stringify(msg.data)}`, "info");
}

// Atualiza a barra superior e o painel esquerdo sem poluir o log central
function updateWorldView(snap) {
    statTick.textContent = `Tick: ${snap.tick || snap.tick_count}`;
    if (snap.game_time) statTime.textContent = `Game Time: ${snap.game_time}`;

    // Aqui tu formatas os NPCs e Players recebidos no snapshot de forma bonita
    worldData.innerHTML = `<pre>${JSON.stringify(snap, null, 2)}</pre>`;
}

// Envia comandos estruturados para o AdminRouter do Go
function sendAction(actionType, payload) {
    const cmd = {
        type: actionType,
        payload: payload
    };
    socket.send(JSON.stringify(cmd));
    appendLog(`> Comando enviado: ${actionType}`, "cmd");
}

function appendLog(message, cssClass) {
    const div = document.createElement('div');
    div.className = `log-entry log-${cssClass}`;
    div.textContent = `[${new Date().toLocaleTimeString()}] ${message}`;
    logContainer.appendChild(div);
    logContainer.scrollTop = logContainer.scrollHeight; // Auto-scroll
}

// Inicia a aplicação
connect();