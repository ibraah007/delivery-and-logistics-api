// Quick layout fix for keyboard interaction
document.addEventListener('keydown', (e) => {
    if (document.getElementById('authOverlay').style.display !== "none") {
        if (e.key >= '0' && e.key <= '9') pressKey(e.key);
        if (e.key === 'Enter') verifyPin();
        if (e.key === 'Backspace' || e.key === 'Escape') clearPin();
    }
});
