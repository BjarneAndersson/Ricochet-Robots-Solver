function sendGetRequest(url) {
    return fetch(url, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (response.status === 204 || response.headers.get('content-length') === '0') {
            return {};
        }
        return response.json();
    })
    .catch(error => {
        console.error('Error in GET request:', error);
        throw error; // Ensure the error is propagated
    });
}

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
        if (response.status !== 200 || response.headers.get('content-length') === '0') {
            return { };
        }
        return response.json();
    })
    .catch(error => {
        console.error('Error in POST request:', error);
        throw error; // Ensure the error is propagated
    });
}

/**
 * Solves the board and updates the tracking data.
 * @returns {Promise<void>} - A promise that resolves when the
 */
function solve() {
    // Select all robots and cells
    const robots = document.querySelectorAll('.robot');

    // Make robots draggable
    robots.forEach(robot => {
        robot.setAttribute('draggable', false);
    });

    sendPostRequest('/api/solve')
        .then(data => {
            console.log(data);
            // update tracking data
            updateTrackingData(data.TrackingData);

            // update the path
            updatePath(data.Path);
        })
        .catch(error => console.error('Error:', error));

    // Make robots draggable
    robots.forEach(robot => {
        robot.setAttribute('draggable', true);
    });
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

function changeTarget(body) {
    if (Object.keys(body).length === 0) {
        body = {
            "color": 'multicolor',
            "symbol": 'vertex'
        }
    }
    sendPostRequest('/api/change-target', body)
        .then(data => { console.log(data); })
        .then(() => location.reload())
        .catch(error => console.error('Error:', error));
    // The location.reload() method is used to reload the current page.
}

function updateTrackingData(trackingData) {
    document.getElementById('init-board-states').textContent = trackingData.InitializedBoardStates;
    document.getElementById('evaluated-board-states').textContent = trackingData.EvaluatedBoardStates;
    document.getElementById('duration').textContent = trackingData.Duration;
}
