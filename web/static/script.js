const handleListClick = async event => {
  const item = event.target.closest('[data-id]');

  if (!item) {
    return;
  }

  const id = item.dataset.id;
  const isChecked = item.classList.contains('checked');

  item.classList.toggle('checked');
  item.style.pointerEvents = 'none';

  try {
    const response = await fetch(`/api/items/${id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ checked: !isChecked }),
    });

    if (!response.ok) {
      item.classList.toggle('checked');
      console.error('Failed to update item');
    }
  } catch (error) {
    item.classList.toggle('checked');
    console.error('Failed to update item:', error);
  } finally {
    item.style.pointerEvents = 'auto';
  }
};

const handleSubmit = list => async event => {
  event.preventDefault();

  const input = document.getElementById('input');
  const name = input.value.trim();

  if (!name) {
    return;
  }

  if (input) {
    input.disabled = true;
  }

  try {
    const response = await fetch('/api/items', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name }),
    });

    if (response.ok) {
      const { id } = await response.json();

      const li = document.createElement('li');
      li.className = 'item';
      li.dataset.id = id;
      li.textContent = name;
      list.appendChild(li);

      input.value = '';
    } else {
      console.error('Failed to add item: Server error');
    }
  } catch (error) {
    console.error('Failed to add item:', error);
  } finally {
    if (input) input.disabled = false;
  }
};

const init = () => {
  const list = document.getElementById('list');
  const form = document.getElementById('form');

  list.addEventListener('click', handleListClick);
  form.addEventListener('submit', handleSubmit(list));
};

document.addEventListener('DOMContentLoaded', init);
