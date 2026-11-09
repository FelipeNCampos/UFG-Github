document.getElementById('openPopup').addEventListener('click', function() {
    var popupElement = document.getElementById('popup');
    popupElement.style.display = 'grid';
});

document.getElementById('voltar').addEventListener('click', function() {
    window.location.href = "/home"
});
