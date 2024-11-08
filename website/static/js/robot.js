document.addEventListener('DOMContentLoaded', () => {
    // Select all robots and cells
    const robots = document.querySelectorAll('.robot');
    const cells = document.querySelectorAll('.cell');

    // Make robots draggable
    robots.forEach(robot => {
        robot.setAttribute('draggable', true);

        // Start dragging a robot
        robot.addEventListener('dragstart', (event) => {
            event.dataTransfer.setData('text/plain', event.target.id);
        });
    });

    // Allow cells to be drop targets
    cells.forEach(cell => {
        cell.addEventListener('dragover', (event) => {
            event.preventDefault();
        });

        cell.addEventListener('drop', (event) => {
            event.preventDefault();

            // Get the robot ID from the dragged element
            const robotId = event.dataTransfer.getData('text/plain');
            const robot = document.getElementById(robotId);

            // Check if the cell is within the no-drop zone (row 7 or 8 and column 7 or 8)
            const row = parseInt(cell.getAttribute('data-row'));
            const column = parseInt(cell.getAttribute('data-column'));
            if ((row === 7 || row === 8) && (column === 7 || column === 8)) {
                console.log("Dropping in the center is not allowed.");
                return; // Prevent dropping in the no-drop zone
            }

            // Check if the cell already contains a robot
            if (cell.querySelector('.robot')) {
                console.log("This cell already has a robot.");
                return; // Prevent stacking robots
            }

            // Append the robot to the new cell and get its new position
            if (robot && cell) {
                cell.appendChild(robot);

                // Get robot color from its class
                const color = robot.classList.contains('yellow') ? 'Yellow' :
                    robot.classList.contains('red') ? 'Red' :
                        robot.classList.contains('green') ? 'Green' :
                            robot.classList.contains('blue') ? 'Blue' : 'Unknown';

                // Get the new position
                const position = getCellPosition(cell);

                // Call the changeRobot function with the new data
                changeRobot({ color, position });
            }
        });
    });
})

// Function to call when a robot is moved
function changeRobot(data) {
    console.log(data); // Log the data or handle it as needed
    // Send a POST request to the server with the new robot data
    sendPostRequest('/api/change-robot', data)
        .then(data => { console.log(data); })
        .then(() => location.reload())
        .catch(error => console.error('Error:', error));
}

// Helper function to get row and column directly from cell attributes
function getCellPosition(cell) {
    const row = parseInt(cell.getAttribute('data-row'));
    const column = parseInt(cell.getAttribute('data-column'));
    return { row, column };
}
