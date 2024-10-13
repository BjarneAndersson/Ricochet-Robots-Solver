/**
 * Sends a POST request to the specified URL.
 *
 * @param {string} url - The URL to send the POST request to.
 * @returns {Promise<Object>} - A promise that resolves with the response data.
 */
function sendPostRequest(url) {
    return fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({})
    })
        .then(response => {
            if (response.status === 204 || response.headers.get('content-length') === '0') {
                return {};
            }
            return response.json();
        });
}

/**
 * Solves the board and updates the tracking data.
 */
function solve() {
    sendPostRequest('/api/solve')
        .then(data => {
            console.log(data);
            const trackingData = data.TrackingData;
            document.getElementById('init-board-states').textContent = trackingData.InitializedBoardStates;
            document.getElementById('evaluated-board-states').textContent = trackingData.EvaluatedBoardStates;
            document.getElementById('duration').textContent = trackingData.Duration;
        })
        .catch(error => console.error('Error:', error));
}
