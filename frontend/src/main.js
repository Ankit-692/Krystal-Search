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

input.addEventListener('input', async (e) => {
  const query = e.target.value.trim();
  if (query.length < 2) {
    resultsArea.innerHTML = '';
    selectedIndex = -1;
    return;
  }
  try {
    const results = await window.go.main.App.Search(query);
    resultsArea.innerHTML = results.map((item) => `
      <div class="result-item" data-path="${item.Path}" data-title="${item.Title}">
        <div class="result-icon"><img src="${item.Icon}"/></div>
        <div class="result-title">${item.Title}</div>
        <div class="result-desc">${item.Path}</div>
      </div>
    `).join('');
    selectedIndex = -1;
  } catch (err) {
    console.error(err);
  }
});

resultsArea.addEventListener('click', (e) => {
  const item = e.target.closest('.result-item');
  if (item) LaunchApp(item);
});

input.addEventListener('keydown', (e) => {
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
}// window.addEventListener('keydown', (e) => {
//   if (e.key === 'Escape') {
//     window.runtime.WindowHide();
//   }
// });
//
// window.onblur = () => {
//   window.runtime.WindowHide();
// };
//

