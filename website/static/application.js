/**
 * Sends a POST request to the specified URL.
 *
 * @param {string} url - The URL to send the POST request to.
 * @param body - The body of the request.
 * @returns {Promise<Object>} - A promise that resolves with the response data.
 */
function sendPostRequest(url, body = {}) {
    return fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
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
 * @returns {Promise<void>} - A promise that resolves when the
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

/**
 * Change a quarter board.
 * @param body - The body of the request.
 * @returns {Promise<void>} - A promise that resolves when the board is changed.
 */
function changeQuarterBoard(body = {}) {
    if (Object.keys(body).length === 0) {
        body = {
            "yellow": {
                "position": 'top-right',
                "version": '1'
            },
            "red": {
                "position": 'top-left',
                "version": '1'
            },
            "green": {
                "position": 'bottom-left',
                "version": '1'
            },
            "blue": {
                "position": 'bottom-right',
                "version": '1'
            }
        }
    }
    sendPostRequest('/api/change-quarter-boards', body)
        .then(data => { console.log(data); })
        .then(() => location.reload())
        .catch(error => console.error('Error:', error));
    // The location.reload() method is used to reload the current page.

    // Load the new board into the solver
    sendPostRequest('api/load')
        .catch(error => console.error('Error:', error));
}
