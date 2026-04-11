import './style.css';
import './app.css';

const input = document.getElementById('search');
const resultsArea = document.getElementById('results');
let selectedIndex = -1;

function getItems() {
  return resultsArea.querySelectorAll('.result-item');
}

function setSelected(index) {
  const items = getItems();
  if (!items.length) return;

  selectedIndex = Math.max(0, Math.min(index, items.length - 1));

  items.forEach((item, i) => {
    item.classList.toggle('selected', i === selectedIndex);
  });

  items[selectedIndex].scrollIntoView({ block: 'nearest' });
}

let debounceTimer = null;

input.addEventListener('input', async (e) => {
  const query = e.target.value.trim();

  if (query.length <= 2) {
    resultsArea.innerHTML = '';
    selectedIndex = -1;
    return;
  }

  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(async () => {
    try {
      if (query.startsWith("f:")) {
        var trimmerQuery = query.slice(2).trim();
        const results = await window.go.main.App.FileSearch(trimmerQuery);
        resultsArea.innerHTML = results.map((item) => `
          <div class="result-item" data-path="${item.Path}" data-title="${item.Name}">
            <div class="result-icon"><img src="${item.Icon}"/></div>
            <div class="result-title">${item.Name}</div>
            <div class="result-desc">${item.Path}</div>
          </div>
        `).join('');
        selectedIndex = -1;
      } else if (query.startsWith("t:")) {
        resultsArea.innerHTML = `
          <div class="terminal-block">
            <div class="terminal-output">Press Enter to Run....</div>
          </div>`;
      } else {
        const results = await window.go.main.App.Search(query);
        resultsArea.innerHTML = results.map((item) => `
          <div class="result-item" data-path="${item.Path}" data-title="${item.Title}">
            <div class="result-icon"><img src="${item.Icon}"/></div>
            <div class="result-title">${item.Title}</div>
            <div class="result-desc">${item.Path}</div>
          </div>
        `).join('');
        selectedIndex = -1;
      }
    } catch (err) {
      console.error(err);
    }
  }, 200);
});


resultsArea.addEventListener('click', (e) => {
  const item = e.target.closest('.result-item');
  if (item) LaunchApp(item);
});

input.addEventListener('keydown', async (e) => {
  const items = getItems();

  if (e.key === 'ArrowDown') {
    e.preventDefault();
    if (!items.length) return;
    setSelected(selectedIndex < 0 ? 0 : selectedIndex + 1);

  } else if (e.key === 'ArrowUp') {
    e.preventDefault();
    if (!items.length) return;
    setSelected(selectedIndex <= 0 ? 0 : selectedIndex - 1);

  } else if (e.key === 'Enter') {
    if (!items.length) {
      var query = input.value.trim();
      if (query.startsWith("t:")) {
        var trimmedQuery = query.slice(2).trim();
        if (!trimmedQuery) return;

        if (trimmedQuery.startsWith("sudo")) {
          showPasswordPrompt(trimmedQuery);
        }
        else {
          await runCommand(trimmedQuery);
        }
      }
      return;
    }
    const target = selectedIndex >= 0
      ? items[selectedIndex]
      : resultsArea.firstElementChild;
    LaunchApp(target);
  }
});

window.runtime.EventsOn("focus_search", () => {
  input.focus();
  input.select();
});

function LaunchApp(item) {
  if (item) {
    const path = item.getAttribute('data-path');
    const title = item.getAttribute('data-title');
    window.go.main.App.Launch({ Title: title, Path: path });
  }
}

window.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    window.runtime.WindowHide();
    input.innerHTML = '';
    resultsArea.innerHTML = '';
  }
});

window.onblur = () => {
  window.runtime.WindowHide();
  input.innerHTML = '';
  resultsArea.innerHTML = '';
};


function escapeHtml(str) {
  if (!str) return '';
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}

async function runCommand(trimmedQuery, password = "") {
  try {
    const output = await window.go.main.App.RunCommand(trimmedQuery, password);
    resultsArea.innerHTML = `
      <div class="terminal-block">
        <div class="terminal-prompt">$ ${trimmedQuery}</div>
        <pre class="terminal-output">${escapeHtml(output) || '(no output)'}</pre>
      </div>
    `;
  } catch (err) {
    resultsArea.innerHTML = `<div class="terminal-output terminal-error">Error: ${err}</div>`;
  }
}

function showPasswordPrompt(command) {
  resultsArea.innerHTML = `
    <div class="terminal-block">
      <div class="terminal-prompt">$ ${command}</div>
      <div class="sudo-row">
        <span class="terminal-prompt">[sudo] password:</span>
        <input id="sudo-input" type="password" autocomplete="off" spellcheck="false" />
      </div>
    </div>
  `;

  const sudoInput = document.getElementById('sudo-input');
  sudoInput.focus();

  sudoInput.addEventListener('keydown', async (e) => {
    if (e.key === 'Enter') {
      const password = sudoInput.value;
      sudoInput.disabled = true;
      await runCommand(command, password);
    }
  });
}
