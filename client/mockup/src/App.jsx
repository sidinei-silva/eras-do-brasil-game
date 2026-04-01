import { useState, useEffect, useRef } from "react";

const C = {
  bg: "#12100e", bgPanel: "#1a1710", bgAlt: "#1f1c16", bgInput: "#0f0e0c",
  border: "#3d3425", borderLight: "#5a4d3a", borderGold: "#8b7340",
  text: "#d4c4a0", textDim: "#8a7d66", textBright: "#f0e6cc",
  gold: "#c9a84c", goldDim: "#8b7340", red: "#a83232", redBright: "#d44",
  blue: "#3a6b8a", blueBright: "#5a9aba", green: "#4a7a3a", greenBright: "#6a5",
  purple: "#7a5a8a", purpleBright: "#a88aba",
};

const font = "'Crimson Pro', Georgia, serif";
const mono = "'Courier Prime', monospace";

// ─── Shared Components ───

function Bar({ current, max, color, label, height = 6 }) {
  return (
    <div style={{ marginBottom: 4 }}>
      <div style={{ display: "flex", justifyContent: "space-between", fontSize: 10, color: C.textDim, marginBottom: 2 }}>
        <span>{label}</span><span>{current}/{max}</span>
      </div>
      <div style={{ height, background: C.bgInput, borderRadius: 3, overflow: "hidden" }}>
        <div style={{ width: `${(current / max) * 100}%`, height: "100%", background: color, borderRadius: 3 }} />
      </div>
    </div>
  );
}

function Btn({ children, active, gold, small, onClick, style: s }) {
  const base = {
    padding: small ? "4px 10px" : "8px 14px",
    fontSize: small ? 11 : 12, fontWeight: 600, fontFamily: font,
    background: active ? `${C.gold}22` : gold ? `${C.gold}15` : "transparent",
    color: active ? C.gold : gold ? C.gold : C.textDim,
    border: `1px solid ${active ? C.gold + "66" : gold ? C.gold + "33" : C.border}`,
    borderRadius: 4, cursor: "pointer", transition: "all .15s", whiteSpace: "nowrap", ...s,
  };
  return <button style={base} onClick={onClick}>{children}</button>;
}

function Panel({ title, children, style: s }) {
  return (
    <div style={{ background: C.bgPanel, border: `1px solid ${C.border}`, borderRadius: 6, overflow: "hidden", ...s }}>
      {title && <div style={{ padding: "8px 12px", fontSize: 10, color: C.textDim, letterSpacing: 1.5, textTransform: "uppercase", fontWeight: 600, borderBottom: `1px solid ${C.border}` }}>{title}</div>}
      <div style={{ padding: 12 }}>{children}</div>
    </div>
  );
}

function Divider() { return <div style={{ height: 1, background: C.border, margin: "12px 0" }} />; }

function Tag({ children, color = C.gold }) {
  return <span style={{ fontSize: 10, padding: "2px 6px", background: `${color}18`, color, border: `1px solid ${color}33`, borderRadius: 3 }}>{children}</span>;
}

// ─── Minimap ───

const nodes = [
  { id: "vila", name: "Vila de São Tomé", icon: "🏠", x: 50, y: 55, lv: "Hub", visited: true, current: true },
  { id: "floresta", name: "Floresta do Norte", icon: "🌲", x: 45, y: 20, lv: "1-3", visited: true },
  { id: "rio", name: "Rio das Marés", icon: "🌊", x: 82, y: 45, lv: "1-2", visited: true },
  { id: "mina", name: "Mina de Ouro", icon: "⛏", x: 90, y: 18, lv: "2-4" },
  { id: "toca", name: "Toca da Fera", icon: "🐺", x: 22, y: 8, lv: "3-4" },
  { id: "camp", name: "Acampamento", icon: "⚔", x: 12, y: 30, lv: "3-5" },
  { id: "pico", name: "Pico da Neblina", icon: "🏔", x: 15, y: 70, lv: "3-5" },
  { id: "ruinas", name: "Ruínas Queimadas", icon: "🔥", x: 60, y: 85, lv: "4-5" },
  { id: "ruptura", name: "A Ruptura", icon: "🌀", x: 35, y: 48, lv: "5+" },
];
const edges = [["vila","floresta"],["vila","rio"],["vila","pico"],["vila","ruinas"],["rio","mina"],["floresta","toca"],["floresta","camp"],["floresta","ruptura"]];

function MiniMap({ size = "small", onSelect }) {
  const [hov, setHov] = useState(null);
  const big = size === "big";
  return (
    <div style={{ position: "relative", width: "100%", aspectRatio: big ? "1.6" : "1.2", background: `${C.bgInput}88`, borderRadius: 6, border: `1px solid ${C.border}`, overflow: "hidden" }}>
      <svg width="100%" height="100%" viewBox="0 0 100 100" style={{ position: "absolute" }}>
        {edges.map(([a, b], i) => {
          const pa = nodes.find(n => n.id === a), pb = nodes.find(n => n.id === b);
          return <line key={i} x1={pa.x} y1={pa.y} x2={pb.x} y2={pb.y} stroke={C.border} strokeWidth={big ? 0.4 : 0.5} strokeDasharray="2,2" />;
        })}
      </svg>
      {nodes.map(n => (
        <div key={n.id} onMouseEnter={() => setHov(n.id)} onMouseLeave={() => setHov(null)} onClick={() => onSelect?.(n)}
          style={{ position: "absolute", left: `${n.x}%`, top: `${n.y}%`, transform: "translate(-50%,-50%)", cursor: "pointer", display: "flex", flexDirection: "column", alignItems: "center", opacity: n.visited ? 1 : 0.35, transition: "all .2s", filter: hov === n.id ? "brightness(1.3)" : "none" }}>
          <div style={{ width: n.current ? (big ? 36 : 28) : (big ? 28 : 22), height: n.current ? (big ? 36 : 28) : (big ? 28 : 22), display: "flex", alignItems: "center", justifyContent: "center", fontSize: n.current ? (big ? 20 : 16) : (big ? 14 : 12), background: n.current ? `${C.gold}22` : "transparent", border: n.current ? `2px solid ${C.gold}` : `1px solid ${C.border}`, borderRadius: "50%", boxShadow: n.current ? `0 0 12px ${C.gold}33` : "none" }}>
            {n.visited ? n.icon : "?"}
          </div>
          {(hov === n.id || n.current || big) && <div style={{ fontSize: big ? 8 : 7, color: n.current ? C.gold : C.textDim, whiteSpace: "nowrap", marginTop: 2, fontWeight: n.current ? 700 : 400 }}>{n.name} {big && <span style={{ color: C.textDim }}>Nv.{n.lv}</span>}</div>}
        </div>
      ))}
    </div>
  );
}

// ─── Character Sidebar ───

function CharSidebar() {
  return (
    <div style={{ width: 190, flexShrink: 0, background: C.bgPanel, borderRight: `1px solid ${C.border}`, display: "flex", flexDirection: "column", padding: 12, overflow: "auto", fontSize: 12 }}>
      <div style={{ width: 70, height: 70, margin: "0 auto 10px", background: `linear-gradient(135deg,${C.bgAlt},${C.bgInput})`, border: `2px solid ${C.borderGold}`, borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 32 }}>⚔</div>
      <div style={{ textAlign: "center", marginBottom: 10 }}>
        <div style={{ color: C.textBright, fontWeight: 700, fontSize: 14 }}>Kaira</div>
        <div style={{ color: C.gold, fontSize: 10, letterSpacing: 1 }}>GUERREIRO TRIBAL</div>
        <div style={{ color: C.textDim, fontSize: 10 }}>Nível 4 · Indígena</div>
      </div>
      <Bar current={32} max={40} color={C.redBright} label="Vida" />
      <Bar current={340} max={500} color={C.blueBright} label="XP" />
      <Divider />
      {[["Força","14 (+2)"],["Vigor","12 (+1)"],["Astúcia","10 (+0)"],["Sabedoria","8 (-1)"],["Presença","11 (+0)"],["Defesa","13"]].map(([k,v]) => (
        <div key={k} style={{ display: "flex", justifyContent: "space-between", fontSize: 11, color: C.textDim, marginBottom: 2 }}><span>{k}</span><span style={{ color: C.text }}>{v}</span></div>
      ))}
      <Divider />
      <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1, textTransform: "uppercase", marginBottom: 4 }}>Equipamento</div>
      <div style={{ fontSize: 11, lineHeight: 1.8, color: C.text }}>
        <div>🗡 Machadinha de Obsidiana</div>
        <div>🛡 Escudo de Casca</div>
        <div>👕 Couraça de Couro</div>
      </div>
    </div>
  );
}

// ─── VIEWS ───

// 1. Landing Page
function LandingPage({ onEnter }) {
  return (
    <div style={{ minHeight: "100vh", background: `linear-gradient(180deg, #0a0908 0%, ${C.bg} 40%, #1a1208 100%)`, display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", padding: 40, textAlign: "center" }}>
      <div style={{ fontSize: 11, letterSpacing: 4, color: C.goldDim, textTransform: "uppercase", marginBottom: 8 }}>Um MUD Moderno · Brasil Colonial · 1497</div>
      <h1 style={{ fontSize: 52, fontWeight: 300, color: C.gold, letterSpacing: 6, margin: "0 0 4px", fontFamily: font }}>ERAS DO BRASIL</h1>
      <div style={{ width: 120, height: 1, background: `linear-gradient(90deg, transparent, ${C.gold}, transparent)`, margin: "12px 0 20px" }} />
      <p style={{ fontSize: 16, color: C.text, maxWidth: 520, lineHeight: 1.7, marginBottom: 32 }}>
        A Raiz do Mundo entrou em colapso. Portais se abriram, criaturas despertaram, e a terra nunca mais será a mesma. Ninguém foi escolhido — você se destaca pelo que faz, não pelo que é.
      </p>
      <div style={{ display: "flex", gap: 12, marginBottom: 40 }}>
        {["Mundo Persistente 24/7", "Combate Tático D20", "Missões Competitivas", "Inimigos que Evoluem", "Zero Instalação"].map(f => (
          <div key={f} style={{ padding: "6px 14px", fontSize: 11, color: C.textDim, border: `1px solid ${C.border}`, borderRadius: 20 }}>{f}</div>
        ))}
      </div>
      <div style={{ display: "flex", flexDirection: "column", gap: 10, width: 300 }}>
        <input placeholder="Nome de aventureiro" style={{ padding: "10px 14px", fontSize: 14, background: C.bgInput, color: C.textBright, border: `1px solid ${C.border}`, borderRadius: 6, outline: "none", fontFamily: font, textAlign: "center" }} />
        <input placeholder="Senha" type="password" style={{ padding: "10px 14px", fontSize: 14, background: C.bgInput, color: C.textBright, border: `1px solid ${C.border}`, borderRadius: 6, outline: "none", fontFamily: font, textAlign: "center" }} />
        <button onClick={onEnter} style={{ padding: "12px", fontSize: 14, fontWeight: 700, background: `${C.gold}22`, color: C.gold, border: `2px solid ${C.gold}66`, borderRadius: 6, cursor: "pointer", letterSpacing: 2, fontFamily: font, marginTop: 4 }}>ENTRAR NO MUNDO</button>
        <button style={{ padding: "10px", fontSize: 12, background: "transparent", color: C.textDim, border: `1px solid ${C.border}`, borderRadius: 6, cursor: "pointer", fontFamily: font }}>Criar Personagem</button>
      </div>
      <div style={{ marginTop: 48, display: "flex", gap: 32, fontSize: 12, color: C.textDim }}>
        <span>🟢 47 aventureiros online</span>
        <span>☀ Tarde no mundo</span>
        <span>Temporada 1 — A Primeira Ruptura</span>
      </div>
    </div>
  );
}

// 2. Exploration
function ExplorationView() {
  const lines = [
    { t: "loc", v: "━━━ Vila de São Tomé ━━━" },
    { t: "desc", v: "O sol da tarde aquece as paredes de pau-a-pique. O Ferreiro Tomás bate seu martelo na forja. O cheiro de peixe seco e mandioca vem da taverna." },
    { t: "s", v: "" },
    { t: "pres", v: "Também estão aqui: Domingos (Missionário Nv.2), Iara (Caçadora Nv.3)" },
    { t: "npc", v: "Você vê: Ferreiro Tomás (NPC), Curandeira Naila (NPC), Guarda da Vila (NPC)" },
    { t: "s", v: "" },
    { t: "exit", v: "Saídas: [Norte] Floresta (~10 min) · [Leste] Rio (~5 min) · [Oeste] Pico (~15 min)" },
    { t: "s", v: "" },
    { t: "sys", v: "[Tarde] O Ferreiro Tomás guardou as ferramentas e foi à taverna." },
    { t: "sys", v: "[Tarde] Curandeira Naila colheu ervas no jardim da capela." },
    { t: "s", v: "" },
    { t: "cmd", v: "> olhar ao redor" },
    { t: "desc", v: "A praça central é simples: terra batida, um poço de pedra ao centro, e três construções — a forja, a taverna, e a capela. Um mural de madeira entalhada mostra um mapa rudimentar da região." },
  ];
  const colors = { loc: C.gold, desc: C.text, pres: C.blueBright, npc: C.greenBright, exit: C.gold, sys: C.textDim, cmd: C.goldDim, s: "transparent" };
  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
      <div style={{ flex: 1, overflow: "auto", padding: "16px 24px", display: "flex", flexDirection: "column", gap: 4 }}>
        {lines.map((l, i) => (
          <div key={i} style={{ color: colors[l.t], fontSize: l.t === "loc" ? 15 : l.t === "sys" ? 12 : 13, fontWeight: l.t === "loc" ? 700 : 400, fontStyle: l.t === "sys" ? "italic" : "normal", fontFamily: l.t === "cmd" ? mono : font, textAlign: l.t === "loc" ? "center" : "left", letterSpacing: l.t === "loc" ? 2 : 0, height: l.t === "s" ? 8 : "auto", lineHeight: l.t === "desc" ? 1.7 : 1.5 }}>{l.v}</div>
        ))}
      </div>
      <InputBar />
    </div>
  );
}

// 3. Combat
function CombatView() {
  return (
    <div style={{ flex: 1, display: "flex", flexDirection: "column", overflow: "auto" }}>
      {/* Initiative */}
      <div style={{ padding: "8px 16px", background: C.bgAlt, borderBottom: `1px solid ${C.border}`, display: "flex", alignItems: "center", gap: 8, fontSize: 11 }}>
        <span style={{ color: C.textDim, letterSpacing: 1, textTransform: "uppercase", fontSize: 10 }}>Iniciativa:</span>
        {[{ n: "Kaira", active: true }, { n: "Lobo Alfa", active: false }, { n: "Lobo", active: false }, { n: "Domingos", active: false }].map((c, i) => (
          <span key={i} style={{ padding: "3px 10px", background: c.active ? `${C.gold}22` : "transparent", border: `1px solid ${c.active ? C.gold + "66" : C.border}`, borderRadius: 3, color: c.active ? C.gold : C.textDim }}>{i + 1}. {c.n}</span>
        ))}
      </div>

      <div style={{ flex: 1, display: "flex", padding: 16, gap: 16, overflow: "auto" }}>
        {/* Enemies */}
        <div style={{ flex: 1 }}>
          <Panel title="Inimigos">
            <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
              {[{ name: "🐺 Lobo Alfa (Veterano)", hp: 18, maxHp: 22, lv: 4, status: "Enfurecido" },
                { name: "🐺 Lobo", hp: 6, maxHp: 8, lv: 2, status: null }].map((e, i) => (
                <div key={i} style={{ padding: 10, background: C.bgInput, borderRadius: 4, border: `1px solid ${C.border}` }}>
                  <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                    <span style={{ color: C.textBright, fontWeight: 600, fontSize: 13 }}>{e.name}</span>
                    <span style={{ color: C.textDim, fontSize: 11 }}>Nv.{e.lv}</span>
                  </div>
                  <Bar current={e.hp} max={e.maxHp} color={C.redBright} label="Vida" />
                  {e.status && <Tag color={C.redBright}>{e.status}</Tag>}
                </div>
              ))}
            </div>
          </Panel>

          {/* Combat Log */}
          <Panel title="Log de Combate" style={{ marginTop: 12 }}>
            <div style={{ fontSize: 12, display: "flex", flexDirection: "column", gap: 4, maxHeight: 160, overflow: "auto" }}>
              {[
                { c: C.gold, t: "⚔ Seu turno! Escolha uma ação." },
                { c: C.textDim, t: "Lobo atacou Domingos → 1d20+2 = 14 vs Defesa 12 → Acerto! 4 de dano." },
                { c: C.greenBright, t: "Domingos usou Iluminação Divina → Lobo Alfa fica Cego por 1 turno." },
                { c: C.blueBright, t: "Kaira atacou Lobo → 1d20+4 = 18 vs Defesa 11 → Acerto! 7 de dano (Cortante)." },
                { c: C.redBright, t: "Lobo Alfa uivou → Todos os aliados ganham +1 Dano por 2 turnos (Enfurecido)." },
              ].map((l, i) => <div key={i} style={{ color: l.c }}>{l.t}</div>)}
            </div>
          </Panel>
        </div>

        {/* Actions */}
        <div style={{ width: 260 }}>
          <Panel title="Suas Ações (Turno 3)">
            <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
              {[
                { icon: "⚔", name: "Atacar", desc: "1d8+2 Cortante (Machadinha)", target: true },
                { icon: "🛡", name: "Defender", desc: "+2 Defesa até próximo turno", target: false },
                { icon: "💥", name: "Golpe Tribal", desc: "1d10+3 Cortante · Recarga: 2 turnos", target: true },
                { icon: "🧪", name: "Usar Item", desc: "Poção de Cura (1d6+2 HP) · 2 restantes", target: false },
                { icon: "🏃", name: "Fugir", desc: "Teste Astúcia CD 12 · Perde turno se falhar", target: false },
              ].map((a, i) => (
                <div key={i} style={{ padding: "8px 10px", background: C.bgInput, borderRadius: 4, border: `1px solid ${C.border}`, cursor: "pointer", transition: "border-color .15s" }}
                  onMouseEnter={e => e.currentTarget.style.borderColor = C.borderGold}
                  onMouseLeave={e => e.currentTarget.style.borderColor = C.border}>
                  <div style={{ color: C.textBright, fontSize: 13, fontWeight: 600 }}>{a.icon} {a.name}</div>
                  <div style={{ color: C.textDim, fontSize: 11 }}>{a.desc}</div>
                  {a.target && <div style={{ color: C.gold, fontSize: 10, marginTop: 2 }}>Alvo: Lobo Alfa ▾</div>}
                </div>
              ))}
            </div>
          </Panel>

          <Panel title="D20 — Última Rolagem" style={{ marginTop: 12 }}>
            <div style={{ textAlign: "center" }}>
              <div style={{ fontSize: 42, color: C.gold, fontWeight: 700 }}>18</div>
              <div style={{ fontSize: 11, color: C.textDim }}>1d20 (14) + Força (+4) = 18</div>
              <div style={{ fontSize: 12, color: C.greenBright, marginTop: 4 }}>✓ Acerto! (vs Defesa 11)</div>
            </div>
          </Panel>
        </div>
      </div>
    </div>
  );
}

// 4. NPC Dialogue
function DialogueView() {
  const [selected, setSelected] = useState(null);
  return (
    <div style={{ flex: 1, display: "flex", padding: 16, gap: 16, overflow: "auto" }}>
      <div style={{ flex: 1 }}>
        <Panel>
          <div style={{ display: "flex", gap: 16, marginBottom: 16 }}>
            <div style={{ width: 80, height: 80, background: C.bgInput, border: `2px solid ${C.borderGold}`, borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 40, flexShrink: 0 }}>🔨</div>
            <div>
              <div style={{ color: C.gold, fontWeight: 700, fontSize: 16 }}>Ferreiro Tomás</div>
              <div style={{ color: C.textDim, fontSize: 12 }}>Artesão · Vila de São Tomé</div>
              <div style={{ display: "flex", gap: 6, marginTop: 6 }}><Tag>Afinidade: Amigável</Tag><Tag color={C.greenBright}>Missão Disponível</Tag></div>
            </div>
          </div>

          <div style={{ background: C.bgInput, borderRadius: 6, padding: 16, borderLeft: `3px solid ${C.borderGold}`, marginBottom: 16 }}>
            <div style={{ color: C.textBright, fontSize: 14, lineHeight: 1.7, fontStyle: "italic" }}>
              "Ah, mais um aventureiro! Olha, estou precisando de minério de ferro da Mina. Os Bandeirantes tomaram conta de lá, então não é tarefa fácil. Mas se trouxer 5 unidades, forjo uma boa espada pra você. O que me diz?"
            </div>
          </div>

          <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
            <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1, textTransform: "uppercase", marginBottom: 2 }}>Suas Respostas</div>
            {[
              { id: 1, text: '"Aceito o desafio. Quantos Bandeirantes tem na Mina?"', tag: "Aceitar missão", tagColor: C.greenBright },
              { id: 2, text: '"E se eu trouxer minério de qualidade superior?"', tag: "Negociar", tagColor: C.gold },
              { id: 3, text: '"Agora não posso, mas talvez depois."', tag: "Recusar", tagColor: C.textDim },
              { id: 4, text: '"Me conte sobre os Bandeirantes de Sangue."', tag: "Informação", tagColor: C.blueBright },
            ].map(opt => (
              <div key={opt.id} onClick={() => setSelected(opt.id)}
                style={{ padding: "10px 14px", background: selected === opt.id ? `${C.gold}11` : C.bgInput, border: `1px solid ${selected === opt.id ? C.gold + "44" : C.border}`, borderRadius: 4, cursor: "pointer", display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <span style={{ color: C.textBright, fontSize: 13 }}>{opt.text}</span>
                <Tag color={opt.tagColor}>{opt.tag}</Tag>
              </div>
            ))}
          </div>

          {selected && <div style={{ marginTop: 12, textAlign: "right" }}><Btn gold>Confirmar Resposta →</Btn></div>}
        </Panel>
      </div>

      <div style={{ width: 220 }}>
        <Panel title="Comércio Disponível">
          {[{ n: "Espada de Ferro", p: "120 UC", q: "Média" }, { n: "Kit de Reparo", p: "30 UC", q: "Comum" }, { n: "Poção de Cura", p: "25 UC", q: "Comum" }].map((item, i) => (
            <div key={i} style={{ padding: "6px 0", borderBottom: i < 2 ? `1px solid ${C.border}` : "none", fontSize: 12 }}>
              <div style={{ color: C.textBright }}>{item.n}</div>
              <div style={{ display: "flex", justifyContent: "space-between", color: C.textDim, fontSize: 11 }}><span>{item.q}</span><span style={{ color: C.gold }}>{item.p}</span></div>
            </div>
          ))}
        </Panel>
        <Panel title="Serviços" style={{ marginTop: 12 }}>
          <div style={{ fontSize: 12, display: "flex", flexDirection: "column", gap: 4 }}>
            <div style={{ color: C.text }}>⚒ Reparar equipamento</div>
            <div style={{ color: C.text }}>🔨 Forjar item (se tiver materiais)</div>
            <div style={{ color: C.text }}>📦 Ver estoque completo</div>
          </div>
        </Panel>
      </div>
    </div>
  );
}

// 5. Story Choice (Voting)
function StoryChoiceView() {
  const totalVotes = 47;
  const choiceA = 28, choiceB = 19;
  return (
    <div style={{ flex: 1, display: "flex", alignItems: "center", justifyContent: "center", padding: 32 }}>
      <div style={{ maxWidth: 640, width: "100%" }}>
        <div style={{ textAlign: "center", marginBottom: 24 }}>
          <Tag color={C.purpleBright}>Evento de Temporada — Fase de Decisão</Tag>
          <h2 style={{ color: C.gold, fontSize: 24, fontWeight: 400, margin: "12px 0 4px", fontFamily: font }}>O Destino da Vila de São Tomé</h2>
          <div style={{ color: C.textDim, fontSize: 13 }}>Os Bandeirantes marcham em direção à Vila. A comunidade deve decidir.</div>
        </div>
        <div style={{ background: C.bgPanel, border: `1px solid ${C.border}`, borderRadius: 8, padding: 20, marginBottom: 16 }}>
          <p style={{ color: C.text, fontSize: 14, lineHeight: 1.7, marginBottom: 16 }}>
            Os batedores trouxeram notícias: uma coluna de Bandeirantes de Sangue avança pela Floresta do Norte. Chegarão em 2 dias. A Vila não tem muralhas — mas tem o terreno e a Curandeira Naila conhece os caminhos secretos da floresta.
          </p>
          <div style={{ display: "flex", gap: 12 }}>
            {[{ label: "A — Defender a Vila", desc: "Construir barricadas, armar emboscadas na floresta. Arriscado, mas preserva a Vila.", votes: choiceA, color: C.greenBright },
              { label: "B — Evacuar para o Pico", desc: "Mover NPCs e recursos para o Pico da Neblina. Seguro, mas a Vila será destruída.", votes: choiceB, color: C.blueBright }].map((ch, i) => (
              <div key={i} style={{ flex: 1, padding: 16, background: C.bgInput, border: `1px solid ${C.border}`, borderRadius: 6, cursor: "pointer" }}>
                <div style={{ color: ch.color, fontWeight: 700, fontSize: 14, marginBottom: 4 }}>{ch.label}</div>
                <div style={{ color: C.textDim, fontSize: 12, marginBottom: 12 }}>{ch.desc}</div>
                <Bar current={ch.votes} max={totalVotes} color={ch.color} label={`${ch.votes} votos (${Math.round(ch.votes / totalVotes * 100)}%)`} height={8} />
              </div>
            ))}
          </div>
        </div>
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", fontSize: 12, color: C.textDim }}>
          <span>⏳ Votação encerra em: 1d 14h 22min</span>
          <span>{totalVotes} aventureiros votaram</span>
        </div>
      </div>
    </div>
  );
}

// 6. Gathering
function GatheringView() {
  return (
    <div style={{ flex: 1, display: "flex", padding: 16, gap: 16, overflow: "auto" }}>
      <div style={{ flex: 1 }}>
        <Panel title="Coleta — Floresta do Norte">
          <div style={{ color: C.text, fontSize: 13, lineHeight: 1.6, marginBottom: 16 }}>
            Você se ajoelha diante de um veio de Madeira de Ferro na base de uma árvore centenária. As fibras brilham com um tom avermelhado sob a luz filtrada pela copa.
          </div>
          <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            {[{ res: "🪵 Madeira de Ferro", qty: "3/5 restantes", skill: "Carpintaria Nv.1", cd: "CD 10", time: "~30s", rarity: "Comum" },
              { res: "🌿 Erva Medicinal", qty: "4/4 restantes", skill: "Herborismo Nv.1", cd: "CD 8", time: "~15s", rarity: "Comum" },
              { res: "🌙 Erva-Lua", qty: "1/2 restantes", skill: "Herborismo Nv.2", cd: "CD 14", time: "~45s", rarity: "Incomum" },
              { res: "🍄 Cogumelo Luminescente", qty: "2/3 restantes", skill: "Herborismo Nv.1", cd: "CD 10", time: "~20s", rarity: "Comum" },
            ].map((r, i) => (
              <div key={i} style={{ padding: 10, background: C.bgInput, border: `1px solid ${C.border}`, borderRadius: 4, display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                <div>
                  <div style={{ color: C.textBright, fontSize: 13, fontWeight: 600 }}>{r.res} <Tag>{r.rarity}</Tag></div>
                  <div style={{ color: C.textDim, fontSize: 11, marginTop: 2 }}>{r.skill} · {r.cd} · {r.time}</div>
                </div>
                <div style={{ textAlign: "right" }}>
                  <div style={{ color: C.textDim, fontSize: 11 }}>{r.qty}</div>
                  <Btn small gold>Coletar</Btn>
                </div>
              </div>
            ))}
          </div>
        </Panel>

        <Panel title="Resultado da Coleta" style={{ marginTop: 12 }}>
          <div style={{ fontSize: 13 }}>
            <div style={{ color: C.gold, marginBottom: 6 }}>🎲 Teste de Carpintaria: 1d20 (12) + Proficiência (+2) = 14 vs CD 10</div>
            <div style={{ color: C.greenBright, marginBottom: 4 }}>✓ Sucesso! Você obteve: 1× Madeira de Ferro (Qualidade Média)</div>
            <div style={{ color: C.textDim, fontSize: 11 }}>+1 XP de Carpintaria · Recurso regenera no próximo período</div>
          </div>
        </Panel>
      </div>

      <div style={{ width: 220 }}>
        <Panel title="Seus Recursos">
          {[["🪵 Madeira de Ferro", "7"], ["🌿 Erva Medicinal", "12"], ["⛏ Minério de Ferro", "2"], ["🌙 Erva-Lua", "1"], ["💰 UC (Ouro)", "340"]].map(([n, q], i) => (
            <div key={i} style={{ display: "flex", justifyContent: "space-between", fontSize: 12, color: C.text, padding: "4px 0", borderBottom: `1px solid ${C.border}22` }}><span>{n}</span><span style={{ color: C.gold }}>{q}</span></div>
          ))}
        </Panel>
        <Panel title="Proficiências" style={{ marginTop: 12 }}>
          {[["Carpintaria", "Nv.2", 65], ["Herborismo", "Nv.1", 30], ["Mineração", "Nv.1", 10]].map(([n, lv, xp], i) => (
            <div key={i} style={{ marginBottom: 6 }}>
              <div style={{ display: "flex", justifyContent: "space-between", fontSize: 11, color: C.text }}><span>{n}</span><span style={{ color: C.gold }}>{lv}</span></div>
              <div style={{ height: 3, background: C.bgInput, borderRadius: 2, overflow: "hidden" }}>
                <div style={{ width: `${xp}%`, height: "100%", background: C.blueBright, borderRadius: 2 }} />
              </div>
            </div>
          ))}
        </Panel>
      </div>
    </div>
  );
}

// 7. Crafting
function CraftingView() {
  return (
    <div style={{ flex: 1, display: "flex", padding: 16, gap: 16, overflow: "auto" }}>
      <div style={{ width: 220 }}>
        <Panel title="Receitas Conhecidas">
          {[{ n: "Machadinha de Obsidiana", cat: "Arma", unlocked: true },
            { n: "Poção de Cura", cat: "Alquimia", unlocked: true },
            { n: "Poção Anti-Fadiga", cat: "Alquimia", unlocked: true },
            { n: "Armadura de Couro Reforçado", cat: "Coureira", unlocked: true },
            { n: "Colar de Dentes", cat: "Artesanato", unlocked: false },
          ].map((r, i) => (
            <div key={i} style={{ padding: "8px 0", borderBottom: `1px solid ${C.border}22`, cursor: "pointer", opacity: r.unlocked ? 1 : 0.4 }}>
              <div style={{ color: C.textBright, fontSize: 12 }}>{r.n}</div>
              <div style={{ color: C.textDim, fontSize: 10 }}>{r.cat}{!r.unlocked && " · 🔒 Bloqueada"}</div>
            </div>
          ))}
        </Panel>
      </div>

      <div style={{ flex: 1 }}>
        <Panel title="Forjar: Machadinha de Obsidiana">
          <div style={{ display: "flex", gap: 16, marginBottom: 16 }}>
            <div style={{ width: 64, height: 64, background: C.bgInput, border: `2px solid ${C.borderGold}`, borderRadius: 8, display: "flex", alignItems: "center", justifyContent: "center", fontSize: 32 }}>🪓</div>
            <div>
              <div style={{ color: C.gold, fontWeight: 700, fontSize: 15 }}>Machadinha de Obsidiana</div>
              <div style={{ color: C.textDim, fontSize: 12 }}>Arma · Dano 1d8+1 Cortante · Qualidade depende da rolagem</div>
              <div style={{ display: "flex", gap: 6, marginTop: 4 }}><Tag>Carpintaria Nv.2</Tag><Tag color={C.greenBright}>CD 12</Tag></div>
            </div>
          </div>

          <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1, textTransform: "uppercase", marginBottom: 6 }}>Ingredientes</div>
          <div style={{ display: "flex", flexDirection: "column", gap: 4, marginBottom: 16 }}>
            {[{ n: "Obsidiana Bruta", need: 2, have: 3, ok: true }, { n: "Madeira de Ferro", need: 1, have: 7, ok: true }, { n: "Couro Curtido", need: 1, have: 0, ok: false }].map((ing, i) => (
              <div key={i} style={{ display: "flex", justifyContent: "space-between", padding: "6px 10px", background: C.bgInput, borderRadius: 4, border: `1px solid ${ing.ok ? C.border : C.red + "44"}` }}>
                <span style={{ color: C.text, fontSize: 12 }}>{ing.n}</span>
                <span style={{ color: ing.ok ? C.greenBright : C.redBright, fontSize: 12 }}>{ing.have}/{ing.need}</span>
              </div>
            ))}
          </div>

          <div style={{ padding: 12, background: `${C.red}11`, border: `1px solid ${C.red}33`, borderRadius: 4, fontSize: 12, color: C.redBright, marginBottom: 12 }}>
            ⚠ Falta 1× Couro Curtido. Obtenha com a proficiência Coureira ou compre de outro jogador.
          </div>

          <div style={{ display: "flex", gap: 8 }}>
            <Btn gold style={{ opacity: 0.5, cursor: "not-allowed", flex: 1 }}>🔨 Forjar (~2 min)</Btn>
          </div>

          <Divider />
          <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1, textTransform: "uppercase", marginBottom: 6 }}>Qualidade Prevista</div>
          <div style={{ display: "flex", gap: 8, fontSize: 11 }}>
            {[["Baixa", "≤8", C.textDim], ["Média", "9-15", C.text], ["Alta", "16-19", C.blueBright], ["Excepcional", "20", C.gold]].map(([q, r, c], i) => (
              <div key={i} style={{ flex: 1, textAlign: "center", padding: "6px 0", background: C.bgInput, borderRadius: 4, border: `1px solid ${C.border}` }}>
                <div style={{ color: c, fontWeight: 600 }}>{q}</div>
                <div style={{ color: C.textDim, fontSize: 10 }}>D20: {r}</div>
              </div>
            ))}
          </div>
        </Panel>
      </div>
    </div>
  );
}

// 8. Inventory
function InventoryView() {
  const items = [
    { n: "Machadinha de Obsidiana", t: "Arma", r: "Incomum", eq: true, dur: "85%", desc: "1d8+1 Cortante" },
    { n: "Escudo de Casca", t: "Escudo", r: "Comum", eq: true, dur: "70%", desc: "+1 Defesa" },
    { n: "Couraça de Couro", t: "Armadura", r: "Comum", eq: true, dur: "60%", desc: "+2 Defesa" },
    { n: "Poção de Cura", t: "Consumível", r: "Comum", eq: false, qty: 3, desc: "Restaura 1d6+2 PV" },
    { n: "Minério de Ferro", t: "Material", r: "Comum", eq: false, qty: 2, desc: "Componente de crafting" },
    { n: "Madeira de Ferro", t: "Material", r: "Comum", eq: false, qty: 7, desc: "Componente de crafting" },
    { n: "Erva-Lua", t: "Material", r: "Incomum", eq: false, qty: 1, desc: "Reduz Fadiga Espiritual" },
    { n: "Ouro da Raiz", t: "Quest", r: "Raro", eq: false, qty: 1, desc: "Moeda de Classe Universal", quest: true },
  ];
  const rarColor = { Comum: C.textDim, Incomum: C.greenBright, Raro: C.blueBright };
  return (
    <div style={{ flex: 1, padding: 16, overflow: "auto" }}>
      <div style={{ display: "flex", gap: 16 }}>
        <div style={{ flex: 1 }}>
          <Panel title={`Inventário (${items.length}/20 slots)`}>
            <div style={{ display: "flex", flexDirection: "column", gap: 4 }}>
              {items.map((item, i) => (
                <div key={i} style={{ display: "flex", justifyContent: "space-between", alignItems: "center", padding: "8px 10px", background: C.bgInput, borderRadius: 4, border: `1px solid ${item.quest ? C.blueBright + "33" : C.border}`, cursor: "pointer" }}>
                  <div>
                    <span style={{ color: C.textBright, fontSize: 13 }}>{item.n} </span>
                    <Tag color={rarColor[item.r]}>{item.r}</Tag>
                    {item.eq && <Tag color={C.gold}>Equipado</Tag>}
                    {item.quest && <Tag color={C.purpleBright}>Quest</Tag>}
                  </div>
                  <div style={{ color: C.textDim, fontSize: 11, textAlign: "right" }}>
                    {item.dur && <div>Durabilidade: {item.dur}</div>}
                    {item.qty && <div>Qtd: {item.qty}</div>}
                  </div>
                </div>
              ))}
            </div>
          </Panel>
        </div>
        <div style={{ width: 200 }}>
          <Panel title="Ouro"><div style={{ color: C.gold, fontSize: 24, fontWeight: 700, textAlign: "center" }}>340 UC</div></Panel>
          <Panel title="Peso" style={{ marginTop: 12 }}><Bar current={14} max={20} color={C.textDim} label="Carga" /></Panel>
          <Panel title="Ações" style={{ marginTop: 12 }}>
            <div style={{ display: "flex", flexDirection: "column", gap: 4 }}>
              <Btn small>Ordenar por tipo</Btn>
              <Btn small>Ordenar por raridade</Btn>
            </div>
          </Panel>
        </div>
      </div>
    </div>
  );
}

// 9. Full Map View
function MapView() {
  const [sel, setSel] = useState(null);
  const regions = [
    { n: "Mata Costeira", lv: "1-5", status: "Atual", desc: "Litoral, Mata Atlântica. Onde a Ruptura começou." },
    { n: "Sertão Distorcido", lv: "5-10", status: "Temporada 2", desc: "Caatinga distorcida pela Raiz." },
    { n: "Serra dos Ecos", lv: "10-15", status: "Temporada 3", desc: "Montanhas com portais temporais." },
    { n: "Pantanal Vivo", lv: "15-20", status: "Temporada 4", desc: "Terreno vivo e mutante." },
    { n: "Coração da Raiz", lv: "20+", status: "Futuro", desc: "Centro do continente. Endgame." },
  ];
  return (
    <div style={{ flex: 1, display: "flex", padding: 16, gap: 16, overflow: "auto" }}>
      <div style={{ flex: 1 }}>
        <Panel title="Mapa Regional — Mata Costeira (Nv. 1-5)">
          <MiniMap size="big" onSelect={n => setSel(n)} />
          {sel && (
            <div style={{ marginTop: 12, padding: 10, background: C.bgInput, borderRadius: 4, border: `1px solid ${C.border}` }}>
              <div style={{ color: C.gold, fontWeight: 700, fontSize: 14 }}>{sel.icon} {sel.name}</div>
              <div style={{ color: C.textDim, fontSize: 12 }}>Nível: {sel.lv} · {sel.visited ? "Explorado" : "Inexplorado"}</div>
              {sel.visited && !sel.current && <Btn small gold style={{ marginTop: 8 }}>Viajar para cá</Btn>}
              {sel.current && <div style={{ color: C.greenBright, fontSize: 11, marginTop: 4 }}>📍 Você está aqui</div>}
            </div>
          )}
        </Panel>
      </div>
      <div style={{ width: 240 }}>
        <Panel title="Mapa-Múndi — Regiões">
          <div style={{ display: "flex", flexDirection: "column", gap: 6 }}>
            {regions.map((r, i) => (
              <div key={i} style={{ padding: 8, background: C.bgInput, borderRadius: 4, border: `1px solid ${i === 0 ? C.gold + "44" : C.border}`, opacity: i === 0 ? 1 : 0.5 }}>
                <div style={{ display: "flex", justifyContent: "space-between" }}>
                  <span style={{ color: i === 0 ? C.gold : C.textDim, fontWeight: 600, fontSize: 12 }}>{r.n}</span>
                  <Tag color={i === 0 ? C.greenBright : C.textDim}>{r.status}</Tag>
                </div>
                <div style={{ color: C.textDim, fontSize: 11 }}>Nv. {r.lv}</div>
                {i === 0 && <div style={{ color: C.textDim, fontSize: 10, marginTop: 2 }}>{r.desc}</div>}
              </div>
            ))}
          </div>
        </Panel>
        <Panel title="Legendas" style={{ marginTop: 12 }}>
          <div style={{ fontSize: 11, color: C.textDim, display: "flex", flexDirection: "column", gap: 3 }}>
            <div>🏠 Hub / Vila segura</div><div>🌲 Floresta / Exploração</div><div>🌊 Água / Maré</div>
            <div>🏔 Montanha / Fadiga</div><div>🔥 Ruínas / Terror</div><div>🌀 Boss / Evento</div>
            <div style={{ marginTop: 4, color: C.text }}>● Explorado</div><div>○ Inexplorado</div>
          </div>
        </Panel>
      </div>
    </div>
  );
}

// 10. Quest Log
function QuestLogView() {
  return (
    <div style={{ flex: 1, padding: 16, overflow: "auto" }}>
      <Panel title="Diário de Missões">
        <div style={{ display: "flex", gap: 8, marginBottom: 16 }}>
          <Btn small active>Ativas (3)</Btn><Btn small>Completas (5)</Btn><Btn small>Falhadas (1)</Btn>
        </div>
        <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
          {[
            { n: "O Pedido do Ferreiro", giver: "Ferreiro Tomás", prog: "0/5 Minérios", zone: "Mina de Ouro", timer: "2d 14h restantes", comp: false, first: null },
            { n: "O Caçador que Não Voltou", giver: "Pajé Aruã", prog: "Etapa 1/3 — Investigar", zone: "Floresta do Norte", timer: null, comp: false, first: null },
            { n: "Peles para a Vila", giver: "Guarda da Vila", prog: "7/10 Peles de Lobo", zone: "Floresta do Norte", timer: "6h restantes", comp: false, first: "Domingos (1° lugar)" },
          ].map((q, i) => (
            <div key={i} style={{ padding: 12, background: C.bgInput, borderRadius: 6, border: `1px solid ${C.border}` }}>
              <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                <span style={{ color: C.gold, fontWeight: 700, fontSize: 14 }}>{q.n}</span>
                {q.timer && <Tag color={C.redBright}>⏳ {q.timer}</Tag>}
              </div>
              <div style={{ fontSize: 12, color: C.textDim, marginBottom: 6 }}>Dado por: {q.giver} · Local: {q.zone}</div>
              <Bar current={parseInt(q.prog)} max={parseInt(q.prog.split("/")[1])} color={C.gold} label={q.prog} height={5} />
              {q.first && <div style={{ fontSize: 11, color: C.blueBright, marginTop: 4 }}>🏆 {q.first} — recompensa parcial disponível</div>}
            </div>
          ))}
        </div>
      </Panel>
    </div>
  );
}

// Input bar shared
function InputBar() {
  return (
    <div style={{ padding: "10px 16px", background: C.bgPanel, borderTop: `1px solid ${C.border}` }}>
      <div style={{ display: "flex", gap: 8, marginBottom: 6 }}>
        <input placeholder="O que deseja fazer?" style={{ flex: 1, padding: "8px 12px", fontSize: 13, background: C.bgInput, color: C.textBright, border: `1px solid ${C.border}`, borderRadius: 4, outline: "none", fontFamily: mono }} />
        <Btn gold>ENVIAR</Btn>
      </div>
      <div style={{ display: "flex", gap: 4, flexWrap: "wrap" }}>
        {["💬 Falar", "🔍 Explorar", "🛏 Descansar", "⚔ Atacar", "📦 Coletar", "🏃 Fugir"].map(b => <Btn key={b} small>{b}</Btn>)}
      </div>
    </div>
  );
}

// ─── MAIN APP ───

const VIEWS = [
  { id: "explore", label: "Exploração", icon: "🌍" },
  { id: "combat", label: "Combate", icon: "⚔" },
  { id: "dialogue", label: "NPC / Diálogo", icon: "💬" },
  { id: "story", label: "Escolha / Votação", icon: "📖" },
  { id: "gather", label: "Coleta", icon: "🌿" },
  { id: "craft", label: "Crafting", icon: "🔨" },
  { id: "inventory", label: "Inventário", icon: "📦" },
  { id: "map", label: "Mapa", icon: "🗺" },
  { id: "quests", label: "Diário", icon: "📜" },
];

export default function ErasDoBrasilComplete() {
  const [view, setView] = useState("landing");
  const [eventIdx, setEventIdx] = useState(0);
  const events = ["☀ Tarde — sol se põe sobre a Mata Costeira", "🐺 Lobo Veterano avistado na Floresta", "⚒ Ferreiro Tomás abriu a forja", "📜 Domingos completou uma missão", "🌿 Recursos regeneraram no Rio das Marés"];

  useEffect(() => { const t = setInterval(() => setEventIdx(i => (i + 1) % events.length), 4000); return () => clearInterval(t); }, []);

  if (view === "landing") return <LandingPage onEnter={() => setView("explore")} />;

  const ViewComponent = { explore: ExplorationView, combat: CombatView, dialogue: DialogueView, story: StoryChoiceView, gather: GatheringView, craft: CraftingView, inventory: InventoryView, map: MapView, quests: QuestLogView }[view];

  return (
    <div style={{ width: "100%", height: "100vh", background: C.bg, fontFamily: font, color: C.text, display: "flex", flexDirection: "column", overflow: "hidden" }}>
      <style>{`@import url('https://fonts.googleapis.com/css2?family=Crimson+Pro:ital,wght@0,300;0,400;0,600;0,700;1,400&family=Courier+Prime&display=swap');*{box-sizing:border-box;margin:0;padding:0}::-webkit-scrollbar{width:6px}::-webkit-scrollbar-track{background:${C.bgInput}}::-webkit-scrollbar-thumb{background:${C.border};border-radius:3px}`}</style>

      {/* Top bar */}
      <div style={{ height: 36, background: C.bgPanel, borderBottom: `1px solid ${C.border}`, display: "flex", alignItems: "center", justifyContent: "space-between", padding: "0 12px", fontSize: 12, flexShrink: 0 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
          <span style={{ color: C.gold, fontWeight: 700, letterSpacing: 2, fontSize: 12 }}>ERAS DO BRASIL</span>
          <span style={{ color: C.textDim }}>·</span>
          <span style={{ color: C.textDim }}>Mata Costeira</span>
        </div>
        <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
          <span style={{ color: C.green, fontSize: 11 }}>● 47 online</span>
          <span style={{ color: C.gold, fontSize: 11 }}>☀ Tarde</span>
          <span style={{ color: C.textDim, fontSize: 11 }}>T1 — A Primeira Ruptura</span>
        </div>
      </div>

      {/* Nav tabs */}
      <div style={{ display: "flex", background: C.bgAlt, borderBottom: `1px solid ${C.border}`, padding: "0 8px", flexShrink: 0, overflow: "auto" }}>
        {VIEWS.map(v => (
          <button key={v.id} onClick={() => setView(v.id)}
            style={{ padding: "8px 14px", fontSize: 11, background: view === v.id ? C.bgPanel : "transparent", color: view === v.id ? C.gold : C.textDim, border: "none", borderBottom: view === v.id ? `2px solid ${C.gold}` : "2px solid transparent", cursor: "pointer", fontFamily: font, whiteSpace: "nowrap", transition: "all .15s" }}>
            {v.icon} {v.label}
          </button>
        ))}
        <div style={{ flex: 1 }} />
        <button onClick={() => setView("landing")} style={{ padding: "8px 14px", fontSize: 11, background: "transparent", color: C.textDim, border: "none", cursor: "pointer", fontFamily: font }}>🚪 Sair</button>
      </div>

      {/* Main */}
      <div style={{ flex: 1, display: "flex", overflow: "hidden" }}>
        {["explore", "combat", "gather"].includes(view) && <CharSidebar />}
        <ViewComponent />
        {view === "explore" && (
          <div style={{ width: 220, flexShrink: 0, background: C.bgPanel, borderLeft: `1px solid ${C.border}`, padding: 12, overflow: "auto" }}>
            <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1.5, textTransform: "uppercase", marginBottom: 8, fontWeight: 600 }}>Mata Costeira</div>
            <MiniMap />
            <Divider />
            <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1.5, textTransform: "uppercase", marginBottom: 6, fontWeight: 600 }}>Missões</div>
            {[["O Pedido do Ferreiro", "0/5"], ["Peles para a Vila", "7/10"]].map(([n, p], i) => (
              <div key={i} style={{ fontSize: 11, padding: "4px 0", color: C.text }}>{n} <span style={{ color: C.gold }}>{p}</span></div>
            ))}
            <Divider />
            <div style={{ fontSize: 10, color: C.textDim, letterSpacing: 1.5, textTransform: "uppercase", marginBottom: 6, fontWeight: 600 }}>Na zona</div>
            <div style={{ fontSize: 11 }}>
              <div><span style={{ color: C.green }}>●</span> <span style={{ color: C.blueBright }}>Domingos</span> <span style={{ color: C.textDim }}>Nv.2</span></div>
              <div><span style={{ color: C.green }}>●</span> <span style={{ color: C.blueBright }}>Iara</span> <span style={{ color: C.textDim }}>Nv.3</span></div>
            </div>
          </div>
        )}
      </div>

      {/* Bottom */}
      <div style={{ height: 28, background: C.bgPanel, borderTop: `1px solid ${C.border}`, display: "flex", alignItems: "center", padding: "0 12px", fontSize: 11, color: C.textDim, flexShrink: 0 }}>
        <span style={{ color: C.borderGold, marginRight: 8, fontSize: 10, letterSpacing: 1 }}>MUNDO</span>
        {events[eventIdx]}
      </div>
    </div>
  );
}
