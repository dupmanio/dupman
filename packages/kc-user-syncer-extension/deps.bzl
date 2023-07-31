load("@rules_jvm_external//:defs.bzl", "maven_install")

def kc_user_syncer_deps():
    keycloak_version = "21.1.1"
    maven_install(
        artifacts = [
            "org.keycloak:keycloak-server-spi:%s" % keycloak_version,
            "org.keycloak:keycloak-server-spi-private:%s" % keycloak_version,
            "org.keycloak:keycloak-services:%s" % keycloak_version,
        ],
        repositories = [
            "https://repo1.maven.org/maven2",
        ],
    )
