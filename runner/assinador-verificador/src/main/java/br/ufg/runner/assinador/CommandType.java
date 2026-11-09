package br.ufg.runner.assinador;

enum CommandType {
    SIGN("sign"),
    VALIDATE("validate"),
    SERVER("server");

    private final String cliName;

    CommandType(String cliName) {
        this.cliName = cliName;
    }

    String cliName() {
        return cliName;
    }

    static CommandType fromCliName(String rawValue) throws ValidationException {
        for (CommandType value : values()) {
            if (value.cliName.equalsIgnoreCase(rawValue)) {
                return value;
            }
        }

        throw new ValidationException("Comando invalido: " + rawValue + ". Use sign, validate ou server.");
    }
}

