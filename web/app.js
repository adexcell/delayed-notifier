const API_BASE = 'http://localhost:8080/api/v1/notifies';;

// Elements
const createForm = document.getElementById('createForm');
const btnCheckStatus = document.getElementById('btnCheckStatus');
const btnCancel = document.getElementById('btnCancel');
const notifyIdInput = document.getElementById('notify_id');
const createResult = document.getElementById('createResult');
const manageResult = document.getElementById('manageResult');

// Helper to show results
function showResult(element, isSuccess, message) {
    element.textContent = message;
    element.className = `result show ${isSuccess ? 'success' : 'error'}`;
}

// Create Notification
createForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = createForm.querySelector('button');
    const originalText = btn.textContent;
    btn.textContent = 'Scheduling...';
    btn.disabled = true;

    try {
        const payload = {
            recipient_email: document.getElementById('recipient_email').value,
            subject: document.getElementById('subject').value,
            body: document.getElementById('body').value,
            // Convert to RFC3339 / ISO-8601 string
            scheduled_at: new Date(document.getElementById('scheduled_at').value).toISOString()
        };

        const res = await fetch(API_BASE, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        const data = await res.json();

        if (res.ok) {
            showResult(createResult, true, `Success! ID: ${data.id}`);
            // Auto-fill the manage form for convenience
            notifyIdInput.value = data.id;
            createForm.reset();
        } else {
            showResult(createResult, false, `Error: ${JSON.stringify(data)}`);
        }
    } catch (err) {
        showResult(createResult, false, `Network Error: ${err.message}`);
    } finally {
        btn.textContent = originalText;
        btn.disabled = false;
    }
});

// Check Status
btnCheckStatus.addEventListener('click', async () => {
    const id = notifyIdInput.value.trim();
    if (!id) return showResult(manageResult, false, 'Please enter an ID');

    btnCheckStatus.textContent = 'Checking...';
    btnCheckStatus.disabled = true;

    try {
        const res = await fetch(`${API_BASE}/${id}`);
        if (res.status === 204) {
            showResult(manageResult, true, 'Status: No Content (Cancelled)');
        } else {
            const data = await res.json();
            if (res.ok) {
                showResult(manageResult, true, `Status: ${data.status}`);
            } else {
                showResult(manageResult, false, `Error: ${JSON.stringify(data)}`);
            }
        }
    } catch (err) {
        showResult(manageResult, false, `Network Error: ${err.message}`);
    } finally {
        btnCheckStatus.textContent = 'Check Status';
        btnCheckStatus.disabled = false;
    }
});

// Cancel Notification
btnCancel.addEventListener('click', async () => {
    const id = notifyIdInput.value.trim();
    if (!id) return showResult(manageResult, false, 'Please enter an ID');

    if (!confirm('Are you sure you want to cancel this notification?')) return;

    btnCancel.textContent = 'Cancelling...';
    btnCancel.disabled = true;

    try {
        const res = await fetch(`${API_BASE}/${id}`, {
            method: 'DELETE'
        });

        if (res.ok || res.status === 204) {
            showResult(manageResult, true, 'Notification cancelled successfully.');
        } else {
            const data = await res.json().catch(() => ({}));
            showResult(manageResult, false, `Error: ${res.status} ${JSON.stringify(data)}`);
        }
    } catch (err) {
        showResult(manageResult, false, `Network Error: ${err.message}`);
    } finally {
        btnCancel.textContent = 'Cancel Notification';
        btnCancel.disabled = false;
    }
});
