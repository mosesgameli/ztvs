import java.util.Scanner;
import java.util.UUID;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);
        if (scanner.hasNextLine()) {
            String line = scanner.nextLine();
            // Simplified JSON parsing for example purposes
            if (line.contains("\"method\":\"handshake\"")) {
                String id = extractId(line);
                String resp = "{" +
                        "\"jsonrpc\":\"2.0\"," +
                        "\"id\":\"" + id + "\"," +
                        "\"result\":{" +
                        "\"name\":\"example-java\"," +
                        "\"version\":\"1.0.0\"," +
                        "\"api_version\":1," +
                        "\"checks_supported\":[\"java_check\"]" +
                        "}}";
                System.out.println(resp);
            } else if (line.contains("\"method\":\"run_check\"")) {
                String id = extractId(line);
                String resp = "{" +
                        "\"jsonrpc\":\"2.0\"," +
                        "\"id\":\"" + id + "\"," +
                        "\"result\":{" +
                        "\"status\":\"pass\"," +
                        "\"finding\":{" +
                        "\"id\":\"F-JAVA-001\"," +
                        "\"check_id\":\"java_check\"," +
                        "\"severity\":\"info\"," +
                        "\"title\":\"Java Plugin Running\"," +
                        "\"description\":\"Java-based plugin is communicating correctly.\"," +
                        "\"evidence\":{\"runtime\":\"java\",\"version\":\"" + System.getProperty("java.version")
                        + "\"}," +
                        "\"remediation\":\"None\"" +
                        "}}}";
                System.out.println(resp);
            }
        }
    }

    private static String extractId(String json) {
        // Very basic extraction for example purposes
        int idIndex = json.indexOf("\"id\":\"");
        if (idIndex == -1)
            return "1";
        int start = idIndex + 6;
        int end = json.indexOf("\"", start);
        return json.substring(start, end);
    }
}
