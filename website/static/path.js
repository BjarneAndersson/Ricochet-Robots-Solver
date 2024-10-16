function updatePath(solveResponse) {
    const slider = document.getElementById('slider');
    const sliderValues = document.getElementById('slider-values');

    slider.hidden = false; // Show the slider

    // Update the slider range based on the number of path steps
    const pathSteps = solveResponse.length - 1; // Drop the last path element
    slider.max = pathSteps - 1;

    // Function to get the arrow character based on the direction
    function getArrow(direction) {
        switch (direction) {
            case 'up': return '↑';
            case 'down': return '↓';
            case 'left': return '←';
            case 'right': return '→';
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
                const arrow = document.createElement('span');
                arrow.textContent = getArrow(step.moved_direction);
                if (step.moved_robot) {
                    arrow.classList.add('arrow', step.moved_robot);
                } else {
                    arrow.classList.add('arrow');
                }
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

    // Function to update the active span
    function updateActiveSpan() {
        const spans = sliderValues.querySelectorAll('.slider-value');
        spans.forEach(span => {
            span.classList.remove('active');
            if (parseInt(span.textContent) === parseInt(slider.value)) {
                span.classList.add('active');
            }
        });

        // Update the path display based on the current slider value
        const currentStep = solveResponse[parseInt(slider.value)];
        console.log(currentStep);
    }

    // Create the slider values initially
    createSliderValues();

    // Add event listener to update the active span when the slider changes
    slider.addEventListener('input', updateActiveSpan);

    // Initial update of the path display
    updateActiveSpan();
}