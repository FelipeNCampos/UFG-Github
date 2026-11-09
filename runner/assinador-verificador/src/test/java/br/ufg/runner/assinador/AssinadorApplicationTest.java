package br.ufg.runner.assinador;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.net.ServerSocket;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

public final class AssinadorApplicationTest {
    public static void main(String[] args) throws Exception {
        AssinadorApplicationTest suite = new AssinadorApplicationTest();

        suite.shouldSignAndValidateInDirectMode();
        suite.shouldSignAndValidateViaHttpMode();
        suite.shouldStartStatusAndStopServer();
        suite.shouldStartServerWithInactivityTimeout();
        suite.shouldRejectPayloadWithoutResourceType();
        suite.shouldRejectMalformedJson();
        suite.shouldRejectPayloadWithoutResourceTypeViaHttp();
        suite.shouldRejectInvalidMode();
        suite.shouldRejectMissingAlias();
        suite.shouldRejectInvalidAlias();
        suite.shouldRejectInvalidPkcs11LibExtension();
        suite.shouldRejectNonNumericPkcs11Slot();
        suite.shouldRejectValidateWithAlias();
        suite.shouldRejectDuplicateFlag();
        suite.shouldRejectUnknownFlag();
        suite.shouldRejectTimeoutOnStatusCommand();

        System.out.println("All assinador-verificador tests passed.");
    }

    // --- Fluxos principais ---

    private void shouldSignAndValidateInDirectMode() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-direct");
        Path inputPath = tempDir.resolve("entrada.json");
        Path signedPath = tempDir.resolve("assinado.json");
        Path validationPath = tempDir.resolve("validacao.json");

        Files.writeString(inputPath, """
            {
              "resourceType": "Bundle",
              "id": "bundle-001",
              "entry": []
            }
            """, StandardCharsets.UTF_8);

        InvocationResult signResult = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", signedPath.toString(),
            "--mode", "direct",
            "--alias", "test-signer",
            "--pkcs11-lib", "token.dll",
            "--pkcs11-slot", "0"
        );

        assertEquals(0, signResult.exitCode(), "sign (direct) should exit with success");
        assertTrue(Files.exists(signedPath), "sign should create output file");
        assertContains(signResult.stdout(), "\"operation\": \"sign\"", "sign stdout should contain operation");
        assertContains(signResult.stdout(), "\"mode\": \"direct\"", "sign stdout should report direct mode");
        assertContains(signResult.stdout(), "\"signatureScheme\": \"SIMULATED-SIGNATURE-V1\"", "sign stdout should expose a deterministic scheme");
        assertTrue(!signResult.stdout().contains("generatedAt"), "sign stdout should not expose volatile timestamps");
        assertContains(Files.readString(signedPath), "\"signature\": \"SIMULATED-SIGNATURE-", "output file should contain simulated signature");

        InvocationResult validateResult = run(
            "validate",
            "--pathin", signedPath.toString(),
            "--pathout", validationPath.toString(),
            "--mode", "direct"
        );

        assertEquals(0, validateResult.exitCode(), "validate (direct) should exit with success");
        assertTrue(Files.exists(validationPath), "validate should create output file");
        assertContains(validateResult.stdout(), "\"valid\": true", "validate should report signature as valid");
        assertContains(validateResult.stdout(), "\"mode\": \"direct\"", "validate stdout should report direct mode");
        assertContains(validateResult.stdout(), "\"validationScheme\": \"SIMULATED-SIGNATURE-V1\"", "validate stdout should expose a deterministic scheme");
        assertTrue(!validateResult.stdout().contains("checkedAt"), "validate stdout should not expose volatile timestamps");
    }

    private void shouldSignAndValidateViaHttpMode() throws Exception {
        int port = findFreePort();

        InvocationResult startResult = run("server", "start", "--port", Integer.toString(port));
        assertEquals(ExitCode.SERVER_RUNNING.value(), startResult.exitCode(), "server start should keep process alive");

        try {
            Path tempDir = Files.createTempDirectory("assinador-test-http");
            Path inputPath = tempDir.resolve("entrada.json");
            Path signedPath = tempDir.resolve("assinado.json");
            Path validationPath = tempDir.resolve("validacao.json");

            Files.writeString(inputPath, """
                {
                  "resourceType": "Composition",
                  "id": "comp-001"
                }
                """, StandardCharsets.UTF_8);

            InvocationResult signResult = run(
                "sign",
                "--pathin", inputPath.toString(),
                "--pathout", signedPath.toString(),
                "--mode", "http",
                "--port", Integer.toString(port),
                "--alias", "http-signer",
                "--pkcs11-lib", "token.so",
                "--pkcs11-slot", "1"
            );

            assertEquals(0, signResult.exitCode(), "sign (http) should exit with success");
            assertTrue(Files.exists(signedPath), "sign via http should create output file");
            assertContains(signResult.stdout(), "\"operation\": \"sign\"", "sign http stdout should contain operation");
            assertContains(signResult.stdout(), "\"mode\": \"http\"", "sign http stdout should report http mode");

            InvocationResult validateResult = run(
                "validate",
                "--pathin", signedPath.toString(),
                "--pathout", validationPath.toString(),
                "--mode", "http",
                "--port", Integer.toString(port)
            );

            assertEquals(0, validateResult.exitCode(), "validate (http) should exit with success");
            assertTrue(Files.exists(validationPath), "validate via http should create output file");
            assertContains(validateResult.stdout(), "\"valid\": true", "validate http should report signature as valid");
            assertContains(validateResult.stdout(), "\"mode\": \"http\"", "validate http stdout should report http mode");
        } finally {
            run("server", "stop", "--port", Integer.toString(port));
        }
    }

    private void shouldStartStatusAndStopServer() throws Exception {
        int port = findFreePort();

        InvocationResult startResult = run("server", "start", "--port", Integer.toString(port));

        try {
            assertEquals(ExitCode.SERVER_RUNNING.value(), startResult.exitCode(), "server start should keep the process alive");
            assertContains(startResult.stdout(), "\"operation\": \"server-start\"", "server start should report startup data");

            InvocationResult statusResult = run("server", "status", "--port", Integer.toString(port));

            assertEquals(0, statusResult.exitCode(), "server status should succeed while running");
            assertContains(statusResult.stdout(), "\"running\": true", "server status should report the service as running");

            InvocationResult stopResult = run("server", "stop", "--port", Integer.toString(port));

            assertEquals(0, stopResult.exitCode(), "server stop should succeed");
            assertContains(stopResult.stdout(), "\"server-stop\"", "server stop should confirm shutdown");

            waitForServerToStop(port);

            InvocationResult stoppedStatus = run("server", "status", "--port", Integer.toString(port));

            assertEquals(0, stoppedStatus.exitCode(), "server status should still succeed after stop");
            assertContains(stoppedStatus.stdout(), "\"running\": false", "server status should report the service as stopped");
        } finally {
            try {
                run("server", "stop", "--port", Integer.toString(port));
            } catch (Exception ignored) {
                // Best effort cleanup.
            }
        }
    }

    private void shouldStartServerWithInactivityTimeout() throws Exception {
        int port = findFreePort();

        InvocationResult startResult = run("server", "start", "--port", Integer.toString(port), "--timeout", "30");

        try {
            assertEquals(ExitCode.SERVER_RUNNING.value(), startResult.exitCode(), "server start with timeout should succeed");
            assertContains(startResult.stdout(), "\"operation\": \"server-start\"", "server start with timeout should report startup");
            assertContains(startResult.stdout(), "\"inactivityTimeoutMinutes\": 30", "server start should include timeout in response");
        } finally {
            run("server", "stop", "--port", Integer.toString(port));
        }
    }

    // --- Rejeicao de payloads invalidos ---

    private void shouldRejectPayloadWithoutResourceType() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-payload");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");

        Files.writeString(inputPath, "{\"id\":\"missing-resource-type\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "test-signer"
        );

        assertEquals(2, result.exitCode(), "sign should reject non-FHIR-like payloads");
        assertContains(result.stderr(), "resourceType", "stderr should mention the missing field");
        assertContains(result.stderr(), "\"code\": \"INVALID_FHIR_PAYLOAD\"", "stderr should expose a structured error code");
        assertContains(result.stderr(), "\"field\": \"resourceType\"", "stderr should expose the invalid field");
    }

    private void shouldRejectMalformedJson() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-malformed-json");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");

        Files.writeString(inputPath, "{\"resourceType\": \"Bundle\"", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "test-signer"
        );

        assertEquals(2, result.exitCode(), "sign should reject malformed JSON");
        assertContains(result.stderr(), "\"code\": \"INVALID_JSON\"", "stderr should expose malformed JSON code");
    }

    private void shouldRejectPayloadWithoutResourceTypeViaHttp() throws Exception {
        int port = findFreePort();
        run("server", "start", "--port", Integer.toString(port));

        try {
            Path tempDir = Files.createTempDirectory("assinador-test-http-payload");
            Path inputPath = tempDir.resolve("entrada.json");
            Path outputPath = tempDir.resolve("saida.json");

            Files.writeString(inputPath, "{\"id\":\"missing-resource-type\"}", StandardCharsets.UTF_8);

            InvocationResult result = run(
                "sign",
                "--pathin", inputPath.toString(),
                "--pathout", outputPath.toString(),
                "--mode", "http",
                "--port", Integer.toString(port),
                "--alias", "test-signer"
            );

            assertEquals(2, result.exitCode(), "sign via http should reject non-FHIR-like payloads");
            assertContains(result.stderr(), "resourceType", "stderr should mention the missing field");
            assertContains(result.stderr(), "\"code\": \"INVALID_FHIR_PAYLOAD\"", "stderr should expose a structured error code");
            assertContains(result.stderr(), "\"field\": \"resourceType\"", "stderr should expose the invalid field");
        } finally {
            run("server", "stop", "--port", Integer.toString(port));
        }
    }

    // --- Validacao de argumentos CLI ---

    private void shouldRejectInvalidMode() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-mode");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "one-time",
            "--alias", "test-signer"
        );

        assertEquals(2, result.exitCode(), "sign should reject unknown mode");
        assertContains(result.stderr(), "Modo nao suportado", "stderr should mention unsupported mode");
    }

    private void shouldRejectMissingAlias() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-alias");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct"
        );

        assertEquals(2, result.exitCode(), "sign without alias should fail");
        assertContains(result.stderr(), "--alias", "stderr should mention the missing --alias flag");
    }

    private void shouldRejectInvalidAlias() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-alias-invalid");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "ab"
        );

        assertEquals(2, result.exitCode(), "sign with alias too short should fail");
        assertContains(result.stderr(), "alias", "stderr should mention alias validation");
    }

    private void shouldRejectInvalidPkcs11LibExtension() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-pkcs11");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "valid-alias",
            "--pkcs11-lib", "library.jar"
        );

        assertEquals(2, result.exitCode(), "sign with invalid pkcs11 lib extension should fail");
        assertContains(result.stderr(), "PKCS#11", "stderr should mention PKCS#11");
    }

    private void shouldRejectNonNumericPkcs11Slot() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-slot");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "valid-alias",
            "--pkcs11-slot", "abc"
        );

        assertEquals(2, result.exitCode(), "sign with non-numeric pkcs11 slot should fail");
        assertContains(result.stderr(), "slot", "stderr should mention slot");
    }

    private void shouldRejectValidateWithAlias() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-validate-alias");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"signature\":\"x\",\"inputDigestSha256\":\"y\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "validate",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "forbidden-alias"
        );

        assertEquals(2, result.exitCode(), "validate with alias should fail");
        assertContains(result.stderr(), "--alias", "stderr should mention --alias is not accepted");
    }

    private void shouldRejectDuplicateFlag() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-duplicate");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "first-alias",
            "--alias", "second-alias"
        );

        assertEquals(2, result.exitCode(), "sign with duplicate flag should fail");
        assertContains(result.stderr(), "--alias", "stderr should mention the duplicate flag");
    }

    private void shouldRejectUnknownFlag() throws Exception {
        Path tempDir = Files.createTempDirectory("assinador-test-unknown");
        Path inputPath = tempDir.resolve("entrada.json");
        Path outputPath = tempDir.resolve("saida.json");
        Files.writeString(inputPath, "{\"resourceType\":\"Bundle\"}", StandardCharsets.UTF_8);

        InvocationResult result = run(
            "sign",
            "--pathin", inputPath.toString(),
            "--pathout", outputPath.toString(),
            "--mode", "direct",
            "--alias", "valid-alias",
            "--unknown-flag", "value"
        );

        assertEquals(2, result.exitCode(), "sign with unknown flag should fail");
        assertContains(result.stderr(), "Flag nao suportada", "stderr should mention unsupported flag");
    }

    private void shouldRejectTimeoutOnStatusCommand() {
        InvocationResult result = run("server", "status", "--timeout", "5");

        assertEquals(2, result.exitCode(), "server status with --timeout should fail");
        assertContains(result.stderr(), "--timeout", "stderr should mention --timeout is invalid here");
    }

    // --- Utilitarios ---

    private InvocationResult run(String... args) {
        ByteArrayOutputStream stdoutBuffer = new ByteArrayOutputStream();
        ByteArrayOutputStream stderrBuffer = new ByteArrayOutputStream();
        AssinadorApplication application = new AssinadorApplication(
            new PrintStream(stdoutBuffer, true, StandardCharsets.UTF_8),
            new PrintStream(stderrBuffer, true, StandardCharsets.UTF_8),
            new SimulatedSigningService()
        );

        int exitCode = application.run(args);

        return new InvocationResult(
            exitCode,
            stdoutBuffer.toString(StandardCharsets.UTF_8),
            stderrBuffer.toString(StandardCharsets.UTF_8)
        );
    }

    private static Path createValidInput(Path inputPath) throws Exception {
        Files.writeString(
            inputPath,
            """
                {
                  "resourceType": "Bundle",
                  "id": "bundle-001",
                  "entry": []
                }
                """,
            StandardCharsets.UTF_8
        );
        return inputPath;
    }

    private static int findFreePort() throws Exception {
        try (ServerSocket socket = new ServerSocket(0)) {
            socket.setReuseAddress(true);
            return socket.getLocalPort();
        }
    }

    private static void waitForServerToStop(int port) throws Exception {
        for (int attempt = 0; attempt < 20; attempt++) {
            InvocationResult statusResult = new AssinadorApplicationTest().run(
                "server", "status", "--port", Integer.toString(port)
            );

            if (statusResult.stdout().contains("\"running\": false")) {
                return;
            }

            Thread.sleep(50L);
        }

        throw new AssertionError("server did not stop in time");
    }

    private static void assertEquals(int expected, int actual, String message) {
        if (expected != actual) {
            throw new AssertionError(message + " (expected " + expected + ", got " + actual + ")");
        }
    }

    private static void assertEquals(CommandType expected, CommandType actual, String message) {
        if (expected != actual) {
            throw new AssertionError(message + " (expected " + expected + ", got " + actual + ")");
        }
    }

    private static void assertEquals(ServerCommandType expected, ServerCommandType actual, String message) {
        if (expected != actual) {
            throw new AssertionError(message + " (expected " + expected + ", got " + actual + ")");
        }
    }

    private static void assertTrue(boolean condition, String message) {
        if (!condition) {
            throw new AssertionError(message);
        }
    }

    private static void assertContains(String actual, String fragment, String message) {
        if (!actual.contains(fragment)) {
            throw new AssertionError(message + " (missing fragment: \"" + fragment + "\" in: " + actual.strip() + ")");
        }
    }

    private record InvocationResult(int exitCode, String stdout, String stderr) {
    }
}
