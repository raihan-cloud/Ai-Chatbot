function toggleSidebar() {
    document.getElementById('sidebar').classList.toggle('active');
}

const chatDisplay = document.getElementById('chat-display');
const userInput = document.getElementById('user-input');
const sendBtn = document.getElementById('send-btn');

function appendMessage(role, text) {
    // Sembunyikan Hero Section saat mulai chat
    const hero = document.querySelector('.hero-section');
    if (hero) hero.style.display = 'none';

    const msgDiv = document.createElement('div');
    msgDiv.className = `message ${role === 'user' ? 'user-msg' : 'bot-msg'}`;
    msgDiv.innerHTML = text;

    chatDisplay.appendChild(msgDiv);
    chatDisplay.scrollTo({ top: chatDisplay.scrollHeight, behavior: 'smooth' });
}

async function sendMessage() {
    const val = userInput.value.trim();
    if (!val) return;

    appendMessage('user', val);
    userInput.value = '';

    // Efek Typing Sederhana
    const loadingDiv = document.createElement('div');
    loadingDiv.className = 'message bot-msg';
    loadingDiv.innerText = 'Thinking...';
    chatDisplay.appendChild(loadingDiv);

    try {
        const res = await fetch('http://localhost:8080/chat', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ message: val })
        });
        const data = await res.json();
        loadingDiv.remove();
        appendMessage('bot', data.reply);
    } catch (e) {
        loadingDiv.innerText = 'Error connecting to server.';
    }
}

sendBtn.addEventListener('click', sendMessage);
userInput.addEventListener('keypress', (e) => { if(e.key === 'Enter') sendMessage(); });