package br.ufg.runner.assinador;

import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;

public final class AssinadorApplication {
    private final PrintStream stdout;
    private final PrintStream stderr;
    private final SimulatedSigningService service;

    public AssinadorApplication() {
        this(System.out, System.err, new SimulatedSigningService());
    }

    AssinadorApplication(PrintStream stdout, PrintStream stderr, SimulatedSigningService service) {
        this.stdout = stdout;
        this.stderr = stderr;
        this.service = service;
    }

    public static void main(String[] args) {
        int exitCode = new AssinadorApplication().run(args);

        if (exitCode != ExitCode.SERVER_RUNNING.value()) {
            System.exit(exitCode);
        }
    }

    int run(String[] args) {
        try {
            CliArguments arguments = CliArguments.parse(args);
            String result = service.execute(arguments);

            if (arguments.outputPath() != null) {
                Files.writeString(arguments.outputPath(), result, StandardCharsets.UTF_8);
            }
            stdout.println(result);
            if (arguments.commandType() == CommandType.SERVER && arguments.serverCommandType() == ServerCommandType.START) {
                return ExitCode.SERVER_RUNNING.value();
            }

            return ExitCode.SUCCESS.value();
        } catch (ValidationException error) {
            stderr.println(JsonErrorResponse.validation(error));
            return ExitCode.VALIDATION_ERROR.value();
        } catch (Exception error) {
            stderr.println(JsonErrorResponse.runtime(error.getMessage()));
            return ExitCode.RUNTIME_ERROR.value();
        }
    }
}

