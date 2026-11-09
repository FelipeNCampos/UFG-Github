package br.ufg.runner.assinador;

import java.io.InputStream;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.URI;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.security.MessageDigest;
import java.util.HexFormat;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

final class SimulatedSigningService {
    private static final String SIGNATURE_PREFIX = "SIMULATED-SIGNATURE-";
    private static final String SIGNATURE_SCHEME = "SIMULATED-SIGNATURE-V1";
    private static final Pattern SIGNATURE_PATTERN = Pattern.compile("\"signature\"\\s*:\\s*\"([^\"]+)\"");
    private static final Pattern DIGEST_PATTERN = Pattern.compile("\"inputDigestSha256\"\\s*:\\s*\"([^\"]+)\"");
    private static final Pattern SIGNER_PATTERN = Pattern.compile("\"signedBy\"\\s*:\\s*\"([^\"]+)\"");

    String execute(CliArguments arguments) throws Exception {
        return switch (arguments.commandType()) {
            case SIGN -> sign(arguments);
            case VALIDATE -> validate(arguments);
            case SERVER -> executeServer(arguments);
        };
    }

    private String executeServer(CliArguments arguments) throws Exception {
        return switch (arguments.serverCommandType()) {
            case START -> SimulatedServerService.start(arguments.serverPort(), arguments.inactivityTimeoutMinutes());
            case STATUS -> SimulatedServerService.status(arguments.serverPort());
            case STOP -> SimulatedServerService.stop(arguments.serverPort());
        };
    }

    private String sign(CliArguments arguments) throws Exception {
        if (arguments.mode() == OperationMode.HTTP) {
            return signViaHttp(arguments);
        }

        String content = Files.readString(arguments.inputPath(), StandardCharsets.UTF_8);
        return signContent(content, arguments.alias(), arguments.pkcs11Library(), arguments.pkcs11Slot(), arguments.mode().cliName());
    }

    private String signViaHttp(CliArguments arguments) throws Exception {
        String content = Files.readString(arguments.inputPath(), StandardCharsets.UTF_8);

        StringBuilder urlBuilder = new StringBuilder("http://localhost:")
            .append(arguments.serverPort())
            .append("/sign?alias=")
            .append(URLEncoder.encode(arguments.alias(), StandardCharsets.UTF_8));

        if (arguments.pkcs11Library() != null) {
            urlBuilder.append("&pkcs11Library=")
                .append(URLEncoder.encode(arguments.pkcs11Library(), StandardCharsets.UTF_8));
        }

        if (arguments.pkcs11Slot() != null) {
            urlBuilder.append("&pkcs11Slot=")
                .append(URLEncoder.encode(arguments.pkcs11Slot(), StandardCharsets.UTF_8));
        }

        return postToServer(urlBuilder.toString(), content);
    }

    private String validate(CliArguments arguments) throws Exception {
        if (arguments.mode() == OperationMode.HTTP) {
            return validateViaHttp(arguments);
        }

        String content = Files.readString(arguments.inputPath(), StandardCharsets.UTF_8);
        return validateContent(content, arguments.mode().cliName());
    }

    private String validateViaHttp(CliArguments arguments) throws Exception {
        String content = Files.readString(arguments.inputPath(), StandardCharsets.UTF_8);
        String url = "http://localhost:" + arguments.serverPort() + "/validate";
        return postToServer(url, content);
    }

    private static String postToServer(String urlStr, String content) throws Exception {
        HttpURLConnection conn = (HttpURLConnection) URI.create(urlStr).toURL().openConnection();
        conn.setRequestMethod("POST");
        conn.setDoOutput(true);
        conn.setRequestProperty("Content-Type", "application/json; charset=UTF-8");
        conn.setConnectTimeout(5_000);
        conn.setReadTimeout(30_000);

        try (OutputStream out = conn.getOutputStream()) {
            out.write(content.getBytes(StandardCharsets.UTF_8));
        }

        int statusCode = conn.getResponseCode();
        InputStream responseStream = statusCode >= 400 ? conn.getErrorStream() : conn.getInputStream();
        String response = responseStream != null ? new String(responseStream.readAllBytes(), StandardCharsets.UTF_8) : "";

        if (statusCode >= 500) {
            throw new RuntimeException("Erro no servidor HTTP (status " + statusCode + "): " + response.strip());
        }

        if (statusCode >= 400) {
            throw new ValidationException("Erro de validacao no servidor HTTP (status " + statusCode + "): " + response.strip());
        }

        return response;
    }

    // Package-private: also called by SimulatedServerService HTTP handlers.
    static String signContent(
        String content,
        String alias,
        String pkcs11Library,
        String pkcs11Slot,
        String modeName
    ) throws Exception {
        JsonPayloadValidator.validateSignPayload(content);
        JsonPayloadValidator.validateAlias(alias);
        JsonPayloadValidator.validatePkcs11Library(pkcs11Library);
        JsonPayloadValidator.validatePkcs11Slot(pkcs11Slot);

        String digest = sha256(content);
        String signature = SIGNATURE_PREFIX + digest;
        String library = pkcs11Library == null ? "not-informed" : pkcs11Library;
        String slot = pkcs11Slot == null ? "not-informed" : pkcs11Slot;

        return """
            {
              "status": "SUCCESS",
              "operation": "sign",
              "mode": "%s",
              "message": "Assinatura simulada gerada com sucesso.",
              "digestAlgorithm": "SHA-256",
              "signatureScheme": "%s",
              "inputDigestSha256": "%s",
              "signature": "%s",
              "signedBy": "%s",
              "pkcs11Library": "%s",
              "pkcs11Slot": "%s"
            }
            """.formatted(
            JsonEscaper.escape(modeName),
            JsonEscaper.escape(SIGNATURE_SCHEME),
            JsonEscaper.escape(digest),
            JsonEscaper.escape(signature),
            JsonEscaper.escape(alias),
            JsonEscaper.escape(library),
            JsonEscaper.escape(slot)
        );
    }

    // Package-private: also called by SimulatedServerService HTTP handlers.
    static String validateContent(String content, String modeName) throws Exception {
        JsonPayloadValidator.validateValidatePayload(content);

        String signature = requireField("signature", SIGNATURE_PATTERN, content);
        String digest = requireField("inputDigestSha256", DIGEST_PATTERN, content);
        String signer = optionalField(SIGNER_PATTERN, content, "unknown");
        String expectedSignature = SIGNATURE_PREFIX + digest;
        boolean valid = expectedSignature.equals(signature);
        String reason = valid
            ? "Assinatura simulada reconhecida e consistente com o digest."
            : "A assinatura simulada nao corresponde ao digest informado.";

        return """
            {
              "status": "SUCCESS",
              "operation": "validate",
              "mode": "%s",
              "message": "Validacao simulada concluida.",
              "validationScheme": "%s",
              "valid": %s,
              "reason": "%s",
              "inputDigestSha256": "%s",
              "signature": "%s",
              "signedBy": "%s"
            }
            """.formatted(
            JsonEscaper.escape(modeName),
            JsonEscaper.escape(SIGNATURE_SCHEME),
            Boolean.toString(valid),
            JsonEscaper.escape(reason),
            JsonEscaper.escape(digest),
            JsonEscaper.escape(signature),
            JsonEscaper.escape(signer)
        );
    }

    private static String requireField(String fieldName, Pattern pattern, String input) throws ValidationException {
        Matcher matcher = pattern.matcher(input);
        if (!matcher.find()) {
            throw new ValidationException("MISSING_REQUIRED_FIELD", fieldName, "O arquivo informado para validate nao contem o campo obrigatorio \"" + fieldName + "\".");
        }

        return matcher.group(1);
    }

    private static String optionalField(Pattern pattern, String input, String fallback) {
        Matcher matcher = pattern.matcher(input);
        if (matcher.find()) {
            return matcher.group(1);
        }

        return fallback;
    }

    private static String sha256(String payload) throws Exception {
        MessageDigest digest = MessageDigest.getInstance("SHA-256");
        byte[] value = digest.digest(payload.getBytes(StandardCharsets.UTF_8));
        return HexFormat.of().withUpperCase().formatHex(value);
    }
}
