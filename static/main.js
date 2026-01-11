const statusMap = {
    0: 'Ожидает',
    1: 'В обработке',
    2: 'Отправлено',
    3: 'Ошибка',
    4: 'Отменено'
};

document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('notifyForm');
    const refreshBtn = document.getElementById('refreshBtn');

    // Загрузка списка при старте
    fetchNotifications();

    // Создание уведомления
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const data = {
            channel: document.getElementById('channel').value,
            target: document.getElementById('target').value,
            payload: btoa(document.getElementById('payload').value), // Кодируем в base64 для []byte в Go
            scheduled_at: new Date(document.getElementById('scheduledAt').value).toISOString()
        };

        const resp = await fetch('/notify', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        if (resp.ok) {
            form.reset();
            fetchNotifications();
        } else {
            alert('Ошибка при создании');
        }
    });

    refreshBtn.addEventListener('click', fetchNotifications);
});

async function fetchNotifications() {
    const list = document.getElementById('notifyList');
    try {
        // ВАЖНО: Тебе нужно реализовать этот эндпоинт в Go!
        const resp = await fetch('/notify');
        const data = await resp.json();

        list.innerHTML = '';
        data.forEach(n => {
            const row = `
                <tr>
                    <td>${n.id.substring(0, 8)}...</td>
                    <td>${n.channel}</td>
                    <td>${n.target}</td>
                    <td><span class="status-badge status-${n.status}">${statusMap[n.status]}</span></td>
                    <td>${new Date(n.scheduled_at).toLocaleString()}</td>
                    <td>
                        <button onclick="deleteNotify('${n.id}')" class="btn-secondary">Удалить</button>
                    </td>
                </tr>
            `;
            list.insertAdjacentHTML('beforeend', row);
        });
    } catch (e) {
        console.error('Failed to fetch', e);
    }
}

async function deleteNotify(id) {
    if (confirm('Удалить уведомление?')) {
        await fetch(`/notify/${id}`, { method: 'DELETE' });
        fetchNotifications();
    }
}
