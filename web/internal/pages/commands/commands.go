package commands

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/TheLazyTurtle33/sea-core/web/internal/botclient"
)

// --- Cache ---

type cache struct {
	mu        sync.Mutex
	data      []byte
	fetchedAt time.Time
	ttl       time.Duration
}

var commandCache = &cache{ttl: 5 * time.Minute}

func (c *cache) get() ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.data == nil || time.Since(c.fetchedAt) > c.ttl {
		return nil, false
	}
	return c.data, true
}

func (c *cache) set(data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = data
	c.fetchedAt = time.Now()
}

// --- Handler ---

func Page(w http.ResponseWriter, r *http.Request) {
	commandsJSON, ok := commandCache.get()
	if !ok {
		var err error
		commandsJSON, err = botclient.GetCommandsJson()
		if err != nil {
			log.Printf("Error getting commands JSON: %v", err)
			commandsJSON = []byte("[]")
		} else {
			commandCache.set(commandsJSON)
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(renderPage(string(commandsJSON))))
}

func renderPage(commandsJSON string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>lazybot33 — commands</title>
  <link rel="preconnect" href="https://fonts.googleapis.com" />
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
  <link href="https://fonts.googleapis.com/css2?family=Cinzel:wght@400;600;700&family=Crimson+Pro:ital,wght@0,300;0,400;1,300;1,400&display=swap" rel="stylesheet" />
  <style>
    :root {
      --abyss:        #04090f;
      --deep:         #060d18;
      --card:         #071120;
      --card-hover:   #0c1a2e;

      --biolum:       #3ef5c8;
      --biolum-dim:   #1a7a61;
      --biolum-glow:  #3ef5c830;
      --biolum-faint: #3ef5c812;

      --teal:         #22d4b0;

      --text:         #c8dfe8;
      --text-dim:     #6a90a8;
      --text-faint:   #2e4d60;

      --border:       #0f2035;
      --border-glow:  #1a4060;

      --inactive-bg:     #060e18;
      --inactive-text:   #1e3448;
      --inactive-border: #0a1826;
    }

    *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
    html { scroll-behavior: smooth; }

    body {
      background: var(--abyss);
      color: var(--text);
      font-family: 'Crimson Pro', Georgia, serif;
      min-height: 100vh;
      overflow-x: hidden;
      position: relative;
    }

    body::before {
      content: '';
      position: fixed;
      inset: 0;
      background:
        radial-gradient(ellipse 60% 40% at 20% 10%, #0a1f3520 0%, transparent 70%),
        radial-gradient(ellipse 40% 60% at 80% 80%, #071828 0%, transparent 60%),
        linear-gradient(180deg, #060d18 0%, #04090f 100%);
      pointer-events: none;
      z-index: 0;
    }

    #particles {
      position: fixed;
      inset: 0;
      pointer-events: none;
      z-index: 1;
      overflow: hidden;
    }

    .particle {
      position: absolute;
      background: var(--biolum);
      border-radius: 50%;
      opacity: 0;
      animation: drift linear infinite;
    }

    @keyframes drift {
      0%   { transform: translateY(-10px) translateX(0px);  opacity: 0; }
      10%  { opacity: 0.6; }
      90%  { opacity: 0.3; }
      100% { transform: translateY(100vh) translateX(30px); opacity: 0; }
    }

    header, main, footer { position: relative; z-index: 2; }

    /* ---- Header ---- */
    header {
      padding: 3rem 2rem 0;
      text-align: center;
    }

    .header-inner { max-width: 860px; margin: 0 auto; }

    .depth-marker {
      font-family: 'Cinzel', serif;
      font-size: 0.6rem;
      letter-spacing: 0.25em;
      text-transform: uppercase;
      color: var(--biolum-dim);
      margin-bottom: 1.2rem;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 1rem;
    }

    .depth-marker::before,
    .depth-marker::after {
      content: '';
      flex: 1;
      max-width: 80px;
      height: 1px;
    }
    .depth-marker::before { background: linear-gradient(90deg, transparent, var(--biolum-dim)); }
    .depth-marker::after  { background: linear-gradient(90deg, var(--biolum-dim), transparent); }

    h1 {
      font-family: 'Cinzel', serif;
      font-weight: 700;
      font-size: clamp(2rem, 6vw, 3.5rem);
      color: var(--text-dim);
      letter-spacing: 0.06em;
      line-height: 1;
      text-shadow: 0 0 40px var(--biolum-faint);
    }

    h1 em {
      font-style: normal;
      color: var(--biolum-dim);
      text-shadow: 0 0 20px var(--biolum-faint);
    }

    .subtitle {
      margin-top: 0.75rem;
      font-size: 1.1rem;
      color: var(--text-dim);
      font-style: italic;
      font-weight: 300;
    }

    .status-row {
      margin-top: 1.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.6rem;
      font-family: 'Cinzel', serif;
      font-size: 0.65rem;
      letter-spacing: 0.1em;
      color: var(--text-faint);
    }

    .status-dot {
      width: 7px; height: 7px;
      border-radius: 50%;
      background: var(--biolum);
      box-shadow: 0 0 10px var(--biolum), 0 0 20px var(--biolum-glow);
      animation: glow-pulse 3s ease-in-out infinite;
    }

    @keyframes glow-pulse {
      0%, 100% { box-shadow: 0 0 10px var(--biolum), 0 0 20px var(--biolum-glow); }
      50%       { box-shadow: 0 0 4px var(--biolum-dim); }
    }

    .header-rule {
      margin-top: 2.5rem;
      width: 100%;
      height: 1px;
      background: linear-gradient(90deg, transparent, var(--border-glow) 30%, var(--border-glow) 70%, transparent);
      position: relative;
    }

    .header-rule::before {
      content: '✦';
      position: absolute;
      left: 50%;
      top: 50%;
      transform: translate(-50%, -50%);
      background: var(--abyss);
      padding: 0 0.75rem;
      color: var(--biolum-dim);
      font-size: 0.7rem;
    }

    /* ---- Main ---- */
    main {
      max-width: 860px;
      margin: 0 auto;
      padding: 2.5rem 2rem 5rem;
    }

    .section-label {
      font-family: 'Cinzel', serif;
      font-size: 0.6rem;
      letter-spacing: 0.2em;
      text-transform: uppercase;
      color: var(--biolum-dim);
      margin-bottom: 1.2rem;
      display: flex;
      align-items: center;
      gap: 0.75rem;
    }

    .section-label::after {
      content: '';
      flex: 1;
      height: 1px;
      background: linear-gradient(90deg, var(--border-glow), transparent);
    }

    /* ---- Cards ---- */
    .commands-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
    }

    .cmd-card {
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 6px;
      padding: 1.25rem 1.5rem;
      position: relative;
      overflow: hidden;
      transition: background 0.2s, border-color 0.2s, transform 0.2s;
      animation: rise 0.5s ease both;
    }

    .cmd-card:not(.inactive) { cursor: pointer; }

    /* Clipboard toast */
    .copy-toast {
      position: fixed;
      bottom: 2rem;
      left: 50%;
      transform: translateX(-50%) translateY(10px);
      background: #0c1e35;
      border: 1px solid var(--biolum-dim);
      color: var(--biolum);
      font-family: 'Crimson Pro', serif;
      font-style: italic;
      font-size: 0.95rem;
      padding: 0.5rem 1.25rem;
      border-radius: 4px;
      box-shadow: 0 0 20px var(--biolum-faint);
      opacity: 0;
      pointer-events: none;
      transition: opacity 0.2s, transform 0.2s;
      z-index: 999;
    }

    .copy-toast.show {
      opacity: 1;
      transform: translateX(-50%) translateY(0);
    }

    @keyframes rise {
      from { opacity: 0; transform: translateY(12px); }
      to   { opacity: 1; transform: translateY(0); }
    }

    .cmd-card::before {
      content: '';
      position: absolute;
      left: 0; top: 15%; bottom: 15%;
      width: 2px;
      background: linear-gradient(180deg, transparent, var(--biolum-dim), transparent);
      border-radius: 0 2px 2px 0;
      transition: opacity 0.2s;
      opacity: 0.6;
    }

    .cmd-card::after {
      content: '';
      position: absolute;
      inset: 0;
      background: radial-gradient(ellipse 60% 80% at 0% 50%, var(--biolum-faint), transparent 70%);
      opacity: 0;
      transition: opacity 0.3s;
      pointer-events: none;
    }

    .cmd-card:not(.inactive):hover {
      background: var(--card-hover);
      border-color: var(--border-glow);
      transform: translateX(3px);
    }

    .cmd-card:not(.inactive):hover::before { opacity: 1; }
    .cmd-card:not(.inactive):hover::after  { opacity: 1; }

    .cmd-card.inactive {
      background: var(--inactive-bg);
      border-color: var(--inactive-border);
      opacity: 0.45;
    }

    .cmd-card.inactive::before { opacity: 0.2; }

    .cmd-card.inactive .cmd-name,
    .cmd-card.inactive .cmd-usage,
    .cmd-card.inactive .cmd-desc { color: var(--inactive-text); }

    .cmd-card.inactive .cmd-trigger {
      background: #060e18;
      color: var(--inactive-text);
      border-color: var(--inactive-border);
    }

    .cmd-layout {
      display: grid;
      grid-template-columns: 1fr auto;
      grid-template-rows: auto auto;
      gap: 0.4rem 1.5rem;
      align-items: start;
    }

    .cmd-header {
      display: flex;
      align-items: baseline;
      gap: 0.85rem;
      flex-wrap: wrap;
    }

    .cmd-name {
      font-family: 'Cinzel', serif;
      font-weight: 600;
      font-size: 1rem;
      color: #e8f4f8;
      letter-spacing: 0.04em;
    }

    .cmd-usage {
      font-size: 0.85rem;
      color: var(--teal);
      font-style: italic;
      font-weight: 300;
    }

    .cmd-desc {
      font-size: 1.05rem;
      color: #7da8c0;
      font-style: italic;
      font-weight: 300;
      line-height: 1.6;
      grid-column: 1;
    }

    .cmd-triggers {
      grid-column: 2;
      grid-row: 1 / 3;
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: 0.3rem;
      justify-content: center;
    }

    .cmd-trigger {
      font-family: 'Crimson Pro', Georgia, serif;
      font-size: 0.8rem;
      font-style: italic;
      letter-spacing: 0.04em;
      background: #091828;
      border: 1px solid var(--border-glow);
      color: var(--biolum-dim);
      padding: 0.2rem 0.65rem;
      border-radius: 3px;
      white-space: nowrap;
      transition: color 0.15s, background 0.15s, border-color 0.15s;
    }

    .cmd-card:not(.inactive):hover .cmd-trigger {
      color: var(--biolum);
      background: #0a1f35;
      border-color: var(--biolum-dim);
    }

    .inactive-badge {
      font-family: 'Cinzel', serif;
      font-size: 0.55rem;
      letter-spacing: 0.12em;
      text-transform: uppercase;
      color: var(--text-faint);
      border: 1px solid var(--inactive-border);
      padding: 0.15rem 0.45rem;
      border-radius: 2px;
      margin-left: 0.25rem;
    }

    .inactive-section { margin-top: 3rem; }

    .empty {
      text-align: center;
      padding: 5rem 2rem;
      color: var(--text-faint);
      font-style: italic;
      font-size: 1rem;
    }

    /* ---- Footer ---- */
    footer {
      position: relative;
      z-index: 2;
      border-top: 1px solid var(--border);
      padding: 1.5rem 2rem;
      text-align: center;
    }

    .footer-inner {
      font-family: 'Cinzel', serif;
      font-size: 0.6rem;
      letter-spacing: 0.15em;
      color: var(--text-faint);
    }

    .footer-inner a { color: var(--biolum-dim); text-decoration: none; transition: color 0.15s; }
    .footer-inner a:hover { color: var(--biolum); }

    @media (max-width: 560px) {
      .cmd-layout {
        grid-template-columns: 1fr;
        grid-template-rows: auto auto auto;
      }
      .cmd-triggers {
        grid-column: 1;
        grid-row: auto;
        flex-direction: row;
        flex-wrap: wrap;
        justify-content: flex-start;
      }
    }
  </style>
</head>
<body>

<div id="particles"></div>
<div class="copy-toast" id="copy-toast"></div>

<header>
  <div class="header-inner">
    <div class="depth-marker">depth level — commands</div>
    <h1>lazybot<em>33</em></h1>
    <p class="subtitle">a catalogue of commands from the deep ^w^</p>
    <div class="status-row">
      <span class="status-dot"></span>
      <span id="cmd-count">loading...</span>
    </div>
    <div class="header-rule"></div>
  </div>
</header>

<main>
  <div id="active-section">
    <div class="section-label">known commands</div>
    <div class="commands-list" id="active-list"></div>
  </div>

  <div class="inactive-section" id="inactive-section" style="display:none">
    <div class="section-label">dormant commands</div>
    <div class="commands-list" id="inactive-list"></div>
  </div>
</main>

<footer>
  <div class="footer-inner">
    <a href="/">lazyturtle33.live</a>
    &nbsp;✦&nbsp;
    <a href="https://twitch.tv/thelazyturtle33" target="_blank" rel="noopener">twitch.tv/thelazyturtle33</a>
  </div>
</footer>

<script>
  (function() {
    const container = document.getElementById('particles');
    for (let i = 0; i < 35; i++) {
      const p = document.createElement('div');
      p.className = 'particle';
      const size = Math.random() * 2.5 + 0.5;
      p.style.cssText = [
        'width:'  + size + 'px',
        'height:' + size + 'px',
        'left:'   + (Math.random() * 100) + '%',
        'top:'    + (Math.random() * 100) + '%',
        'animation-duration:'  + (Math.random() * 20 + 15) + 's',
        'animation-delay:'     + (Math.random() * -25) + 's',
      ].join(';');
      container.appendChild(p);
    }
  })();

  const RAW = ` + "`" + commandsJSON + "`" + `;

  function buildCard(cmd, delay) {
    const inactive = cmd.is_active === false;
    const card = document.createElement('div');
    card.className = 'cmd-card' + (inactive ? ' inactive' : '');
    card.style.animationDelay = delay + 'ms';

    const triggers = (cmd.triggers || []).map(t =>
      ` + "`<span class=\"cmd-trigger\">${t}</span>`" + `
    ).join('');

    card.innerHTML = ` + "`" + `
      <div class="cmd-layout">
        <div class="cmd-header">
          <span class="cmd-name">${cmd.name}</span>
          <span class="cmd-usage">${cmd.usage || ''}</span>
          ${inactive ? '<span class="inactive-badge">dormant</span>' : ''}
        </div>
        <div class="cmd-triggers">${triggers}</div>
        <div class="cmd-desc">${cmd.description || 'no description provided.'}</div>
      </div>
    ` + "`" + `;

    return card;
  }

  let toastTimer = null;
  function showToast(text) {
    const toast = document.getElementById('copy-toast');
    toast.textContent = '"' + text + '" copied to clipboard';
    toast.classList.add('show');
    clearTimeout(toastTimer);
    toastTimer = setTimeout(() => toast.classList.remove('show'), 2000);
  }

  function render() {
    let commands = [];
    try { commands = JSON.parse(RAW); } catch(e) {}

    const active   = commands.filter(c => c.is_active !== false);
    const inactive = commands.filter(c => c.is_active === false);

    const activeList   = document.getElementById('active-list');
    const inactiveList = document.getElementById('inactive-list');

    if (active.length === 0 && inactive.length === 0) {
      activeList.innerHTML = '<div class="empty">nothing here yet... the deep is quiet.</div>';
    }

    active.forEach((cmd, i) => {
      const card = buildCard(cmd, i * 60);
      const firstTrigger = (cmd.triggers || [])[0];
      if (firstTrigger) {
        card.addEventListener('click', () => {
          navigator.clipboard.writeText(firstTrigger).then(() => showToast(firstTrigger));
        });
      }
      activeList.appendChild(card);
    });

    if (inactive.length > 0) {
      document.getElementById('inactive-section').style.display = '';
      inactive.forEach((cmd, i) => inactiveList.appendChild(buildCard(cmd, i * 60)));
    }

    const total = commands.length;
    document.getElementById('cmd-count').textContent =
      total + (total === 1 ? ' command' : ' commands') + ' known to the bot';
  }

  render();
</script>
</body>
</html>`
}
