function populateTargets(targets) {
    const targetsContainer = document.getElementById('targets-container');
    targetsContainer.innerHTML = ''; // Clear any existing targets

    targets.forEach(target => {
        const targetElement = document.createElement('div');
        targetElement.className = `target ${target.symbol} ${target.color}`;
        targetElement.addEventListener('click', () => changeTarget(target));
        targetsContainer.appendChild(targetElement);
    });
}

// Example usage with a list of targets
const targets = [
    { symbol: 'circle', color: 'yellow' },
    { symbol: 'triangle', color: 'yellow' },
    { symbol: 'square', color: 'yellow' },
    { symbol: 'hexagon', color: 'yellow' },
    { symbol: 'circle', color: 'red' },
    { symbol: 'triangle', color: 'red' },
    { symbol: 'square', color: 'red' },
    { symbol: 'hexagon', color: 'red' },
    { symbol: 'circle', color: 'green' },
    { symbol: 'triangle', color: 'green' },
    { symbol: 'square', color: 'green' },
    { symbol: 'hexagon', color: 'green' },
    { symbol: 'circle', color: 'blue' },
    { symbol: 'triangle', color: 'blue' },
    { symbol: 'square', color: 'blue' },
    { symbol: 'hexagon', color: 'blue' },
    { symbol: 'vertex', color: 'multicolor' }
];

window.onload = () => {
    populateTargets(targets);
    getTarget()
}