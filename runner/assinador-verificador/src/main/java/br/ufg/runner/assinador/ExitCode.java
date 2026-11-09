package br.ufg.runner.assinador;

enum ExitCode {
    SUCCESS(0),
    RUNTIME_ERROR(1),
    VALIDATION_ERROR(2),
    SERVER_RUNNING(3);

    private final int value;

    ExitCode(int value) {
        this.value = value;
    }

    int value() {
        return value;
    }
}

