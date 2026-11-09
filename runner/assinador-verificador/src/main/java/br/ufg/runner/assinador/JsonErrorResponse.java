package br.ufg.runner.assinador;

final class JsonErrorResponse {
    private JsonErrorResponse() {
    }

    static String validation(ValidationException error) {
        return build("VALIDATION_ERROR", error.code(), error.getMessage(), error.field());
    }

    static String runtime(String message) {
        return build("RUNTIME_ERROR", "RUNTIME_ERROR", message, null);
    }

    private static String build(String type, String code, String message, String field) {
        String safeMessage = message != null ? message : "Erro desconhecido.";
        StringBuilder builder = new StringBuilder();
        builder.append("{\n")
            .append("  \"status\": \"ERROR\",\n")
            .append("  \"type\": \"").append(JsonEscaper.escape(type)).append("\",\n")
            .append("  \"code\": \"").append(JsonEscaper.escape(code)).append("\",\n")
            .append("  \"message\": \"").append(JsonEscaper.escape(safeMessage)).append("\"");

        if (field != null && !field.isBlank()) {
            builder.append(",\n")
                .append("  \"field\": \"").append(JsonEscaper.escape(field)).append("\"");
        }

        builder.append("\n}");
        return builder.toString();
    }
}