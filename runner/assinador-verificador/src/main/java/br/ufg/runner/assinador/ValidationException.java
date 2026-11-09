package br.ufg.runner.assinador;

final class ValidationException extends Exception {
    private final String code;
    private final String field;

    ValidationException(String message) {
        this("VALIDATION_ERROR", null, message);
    }

    ValidationException(String code, String message) {
        this(code, null, message);
    }

    ValidationException(String code, String field, String message) {
        super(message);
        this.code = code;
        this.field = field;
    }

    String code() {
        return code;
    }

    String field() {
        return field;
    }
}

