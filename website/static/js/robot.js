document.addEventListener('DOMContentLoaded', () => {
    // Select all robots and cells
    const robots = document.querySelectorAll('.robot');
    const cells = document.querySelectorAll('.cell');
    const slider = document.getElementById('slider');
    const sliderValues = document.getElementById('slider-values');

    // Make robots draggable
    robots.forEach(robot => {
        robot.setAttribute('draggable', true);

        // Start dragging a robot and reset slider to 0
        robot.addEventListener('dragstart', (event) => {
            event.dataTransfer.setData('text/plain', event.target.id);
            slider.value = 0; // Reset slider to 0 on drag
            updateActiveSpan(); // Update robots to initial position
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
                updateTransform(robot, cell); // Apply smooth animation to the new cell
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

    // Function to move the robot to a specified cell
    function changeRobotPosition(robotColor, position) {
        const robot = document.querySelector(`.robot.${robotColor.toLowerCase()}`);
        const targetCell = document.querySelector(`.cell[data-row="${position.row}"][data-column="${position.column}"]`);

        if (robot && targetCell) {
            updateTransform(robot, targetCell); // Smoothly animate robot to target cell
            targetCell.appendChild(robot); // Append robot to target cell for final positioning
        }
    }

    // Function to update the active span and move robots based on slider value
    function updateActiveSpan() {
        const spans = sliderValues.querySelectorAll('.slider-value');
        spans.forEach(span => {
            span.classList.remove('active');
            if (parseInt(span.textContent) === parseInt(slider.value)) {
                span.classList.add('active');
            }
        });

        // Update robot positions based on the current slider value
        const currentStep = solveResponse[parseInt(slider.value)];
        currentStep.robots.forEach(robotData => {
            changeRobotPosition(robotData.color, robotData.position);
        });
    }

    // Function to calculate transform offset and apply animation
    function updateTransform(robot, targetCell) {
        const robotRect = robot.getBoundingClientRect();
        const targetRect = targetCell.getBoundingClientRect();
        const offsetX = targetRect.left - robotRect.left;
        const offsetY = targetRect.top - robotRect.top;

        // Apply CSS transform for smooth transition
        robot.style.transform = `translate(${offsetX}px, ${offsetY}px)`;

        // Wait for the transition to finish before resetting the transform
        robot.addEventListener('transitionend', () => {
            robot.style.transform = ''; // Reset transform for accurate placement
        }, { once: true });
    }

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

    // Add event listener to update the active span and robot positions when the slider changes
    slider.addEventListener('input', () => {
        updateActiveSpan(); // Enable animation on slider change
    });

    // Initial update of the path display and robot positions
    updateActiveSpan();
});
