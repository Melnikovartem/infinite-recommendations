const userId = `user_${Math.random().toString(36).substr(2, 9)}`;
let loading = false;
let offset = 0;

// Reset offset every minute
setInterval(() => {
    offset = 0;
}, 60000);

const calculateInitialLoad = () => {
    const squareSize = 120; // 100px square + 20px margin
    const squaresPerRow = Math.floor(window.innerWidth / squareSize);
    const rows = Math.ceil(window.innerHeight / squareSize);
    return squaresPerRow * rows + 20; // Add buffer
};

const loadRecommendations = async (number) => {
    if (loading) return;
    loading = true;

    const response = await fetch(`http://178.62.197.190:8080/recommend?userId=${userId}&n=${number}&offset=${offset}`);
    const recommendations = await response.json();

    const content = document.getElementById('content');
    recommendations.forEach(recommendation => {
        const square = document.createElement('div');
        square.className = 'square';
        square.style.backgroundColor = recommendation.html_color;
        square.addEventListener('click', () => likeColor(recommendation.html_color));
        content.appendChild(square);
    });

    offset += number; // Update offset
    loading = false;
};

const likeColor = async (color) => {
    await fetch('http://178.62.197.190:8080/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userId: userId,
            html_color: color
        })
    });
};

window.addEventListener('scroll', () => {
    if (window.innerHeight + window.scrollY >= document.documentElement.scrollHeight - 100) {
        loadRecommendations(20); // Load more with buffer
    }
});

// Initial load
loadRecommendations(calculateInitialLoad());
