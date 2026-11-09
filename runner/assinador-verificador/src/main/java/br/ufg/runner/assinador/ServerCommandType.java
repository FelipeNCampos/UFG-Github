package br.ufg.runner.assinador;

enum ServerCommandType {
    START("start"),
    STATUS("status"),
    STOP("stop");

    private final String cliName;

    ServerCommandType(String cliName) {
        this.cliName = cliName;
    }

    String cliName() {
        return cliName;
    }

    static ServerCommandType fromCliName(String rawValue) throws ValidationException {
        for (ServerCommandType value : values()) {
            if (value.cliName.equalsIgnoreCase(rawValue)) {
                return value;
            }
        }

        throw new ValidationException("Acao de server invalida: " + rawValue + ". Use start, status ou stop.\n" + CliArguments.usage());
    }
}
