document.addEventListener('DOMContentLoaded', function () {
    const cpf = document.getElementById('cpf');
    const senha = document.getElementById('senha');
    const entrar = document.getElementById('entrar');

    function validateInput() {
        entrar.disabled = !(cpf.value.trim() && senha.value.trim());
    }

    cpf.addEventListener('input', validateInput);
    senha.addEventListener('input', validateInput);
});
