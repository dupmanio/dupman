load("@rules_jvm_external//:defs.bzl", "maven_install")

def kc_themes_deps():
    maven_install(
        repositories = [
            "https://repo1.maven.org/maven2",
        ],
    )
