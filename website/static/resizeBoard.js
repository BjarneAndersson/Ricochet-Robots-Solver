document.addEventListener("DOMContentLoaded", function() {
    const board = document.querySelector(".board");
    const resizeBoard = () => {
        const container = board.parentElement;
        const size = Math.min(container.clientWidth, container.clientHeight);
        board.style.width = `${size}px`;
        board.style.height = `${size}px`;
    };
    window.addEventListener("resize", resizeBoard);
    resizeBoard();
});