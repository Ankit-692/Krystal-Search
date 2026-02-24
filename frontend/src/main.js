import './style.css';
import './app.css';

const input = document.getElementById('search');
const resultsArea = document.getElementById('results');

input.addEventListener('input', async (e) => {
  const query = e.target.value.trim();
  if (query.length < 2) {
    resultsArea.innerHTML = '';
    return;
  }
  try {
    const results = await window.go.main.App.Search(query);
    resultsArea.innerHTML = results.map((item, index) => `
    <div class="result-item" data-path="${item.Path}" data-title="${item.Title}">
        <div class="result-title">${item.Title}</div>
        <div class="result-desc">${item.Path}</div>
    </div>
`).join('');
  } catch (err) {
    console.error(err);
  }

});

resultsArea.addEventListener('click', (e) => {
  LaunchApp(e);  
});

input.addEventListener('keydown',(e)=>{
  if(e.key === 'Enter'){
    LaunchApp(e);
  }
})

window.runtime.EventsOn("focus_search", () => {
  input.focus();
  input.select();
});


window.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') {
    window.runtime.WindowHide();
  }
});

window.onblur = () => {
  window.runtime.WindowHide();
};

function LaunchApp(e){
const item = e.target.closest('.result-item');
  if (item) {
    const path = item.getAttribute('data-path');
    const title = item.getAttribute('data-title');
    window.go.main.App.Launch({ Title: title, Path: path });
  }
}