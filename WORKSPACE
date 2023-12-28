workspace(name = "com_github_dupmanio_dupman")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

RULES_PROTOBUF_EXTERNAL_TAG = "3.25.0"

RULES_PROTOBUF_EXTERNAL_SHA = "540200ef1bb101cf3f86f257f7947035313e4e485eea1f7eed9bc99dd0e2cb68"

RULES_GO_EXTERNAL_TAG = "v0.39.1"

RULES_GO_EXTERNAL_SHA = "6dc2da7ab4cf5d7bfc7c949776b1b7c733f05e56edc4bcd9022bb249d2e2a996"

RULES_GAZELLE_EXTERNAL_TAG = "v0.32.0"

RULES_GAZELLE_EXTERNAL_SHA = "29218f8e0cebe583643cbf93cae6f971be8a2484cdcfa1e45057658df8d54002"

RULES_JVM_EXTERNAL_TAG = "5.3"

RULES_JVM_EXTERNAL_SHA = "d31e369b854322ca5098ea12c69d7175ded971435e55c18dd9dd5f29cc5249ac"

RULES_OCI_EXTERNAL_TAG = "1.4.3"

RULES_OCI_EXTERNAL_SHA = "d41d0ba7855f029ad0e5ee35025f882cbe45b0d5d570842c52704f7a47ba8668"

http_archive(
    name = "com_google_protobuf",
    sha256 = RULES_PROTOBUF_EXTERNAL_SHA,
    strip_prefix = "protobuf-%s" % RULES_PROTOBUF_EXTERNAL_TAG,
    urls = [
        "https://github.com/protocolbuffers/protobuf/archive/v%s.tar.gz" % RULES_PROTOBUF_EXTERNAL_TAG,
    ],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = RULES_GO_EXTERNAL_SHA,
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/%s/rules_go-%s.zip" % (RULES_GO_EXTERNAL_TAG, RULES_GO_EXTERNAL_TAG),
        "https://github.com/bazelbuild/rules_go/releases/download/%s/rules_go-%s.zip" % (RULES_GO_EXTERNAL_TAG, RULES_GO_EXTERNAL_TAG),
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = RULES_GAZELLE_EXTERNAL_SHA,
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/%s/bazel-gazelle-%s.tar.gz" % (RULES_GAZELLE_EXTERNAL_TAG, RULES_GAZELLE_EXTERNAL_TAG),
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/%s/bazel-gazelle-%s.tar.gz" % (RULES_GAZELLE_EXTERNAL_TAG, RULES_GAZELLE_EXTERNAL_TAG),
    ],
)

http_archive(
    name = "rules_jvm_external",
    sha256 = RULES_JVM_EXTERNAL_SHA,
    strip_prefix = "rules_jvm_external-%s" % RULES_JVM_EXTERNAL_TAG,
    url = "https://github.com/bazelbuild/rules_jvm_external/releases/download/%s/rules_jvm_external-%s.tar.gz" % (RULES_JVM_EXTERNAL_TAG, RULES_JVM_EXTERNAL_TAG),
)

http_archive(
    name = "rules_oci",
    sha256 = RULES_OCI_EXTERNAL_SHA,
    strip_prefix = "rules_oci-%s" % RULES_OCI_EXTERNAL_TAG,
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v%s/rules_oci-v%s.tar.gz" % (RULES_OCI_EXTERNAL_TAG, RULES_OCI_EXTERNAL_TAG),
)

# Setup Protobuf

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# Setup GO

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains(version = "1.19")

gazelle_dependencies()

# Setup Java.

load("@rules_jvm_external//:repositories.bzl", "rules_jvm_external_deps")

rules_jvm_external_deps()

load("@rules_jvm_external//:setup.bzl", "rules_jvm_external_setup")

rules_jvm_external_setup()

load("//packages/kc-user-syncer-extension:deps.bzl", "kc_user_syncer_deps")

kc_user_syncer_deps()

load("//packages/kc-themes:deps.bzl", "kc_themes_deps")

kc_themes_deps()

# Setup OCI.

load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
)

load("@rules_oci//oci:pull.bzl", "oci_pull")

oci_pull(
    name = "distroless_base",
    digest = "sha256:ccaef5ee2f1850270d453fdf700a5392534f8d1a8ca2acda391fbb6a06b81c86",
    image = "gcr.io/distroless/base",
    platforms = [
        "linux/amd64",
        "linux/arm64",
    ],
)
