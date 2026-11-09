package br.ufg.runner.assinador;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.HttpURLConnection;
import java.net.InetSocketAddress;
import java.net.URI;
import java.net.URLDecoder;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicLong;

final class SimulatedServerService {
    private static final ConcurrentHashMap<Integer, ServerEntry> RUNNING_SERVERS = new ConcurrentHashMap<>();

    private SimulatedServerService() {
    }

    static String start(int port, Integer timeoutMinutes) throws Exception {
        if (RUNNING_SERVERS.containsKey(port)) {
            throw new ValidationException("Ja existe um servidor simulado rodando na porta " + port + ".");
        }

        HttpServer httpServer = HttpServer.create(new InetSocketAddress(port), 0);
        AtomicLong lastActivity = new AtomicLong(System.currentTimeMillis());
        AtomicBoolean running = new AtomicBoolean(true);
        String startedAt = Instant.now().toString();

        ServerEntry entry = new ServerEntry(httpServer, startedAt, lastActivity, running);

        httpServer.createContext("/sign", exchange -> {
            lastActivity.set(System.currentTimeMillis());
            handleSign(exchange);
        });
        httpServer.createContext("/validate", exchange -> {
            lastActivity.set(System.currentTimeMillis());
            handleValidate(exchange);
        });
        httpServer.createContext("/status", exchange -> handleStatusRequest(exchange, port, startedAt));
        httpServer.createContext("/shutdown", exchange -> handleShutdownRequest(exchange, port, entry));
        httpServer.setExecutor(null);
        httpServer.start();

        RUNNING_SERVERS.put(port, entry);

        if (timeoutMinutes != null && timeoutMinutes > 0) {
            scheduleInactivityTimeout(port, entry, timeoutMinutes);
        }

        String timeoutField = timeoutMinutes != null
            ? (",\n  \"inactivityTimeoutMinutes\": " + timeoutMinutes)
            : "";

        return """
            {
              "status": "SUCCESS",
              "operation": "server-start",
              "port": %d,
              "message": "Servidor HTTP iniciado.",
              "startedAt": "%s"%s
            }
            """.formatted(port, JsonEscaper.escape(startedAt), timeoutField);
    }

    static String status(int port) {
        ServerEntry entry = RUNNING_SERVERS.get(port);
        if (entry != null) {
            boolean running = entry.running().get();
            if (!running) {
                RUNNING_SERVERS.remove(port);
            }
            return buildStatusResponse(port, running);
        }

        return statusViaHttp(port);
    }

    static String stop(int port) throws Exception {
        ServerEntry entry = RUNNING_SERVERS.remove(port);
        if (entry != null) {
            entry.running().set(false);
            entry.server().stop(0);
            return buildStopResponse(port);
        }

        return stopViaHttp(port);
    }

    private static String statusViaHttp(int port) {
        try {
            HttpURLConnection conn = (HttpURLConnection) URI.create("http://localhost:" + port + "/status").toURL().openConnection();
            conn.setRequestMethod("GET");
            conn.setConnectTimeout(2_000);
            conn.setReadTimeout(5_000);
            int code = conn.getResponseCode();
            if (code == 200) {
                return new String(conn.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
            }
            return buildStatusResponse(port, false);
        } catch (IOException ignored) {
            // Connection refused means the server is not running.
            return buildStatusResponse(port, false);
        } catch (Exception ignored) {
            return buildStatusResponse(port, false);
        }
    }

    private static String stopViaHttp(int port) {
        try {
            HttpURLConnection conn = (HttpURLConnection) URI.create("http://localhost:" + port + "/shutdown").toURL().openConnection();
            conn.setRequestMethod("POST");
            conn.setConnectTimeout(2_000);
            conn.setReadTimeout(5_000);
            conn.getResponseCode();
        } catch (IOException ignored) {
            // Server may already be down — treat as successfully stopped.
        } catch (Exception ignored) {
            // Best effort.
        }
        return buildStopResponse(port);
    }

    private static String buildStatusResponse(int port, boolean running) {
        return """
            {
              "status": "SUCCESS",
              "operation": "server-status",
              "port": %d,
              "running": %s,
              "checkedAt": "%s"
            }
            """.formatted(port, Boolean.toString(running), JsonEscaper.escape(Instant.now().toString()));
    }

    private static String buildStopResponse(int port) {
        return """
            {
              "status": "SUCCESS",
              "operation": "server-stop",
              "port": %d,
              "message": "Servidor encerrado.",
              "stoppedAt": "%s"
            }
            """.formatted(port, JsonEscaper.escape(Instant.now().toString()));
    }

    private static void handleSign(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            sendResponse(exchange, 405, JsonErrorResponse.runtime("Metodo nao permitido. Use POST."));
            return;
        }

        try {
            Map<String, String> params = parseQueryString(exchange.getRequestURI().getQuery());
            String alias = params.get("alias");
            String pkcs11Library = params.get("pkcs11Library");
            String pkcs11Slot = params.get("pkcs11Slot");
            String content = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            String result = SimulatedSigningService.signContent(content, alias, pkcs11Library, pkcs11Slot, "http");
            sendResponse(exchange, 200, result);
        } catch (ValidationException error) {
            sendResponse(exchange, 422, JsonErrorResponse.validation(error));
        } catch (Exception error) {
            sendResponse(exchange, 500, JsonErrorResponse.runtime(error.getMessage()));
        }
    }

    private static void handleValidate(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            sendResponse(exchange, 405, JsonErrorResponse.runtime("Metodo nao permitido. Use POST."));
            return;
        }

        try {
            String content = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
            String result = SimulatedSigningService.validateContent(content, "http");
            sendResponse(exchange, 200, result);
        } catch (ValidationException error) {
            sendResponse(exchange, 422, JsonErrorResponse.validation(error));
        } catch (Exception error) {
            sendResponse(exchange, 500, JsonErrorResponse.runtime(error.getMessage()));
        }
    }

    private static void handleStatusRequest(HttpExchange exchange, int port, String startedAt) throws IOException {
        String response = """
            {
              "status": "SUCCESS",
              "operation": "server-status",
              "port": %d,
              "running": true,
              "startedAt": "%s",
              "checkedAt": "%s"
            }
            """.formatted(port, JsonEscaper.escape(startedAt), JsonEscaper.escape(Instant.now().toString()));
        sendResponse(exchange, 200, response);
    }

    private static void handleShutdownRequest(HttpExchange exchange, int port, ServerEntry entry) throws IOException {
        String stoppedAt = Instant.now().toString();
        String response = """
            {
              "status": "SUCCESS",
              "operation": "server-stop",
              "port": %d,
              "message": "Servidor encerrado via HTTP.",
              "stoppedAt": "%s"
            }
            """.formatted(port, JsonEscaper.escape(stoppedAt));
        sendResponse(exchange, 200, response);

        entry.running().set(false);
        RUNNING_SERVERS.remove(port);

        // Stop the server in a separate daemon thread so the response is flushed first.
        Thread stopThread = new Thread(() -> {
            try {
                Thread.sleep(100);
            } catch (InterruptedException ignored) {
                Thread.currentThread().interrupt();
            }
            entry.server().stop(0);
        });
        stopThread.setDaemon(true);
        stopThread.start();
    }

    private static void sendResponse(HttpExchange exchange, int statusCode, String body) throws IOException {
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=UTF-8");
        exchange.sendResponseHeaders(statusCode, bytes.length);
        try (OutputStream out = exchange.getResponseBody()) {
            out.write(bytes);
        }
    }

    private static Map<String, String> parseQueryString(String query) throws Exception {
        Map<String, String> params = new HashMap<>();
        if (query == null || query.isEmpty()) {
            return params;
        }

        for (String pair : query.split("&")) {
            int idx = pair.indexOf('=');
            if (idx < 0) {
                params.put(URLDecoder.decode(pair, StandardCharsets.UTF_8), "");
            } else {
                String key = URLDecoder.decode(pair.substring(0, idx), StandardCharsets.UTF_8);
                String value = URLDecoder.decode(pair.substring(idx + 1), StandardCharsets.UTF_8);
                params.put(key, value);
            }
        }

        return params;
    }

    private static void scheduleInactivityTimeout(int port, ServerEntry entry, int timeoutMinutes) {
        long checkIntervalMs = 30_000L;
        long timeoutMs = timeoutMinutes * 60_000L;

        ScheduledExecutorService scheduler = Executors.newSingleThreadScheduledExecutor(runnable -> {
            Thread thread = new Thread(runnable);
            thread.setDaemon(true);
            return thread;
        });

        scheduler.scheduleAtFixedRate(() -> {
            if (!entry.running().get()) {
                scheduler.shutdown();
                return;
            }

            long idleMs = System.currentTimeMillis() - entry.lastActivity().get();
            if (idleMs >= timeoutMs) {
                entry.running().set(false);
                RUNNING_SERVERS.remove(port);
                entry.server().stop(0);
                scheduler.shutdown();
            }
        }, checkIntervalMs, checkIntervalMs, TimeUnit.MILLISECONDS);
    }

    private static final class ServerEntry {
        private final HttpServer server;
        private final String startedAt;
        private final AtomicLong lastActivity;
        private final AtomicBoolean running;

        ServerEntry(HttpServer server, String startedAt, AtomicLong lastActivity, AtomicBoolean running) {
            this.server = server;
            this.startedAt = startedAt;
            this.lastActivity = lastActivity;
            this.running = running;
        }

        HttpServer server() {
            return server;
        }

        String startedAt() {
            return startedAt;
        }

        AtomicLong lastActivity() {
            return lastActivity;
        }

        AtomicBoolean running() {
            return running;
        }
    }
}
