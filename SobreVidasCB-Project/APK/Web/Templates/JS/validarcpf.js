$(document).ready(function() {
    $("#cpf").inputmask("999.999.999-99");
});
    // Função para validar CPF
    function validarCPF(cpf) {
        cpf = cpf.replace(/[^\d]+/g, ''); // Remove caracteres não numéricos do CPF
        
        if (cpf.length !== 11 ||
            cpf === '00000000000' ||
            cpf === '11111111111' ||
            cpf === '22222222222' ||
            cpf === '33333333333' ||
            cpf === '44444444444' ||
            cpf === '55555555555' ||
            cpf === '66666666666' ||
            cpf === '77777777777' ||
            cpf === '88888888888' ||
            cpf === '99999999999') {
            return false;
        }

        // Verificação dos dígitos verificadores
        let soma = 0;
        let resto;

        for (let i = 1; i <= 9; i++) {
            soma = soma + parseInt(cpf.substring(i - 1, i)) * (11 - i);
        }
        resto = (soma * 10) % 11;

        if ((resto === 10) || (resto === 11)) {
            resto = 0;
        }
        if (resto !== parseInt(cpf.substring(9, 10))) {
            return false;
        }

        soma = 0;
        for (let i = 1; i <= 10; i++) {
            soma = soma + parseInt(cpf.substring(i - 1, i)) * (12 - i);
        }
        resto = (soma * 10) % 11;

        if ((resto === 10) || (resto === 11)) {
            resto = 0;
        }
        if (resto !== parseInt(cpf.substring(10, 11))) {
            return false;
        }
        return true; // CPF válido
    }

    // Evento de input no campo CPF para validar e habilitar o botão de Entrar
    $("#cpf").on('input', function() {
        const cpf = $(this).val();
        const cpfValido = validarCPF(cpf);

        if (cpfValido) {
            $("#entrar").prop('disabled', false);
        } else {
            $("#entrar").prop('disabled', true);
        }
    });

    // Evento de clique no botão Entrar
    $("#entrar").on('click', function() {
        const cpf = $("#cpf").val();
        const senha = $("#senha").val(); // Supondo que você também tenha um campo de senha

        if (!validarCPF(cpf)) {
            alert("CPF inválido.");
            return;
        }

        // Aqui você pode realizar outras validações antes de redirecionar
        // Por exemplo, validar a senha, consultar no banco de dados, etc.
        
        // Se todas as validações passarem, redirecione para submitlogin.html
        window.location.href = "submitlogin.html";
    });


 