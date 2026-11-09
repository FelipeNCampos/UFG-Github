package br.ufg.runner.assinador;

enum OperationMode {
    DIRECT("direct"),
    HTTP("http");

    private final String cliName;

    OperationMode(String cliName) {
        this.cliName = cliName;
    }

    String cliName() {
        return cliName;
    }

    static OperationMode fromCliName(String rawValue) throws ValidationException {
        for (OperationMode value : values()) {
            if (value.cliName.equalsIgnoreCase(rawValue)) {
                return value;
            }
        }

        throw new ValidationException(
            "Modo nao suportado: " + rawValue + ". Use direct (execucao direta) ou http (servidor HTTP em execucao)."
        );
    }
}
