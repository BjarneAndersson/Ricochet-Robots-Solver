/**
 * Sends a POST request to the specified URL.
 *
 * @param {string} url - The URL to send the POST request to.
 * @returns {Promise<void>} - A promise that resolves when the request is complete.
 */
function sendPostRequest(url) {
    fetch(url, {
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
        })
        .then(data => console.log(data))
        .catch(error => console.error('Error:', error));
}