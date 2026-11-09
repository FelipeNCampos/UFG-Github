package br.ufg.runner.assinador;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

final class JsonPayloadValidator {
    private static final Pattern ALIAS_PATTERN = Pattern.compile("[A-Za-z0-9_-]{3,64}");
    private static final Pattern DIGITS_PATTERN = Pattern.compile("\\d+");
    private static final Pattern PKCS11_LIBRARY_PATTERN = Pattern.compile("(?i).*(\\.dll|\\.so|\\.dylib)$");
    private static final Pattern RESOURCE_TYPE_PATTERN = Pattern.compile("\\\"resourceType\\\"\\s*:\\s*\\\"([^\\\"]+)\\\"");
    private static final Pattern SIGNATURE_PATTERN = Pattern.compile("\\\"signature\\\"\\s*:\\s*\\\"([^\\\"]+)\\\"");
    private static final Pattern DIGEST_PATTERN = Pattern.compile("\\\"inputDigestSha256\\\"\\s*:\\s*\\\"([A-F0-9]{64})\\\"");

    private JsonPayloadValidator() {
    }

    static void validateSignPayload(String payload) throws ValidationException {
        validateObject(payload);
        requireField(payload, RESOURCE_TYPE_PATTERN, "resourceType", "INVALID_FHIR_PAYLOAD", "O JSON de entrada precisa conter o campo resourceType.");
    }

    static void validateValidatePayload(String payload) throws ValidationException {
        validateObject(payload);
        requireField(payload, SIGNATURE_PATTERN, "signature", "MISSING_REQUIRED_FIELD", "O arquivo informado para validate nao contem o campo obrigatorio \"signature\".");
        requireField(payload, DIGEST_PATTERN, "inputDigestSha256", "MISSING_REQUIRED_FIELD", "O arquivo informado para validate nao contem um inputDigestSha256 valido em hexadecimal maiusculo.");
    }

    static void validateAlias(String alias) throws ValidationException {
        if (alias == null || alias.isBlank()) {
            throw new ValidationException("MISSING_REQUIRED_FIELD", "alias", "O parametro alias e obrigatorio para o comando sign.");
        }

        if (!ALIAS_PATTERN.matcher(alias).matches()) {
            throw new ValidationException("INVALID_ALIAS", "alias", "O alias deve ter entre 3 e 64 caracteres e usar apenas letras, numeros, '-' ou '_'.");
        }
    }

    static void validatePkcs11Library(String pkcs11Library) throws ValidationException {
        if (pkcs11Library == null || pkcs11Library.isBlank()) {
            return;
        }

        if (!PKCS11_LIBRARY_PATTERN.matcher(pkcs11Library).matches()) {
            throw new ValidationException("INVALID_PKCS11_LIBRARY", "pkcs11Library", "A biblioteca PKCS#11 deve terminar com .dll, .so ou .dylib.");
        }
    }

    static void validatePkcs11Slot(String pkcs11Slot) throws ValidationException {
        if (pkcs11Slot == null || pkcs11Slot.isBlank()) {
            return;
        }

        if (!DIGITS_PATTERN.matcher(pkcs11Slot).matches()) {
            throw new ValidationException("INVALID_PKCS11_SLOT", "pkcs11Slot", "O slot PKCS#11 deve ser um numero inteiro nao negativo.");
        }
    }

    static String requireStringField(String payload, String fieldName, String errorCode, String message) throws ValidationException {
        Pattern pattern = Pattern.compile("\\\"" + Pattern.quote(fieldName) + "\\\"\\s*:\\s*\\\"([^\\\"]*)\\\"");
        return requireField(payload, pattern, fieldName, errorCode, message);
    }

    static boolean hasBalancedJsonStructure(String payload) {
        boolean inString = false;
        boolean escaped = false;
        int curlyDepth = 0;
        int squareDepth = 0;

        for (int index = 0; index < payload.length(); index++) {
            char character = payload.charAt(index);

            if (inString) {
                if (escaped) {
                    escaped = false;
                    continue;
                }
                if (character == '\\') {
                    escaped = true;
                    continue;
                }
                if (character == '"') {
                    inString = false;
                }
                continue;
            }

            switch (character) {
                case '"' -> inString = true;
                case '{' -> curlyDepth++;
                case '}' -> {
                    curlyDepth--;
                    if (curlyDepth < 0) {
                        return false;
                    }
                }
                case '[' -> squareDepth++;
                case ']' -> {
                    squareDepth--;
                    if (squareDepth < 0) {
                        return false;
                    }
                }
                default -> {
                    // no-op
                }
            }
        }

        return !inString && !escaped && curlyDepth == 0 && squareDepth == 0;
    }

    static void validateObject(String payload) throws ValidationException {
        String normalized = payload == null ? "" : payload.trim();
        if (normalized.isEmpty()) {
            throw new ValidationException("INVALID_JSON", "O conteudo do arquivo de entrada esta vazio.");
        }

        if (!(normalized.startsWith("{") && normalized.endsWith("}"))) {
            throw new ValidationException("INVALID_JSON", "O arquivo de entrada deve conter um objeto JSON.");
        }

        if (!hasBalancedJsonStructure(normalized)) {
            throw new ValidationException("INVALID_JSON", "O arquivo de entrada possui estrutura JSON invalida.");
        }
    }

    private static String requireField(
        String payload,
        Pattern pattern,
        String fieldName,
        String errorCode,
        String message
    ) throws ValidationException {
        Matcher matcher = pattern.matcher(payload);
        if (!matcher.find()) {
            throw new ValidationException(errorCode, fieldName, message);
        }

        String value = matcher.group(1);
        if (value == null || value.isBlank()) {
            throw new ValidationException(errorCode, fieldName, message);
        }

        return value;
    }
}