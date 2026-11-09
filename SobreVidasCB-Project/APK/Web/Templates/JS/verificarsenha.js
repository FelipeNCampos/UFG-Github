function verificarSenhas() {
    var senha1 = document.getElementById("senha").value;
    var senha2 = document.getElementById("senha2").value;

    if (senha1 !== senha2) {
        alert("As senhas não coincidem! Por favor, digite novamente.");
        return;
    }

    // Senhas coincidem, redirecionar para a página de login
    window.location.href = "login.html";
}