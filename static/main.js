let isLogin = false;
        function toggleMode() {
            isLogin = !isLogin;
            document.getElementById('title').innerText = isLogin ? 'Вход' : 'Регистрация';
            document.getElementById('message').innerText = '';
        }

        document.getElementById('authForm').onsubmit = async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const messageEl = document.getElementById('message');
            
            const url = isLogin ? '/auth/login' : '/auth/register';

            try {
                const response = await fetch(url, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });

                const data = await response.json();

                if (response.ok) {
                    messageEl.style.color = 'green';
                    messageEl.innerText = isLogin ? `Успех! Токен: ${data.token.substring(0, 20)}...` : 'Регистрация успешна! Теперь войдите.';
                } else {
                    messageEl.style.color = 'red';
                    messageEl.innerText = `Ошибка: ${data.error || 'Неизвестная ошибка'}`;
                }
            } catch (err) {
                messageEl.style.color = 'red';
                messageEl.innerText = 'Ошибка сети или сервер недоступен';
            }
        };