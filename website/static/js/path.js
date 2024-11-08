function updatePath(solveResponse) {
    const slider = document.getElementById('slider');
    const sliderValues = document.getElementById('slider-values');

    slider.hidden = false; // Show the slider

    // Update the slider range based on the number of path steps
    const pathSteps = solveResponse.length - 1; // Drop the last path element
    slider.max = pathSteps - 1;

    // Function to get the arrow icon based on the direction
    function getArrowIcon(direction) {
        switch (direction) {
            case 'up': return 'fas fa-arrow-up';
            case 'down': return 'fas fa-arrow-down';
            case 'left': return 'fas fa-arrow-left';
            case 'right': return 'fas fa-arrow-right';
            default: return '';
        }
    }

    // Function to create span elements based on slider range
    function createSliderValues() {
        const min = parseInt(slider.min);
        const max = parseInt(slider.max);
        sliderValues.innerHTML = ''; // Clear any existing spans

        for (let i = min; i <= max; i++) {
            const step = solveResponse[i];

            if (i > 0) { // Leave out the arrow for slider-value=0
                const arrow = document.createElement('i');
                arrow.className = `${getArrowIcon(step.moved_direction)} fa-lg`;
                arrow.style.color = step.moved_robot;
                sliderValues.appendChild(arrow);
            }

            const span = document.createElement('span');
            span.textContent = i;
            span.classList.add('slider-value');
            if (i === parseInt(slider.value)) {
                span.classList.add('active');
            }

            sliderValues.appendChild(span);
        }
    }

    // Function to move the robot to a specified cell
    function changeRobotPosition(robotColor, position) {
        const robot = document.querySelector(`.robot.${robotColor.toLowerCase()}`);
        const targetCell = document.querySelector(`.cell[data-row="${position.row}"][data-column="${position.column}"]`);

        if (robot && targetCell) {
            targetCell.appendChild(robot); // Move the robot to the target cell
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

    // Create the slider values initially
    createSliderValues();

    // Add event listener to update the active span and robot positions when the slider changes
    slider.addEventListener('input', updateActiveSpan);

    // Initial update of the path display and robot positions
    updateActiveSpan();
}
