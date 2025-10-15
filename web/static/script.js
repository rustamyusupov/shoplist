const handleSubmit = async event => {
  event.preventDefault();

  const input = form.querySelector('#input');
  const name = input.value.trim();

  if (!name) {
    return;
  }

  try {
    const response = await fetch('/api/items', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name }),
    });

    if (response.ok) {
      const { id } = await response.json();

      const list = document.getElementById('list');
      const li = document.createElement('li');
      li.className = 'item';
      li.dataset.id = id;
      li.textContent = name;
      list.appendChild(li);

      input.value = '';
    }
  } catch (error) {
    console.error('Failed to add item:', error);
  }
};

const handleItemClick = async event => {
  const item = event.currentTarget;
  const id = item.dataset.id;
  const isChecked = item.classList.contains('checked');

  try {
    const response = await fetch(`/api/items/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ checked: !isChecked }),
    });

    if (response.ok) {
      item.classList.toggle('checked');
    }
  } catch (error) {
    console.error('Failed to update item:', error);
  }
};

const init = async () => {
  const form = document.getElementById('form');
  form.addEventListener('submit', handleSubmit);

  const items = document.querySelectorAll('.item');
  items.forEach(item => {
    item.addEventListener('click', handleItemClick);
  });
};

document.addEventListener('DOMContentLoaded', init);
