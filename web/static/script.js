const focusInput = () => {
  const input = document.getElementById('input');
  input.focus();
};

const renderItem = item => {
  const li = document.createElement('li');
  li.className = 'item' + (item.checked ? ' checked' : '');
  li.dataset.id = item.id;

  const text = document.createElement('button');
  text.className = 'text';
  text.textContent = item.name;
  li.appendChild(text);

  if (item.checked) {
    const cross = document.createElement('button');
    cross.className = 'cross';
    cross.textContent = '♻️';
    cross.addEventListener('click', handleCrossClick);

    li.appendChild(cross);
  }

  return li;
};

const renderList = async () => {
  const list = document.getElementById('list');
  const fragment = document.createDocumentFragment();

  const response = await fetch('/api/items');
  const { items } = await response.json();

  list.innerHTML = '';
  fragment.append(...items.map(renderItem));
  list.appendChild(fragment);
};

const handleCrossClick = async event => {
  event.stopPropagation();

  const item = event.target.closest('[data-id]');
  item.style.pointerEvents = 'none';

  await fetch(`/api/items/${item.dataset.id}`, { method: 'DELETE' });
  await renderList();

  item.style.pointerEvents = 'auto';
  focusInput();
};

const handleListClick = async event => {
  if (!event.target.classList.contains('text')) {
    return;
  }

  const item = event.target.closest('[data-id]');
  const isChecked = item.classList.contains('checked');
  item.style.pointerEvents = 'none';

  await fetch(`/api/items/${item.dataset.id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ checked: !isChecked }),
  });
  await renderList();

  item.style.pointerEvents = 'auto';
  focusInput();
};

const handleSubmit = async event => {
  event.preventDefault();

  const input = document.getElementById('input');
  input.disabled = true;

  await fetch('/api/items', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name: input.value.trim() }),
  });
  await renderList();

  input.value = '';
  input.disabled = false;
  focusInput();
};

const init = () => {
  const form = document.getElementById('form');
  const list = document.getElementById('list');

  form.addEventListener('submit', handleSubmit);
  list.addEventListener('click', handleListClick);

  focusInput();
  renderList();
};

document.addEventListener('DOMContentLoaded', init);
