package io.dupman;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;
import org.jboss.logging.Logger;
import org.keycloak.events.Event;
import org.keycloak.events.EventListenerProvider;
import org.keycloak.events.admin.AdminEvent;
import org.keycloak.events.admin.OperationType;
import org.keycloak.events.admin.ResourceType;
import org.keycloak.models.KeycloakSession;
import org.keycloak.models.RealmModel;
import org.keycloak.models.UserModel;

import java.io.IOException;
import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Map;

public class UserSyncerProvider implements EventListenerProvider {

    private static final Logger log = Logger.getLogger(UserSyncerProvider.class);
    private final KeycloakSession session;
    private final HttpClient httpClient;
    private final URI oauthURL;
    private final String oauthCredentials;
    private final URI syncEndpoint;

    public UserSyncerProvider(KeycloakSession session) {
        this.session = session;
        this.oauthURL = URI.create(System.getenv("USER_SYNC_OAUTH_URL"));
        this.oauthCredentials = System.getenv("USER_SYNC_OAUTH_CREDENTIALS");
        this.syncEndpoint = URI.create(System.getenv("USER_SYNC_ENDPOINT"));
        this.httpClient = HttpClient.newBuilder()
                .followRedirects(HttpClient.Redirect.ALWAYS)
                .connectTimeout(Duration.ofSeconds(5))
                .build();
    }

    @Override
    public void onEvent(Event event) {
        String realmId = event.getRealmId();
        String userId = event.getUserId();

        switch (event.getType()) {
            case REGISTER:
                processUser(OperationType.CREATE, realmId, userId);
                break;
            case UPDATE_PROFILE:
                processUser(OperationType.UPDATE, realmId, userId);
                break;
            default:
                break;
        }
    }

    private void processUser(OperationType operationType, String realmId, String userId) {
        RealmModel realm = this.session.realms().getRealm(realmId);
        UserModel user = this.session.users().getUserById(realm, userId);

        try {
            String token = this.getAuthToken();

            switch (operationType) {
                case CREATE:
                    this.createUser(user, token);
                    break;
                case UPDATE:
                    this.updateUser(user, token);
                    break;
                default:
                    break;
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private String getAuthToken() throws Exception {
        Map<String, String> data = Map.of(
                "grant_type", "client_credentials",
                "scope", "user:create user:update"
        );

        String payload = data.entrySet().stream()
                .map(entry -> URLEncoder.encode(entry.getKey(), StandardCharsets.UTF_8) +
                        "=" + URLEncoder.encode(entry.getValue(), StandardCharsets.UTF_8))
                .reduce((s1, s2) -> s1 + "&" + s2)
                .orElse("");

        HttpRequest request = HttpRequest.newBuilder()
                .uri(this.oauthURL)
                .timeout(Duration.ofSeconds(5))
                .header("Content-Type", "application/x-www-form-urlencoded")
                .header("Authorization", "Basic " + this.oauthCredentials)
                .POST(HttpRequest.BodyPublishers.ofString(payload))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        String responseBody = response.body();

        if (response.statusCode() == 200) {
            return (new ObjectMapper())
                    .readTree(responseBody)
                    .get("access_token")
                    .asText();
        }

        throw new Exception("Unable to authenticate service account");
    }

    private void createUser(UserModel user, String token) throws IOException, InterruptedException {
        String payload = this.userToJson(user).toString();
        HttpRequest request = HttpRequest.newBuilder()
                .uri(this.syncEndpoint)
                .timeout(Duration.ofSeconds(5))
                .header("Content-Type", "application/json")
                .header("Authorization", "Bearer " + token)
                .POST(HttpRequest.BodyPublishers.ofString(payload))
                .build();

        // @todo: handle response.
        this.httpClient.send(request, HttpResponse.BodyHandlers.ofString());
    }

    private ObjectNode userToJson(UserModel user) {
        return (new ObjectMapper()).createObjectNode()
                .put("id", user.getId())
                .put("username", user.getUsername())
                .put("email", user.getEmail())
                .put("firstName", user.getFirstName())
                .put("lastName", user.getLastName());
    }

    private void updateUser(UserModel user, String token) throws IOException, InterruptedException {
        String payload = this.userToJson(user).toString();
        HttpRequest request = HttpRequest.newBuilder()
                .uri(this.syncEndpoint)
                .timeout(Duration.ofSeconds(5))
                .header("Content-Type", "application/json")
                .header("Authorization", "Bearer " + token)
                .method("PATCH", HttpRequest.BodyPublishers.ofString(payload))
                .build();

        // @todo: handle response.
        this.httpClient.send(request, HttpResponse.BodyHandlers.ofString());
    }

    @Override
    public void onEvent(AdminEvent event, boolean isPredefinedEvent) {
        if (event.getResourceType().equals(ResourceType.USER) && isPredefinedEvent) {
            this.processUser(event.getOperationType(), event.getRealmId(), this.extractIdFromResourcePath(event.getResourcePath()));
        }
    }

    private String extractIdFromResourcePath(String path) {
        return path.substring(path.lastIndexOf('/') + 1);
    }

    @Override
    public void close() {

    }

}
