load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "co_honnef_go_tools",
        importpath = "honnef.co/go/tools",
        sum = "h1:UoveltGrhghAA7ePc+e+QYDHXrBps2PqFZiHkGR/xK8=",
        version = "v0.0.1-2020.1.4",
    )
    go_repository(
        name = "com_github_antihax_optional",
        importpath = "github.com/antihax/optional",
        sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_armon_go_metrics",
        importpath = "github.com/armon/go-metrics",
        sum = "h1:yCQqn7dwca4ITXb+CbubHmedzaQYHhNhrEXLYUeEe8Q=",
        version = "v0.4.0",
    )

    go_repository(
        name = "com_github_benbjohnson_clock",
        importpath = "github.com/benbjohnson/clock",
        sum = "h1:ip6w0uFQkncKQ979AypyG0ER7mqUSBdKLOgAle/AT8A=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_bsm_ginkgo_v2",
        importpath = "github.com/bsm/ginkgo/v2",
        sum = "h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=",
        version = "v2.12.0",
    )
    go_repository(
        name = "com_github_bsm_gomega",
        importpath = "github.com/bsm/gomega",
        sum = "h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=",
        version = "v1.27.10",
    )

    go_repository(
        name = "com_github_burntsushi_toml",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_burntsushi_xgb",
        importpath = "github.com/BurntSushi/xgb",
        sum = "h1:1BDTz0u9nC3//pOCMdNH+CiXJVYJh5UQNCOBG7jbELc=",
        version = "v0.0.0-20160522181843-27f122750802",
    )

    go_repository(
        name = "com_github_bytedance_sonic",
        importpath = "github.com/bytedance/sonic",
        sum = "h1:6iJ6NqdoxCDr6mbY8h18oSO+cShGSMRGCEo7F2h0x8s=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_cenkalti_backoff_v4",
        importpath = "github.com/cenkalti/backoff/v4",
        sum = "h1:y4OZtCnogmCPw98Zjyt5a6+QwPLGkiQsYW5oUqylYbM=",
        version = "v4.2.1",
    )

    go_repository(
        name = "com_github_census_instrumentation_opencensus_proto",
        importpath = "github.com/census-instrumentation/opencensus-proto",
        sum = "h1:iKLQ0xPNFxR/2hzXZMrBo8f1j86j5WHzznCCQxV/b8g=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:DC2CZ1Ep5Y4k3ZQ899DldepgrayRUGE6BBZ/cd9Cj44=",
        version = "v2.2.0",
    )

    go_repository(
        name = "com_github_chenzhuoyu_base64x",
        importpath = "github.com/chenzhuoyu/base64x",
        sum = "h1:qSGYFH7+jGhDF8vLC+iwCD4WpbV1EBDSzWkJODFLams=",
        version = "v0.0.0-20221115062448-fe3a3abad311",
    )
    go_repository(
        name = "com_github_chromedp_cdproto",
        importpath = "github.com/chromedp/cdproto",
        sum = "h1:aPflPkRFkVwbW6dmcVqfgwp1i+UWGFH6VgR1Jim5Ygc=",
        version = "v0.0.0-20230802225258-3cf4e6d46a89",
    )
    go_repository(
        name = "com_github_chromedp_chromedp",
        importpath = "github.com/chromedp/chromedp",
        sum = "h1:dKtNz4kApb06KuSXoTQIyUC2TrA0fhGDwNZf3bcgfKw=",
        version = "v0.9.2",
    )
    go_repository(
        name = "com_github_chromedp_sysutil",
        importpath = "github.com/chromedp/sysutil",
        sum = "h1:+ZxhTpfpZlmchB58ih/LBHX52ky7w2VhQVKQMucy3Ic=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_chzyer_logex",
        importpath = "github.com/chzyer/logex",
        sum = "h1:Swpa1K6QvQznwJRcfTfQJmTE72DqScAa40E+fbHEXEE=",
        version = "v1.1.10",
    )
    go_repository(
        name = "com_github_chzyer_readline",
        importpath = "github.com/chzyer/readline",
        sum = "h1:fY5BOSpyZCqRo5OhCuC+XN+r/bBCmeuuJtjz+bCNIf8=",
        version = "v0.0.0-20180603132655-2972be24d48e",
    )
    go_repository(
        name = "com_github_chzyer_test",
        importpath = "github.com/chzyer/test",
        sum = "h1:q763qf9huN11kDQavWsoZXJNW3xEE4JJyHa5Q25/sd8=",
        version = "v0.0.0-20180213035817-a1ea475d72b1",
    )
    go_repository(
        name = "com_github_client9_misspell",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_cncf_udpa_go",
        importpath = "github.com/cncf/udpa/go",
        sum = "h1:QQ3GSy+MqSHxm/d8nCtnAiZdYFd45cYZPs8vOOIYKfk=",
        version = "v0.0.0-20220112060539-c52dc94e7fbe",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:/inchEIKaYC1Akx+H+gqO04wryn5h75LSazbRlnya1k=",
        version = "v0.0.0-20230607035331-e9ce68804cb4",
    )

    go_repository(
        name = "com_github_coreos_go_oidc_v3",
        importpath = "github.com/coreos/go-oidc/v3",
        sum = "h1:AKVxfYw1Gmkn/w96z0DbT/B/xFnzTd3MkZvWLjF4n/o=",
        version = "v3.6.0",
    )

    go_repository(
        name = "com_github_coreos_go_semver",
        importpath = "github.com/coreos/go-semver",
        sum = "h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_coreos_go_systemd_v22",
        importpath = "github.com/coreos/go-systemd/v22",
        sum = "h1:D9/bQk5vlXQFZ6Kwuu6zaiXJ9oTPe68++AzAJc1DzSI=",
        version = "v22.3.2",
    )

    go_repository(
        name = "com_github_creack_pty",
        importpath = "github.com/creack/pty",
        sum = "h1:uDmaGzcdjhF4i/plgjmEsriH11Y0o7RKapEf/LDaM3w=",
        version = "v1.1.9",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_dgryski_go_rendezvous",
        importpath = "github.com/dgryski/go-rendezvous",
        sum = "h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=",
        version = "v0.0.0-20200823014737-9f7001d12a5f",
    )

    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:wSUXTlLfiAQRWs2F+p+EKOY9rUyis1MyGqJ2DIk5HpM=",
        version = "v0.11.1",
    )
    go_repository(
        name = "com_github_envoyproxy_protoc_gen_validate",
        importpath = "github.com/envoyproxy/protoc-gen-validate",
        sum = "h1:QkIBuU5k+x7/QXPvPPnWXWlCdaBFApVqftFV6k087DA=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_fatih_color",
        importpath = "github.com/fatih/color",
        sum = "h1:8LOYc1KYPPmyKMuN8QV2DNRWNbLo6LZ0iLs8+mlH53w=",
        version = "v1.13.0",
    )

    go_repository(
        name = "com_github_frankban_quicktest",
        importpath = "github.com/frankban/quicktest",
        sum = "h1:g2rn0vABPOOXmZUj+vbmUp0lPoXEMuhTpIluN0XL9UY=",
        version = "v1.14.4",
    )
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:n+5WquG0fcWoWp6xPWfHdbskMCQaFnG6PfBrh1Ky4HY=",
        version = "v1.6.0",
    )

    go_repository(
        name = "com_github_gabriel_vasile_mimetype",
        importpath = "github.com/gabriel-vasile/mimetype",
        sum = "h1:w5qFW6JKBz9Y393Y4q372O9A7cUSequkh1Q7OhCmWKU=",
        version = "v1.4.2",
    )

    go_repository(
        name = "com_github_gin_contrib_sse",
        importpath = "github.com/gin-contrib/sse",
        sum = "h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gin_contrib_zap",
        importpath = "github.com/gin-contrib/zap",
        sum = "h1:RMSFFJo34XZogV62OgOzvrlaMNmXrNxmJ3bFmMwl6Cc=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gin_gonic_gin",
        importpath = "github.com/gin-gonic/gin",
        sum = "h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_go_gl_glfw",
        importpath = "github.com/go-gl/glfw",
        sum = "h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=",
        version = "v0.0.0-20190409004039-e6da0acd62b1",
    )
    go_repository(
        name = "com_github_go_gl_glfw_v3_3_glfw",
        importpath = "github.com/go-gl/glfw/v3.3/glfw",
        sum = "h1:WtGNWLvXpe6ZudgnXrq0barxBImvnnJoMEhXAzcbM0I=",
        version = "v0.0.0-20200222043503-6f7a984d4dc4",
    )
    go_repository(
        name = "com_github_go_jose_go_jose_v3",
        importpath = "github.com/go-jose/go-jose/v3",
        sum = "h1:s6rrhirfEP/CGIoc6p+PZAeogN2SxKav6Wp7+dyMWVo=",
        version = "v3.0.0",
    )

    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:g01GSCwiDw2xSZfjJ2/T9M+S6pFdcNtFYsp+Y43HYDQ=",
        version = "v1.2.4",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_go_playground_assert_v2",
        importpath = "github.com/go-playground/assert/v2",
        sum = "h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=",
        version = "v2.2.0",
    )
    go_repository(
        name = "com_github_go_playground_locales",
        importpath = "github.com/go-playground/locales",
        sum = "h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_go_playground_universal_translator",
        importpath = "github.com/go-playground/universal-translator",
        sum = "h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=",
        version = "v0.18.1",
    )
    go_repository(
        name = "com_github_go_playground_validator_v10",
        importpath = "github.com/go-playground/validator/v10",
        sum = "h1:vgvQWe3XCz3gIeFDm/HnTIbj6UGmg/+t63MyGU2n5js=",
        version = "v10.14.0",
    )
    go_repository(
        name = "com_github_go_resty_resty_v2",
        importpath = "github.com/go-resty/resty/v2",
        sum = "h1:me+K9p3uhSmXtrBZ4k9jcEAfJmuC8IivWHwaLZwPrFY=",
        version = "v2.7.0",
    )
    go_repository(
        name = "com_github_gobwas_httphead",
        importpath = "github.com/gobwas/httphead",
        sum = "h1:exrUm0f4YX0L7EBwZHuCF4GDp8aJfVeBrlLQrs6NqWU=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gobwas_pool",
        importpath = "github.com/gobwas/pool",
        sum = "h1:xfeeEhW7pwmX8nuLVlqbzVc7udMDrwetjEv+TZIz1og=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_github_gobwas_ws",
        importpath = "github.com/gobwas/ws",
        sum = "h1:F2aeBZrm2NDsc7vbovKrWSogd4wvfAxg0FQ89/iqOTk=",
        version = "v1.2.1",
    )

    go_repository(
        name = "com_github_goccy_go_json",
        importpath = "github.com/goccy/go-json",
        sum = "h1:CrxCmQqYDkv1z7lO7Wbh2HN93uovUHgrECaO5ZrCXAU=",
        version = "v0.10.2",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        importpath = "github.com/gogo/protobuf",
        sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_golang_glog",
        importpath = "github.com/golang/glog",
        sum = "h1:/d3pCKDPWNnvIWe0vVUpNP32qc8U3PDVxySP/y360qE=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt_v5",
        importpath = "github.com/golang-jwt/jwt/v5",
        sum = "h1:1n1XNM9hk7O9mnQoNBGolZvzebBQ7p93ULHRc28XJUE=",
        version = "v5.0.0",
    )

    go_repository(
        name = "com_github_golang_mock",
        importpath = "github.com/golang/mock",
        sum = "h1:l75CXGRSwbaYNpl/Z2X1XIIAMSCquvXgpVZDhwEIJsc=",
        version = "v1.4.4",
    )

    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_github_google_btree",
        importpath = "github.com/google/btree",
        sum = "h1:0udJVsspx3VBr5FwtLhQQtuAsVc79tTq0ocGIPAU6qo=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=",
        version = "v0.5.9",
    )
    go_repository(
        name = "com_github_google_gofuzz",
        importpath = "github.com/google/gofuzz",
        sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_google_martian",
        importpath = "github.com/google/martian",
        sum = "h1:/CP5g8u/VJHijgedC/Legn3BAbAaWPgecwXBIDzw5no=",
        version = "v2.1.0+incompatible",
    )
    go_repository(
        name = "com_github_google_martian_v3",
        importpath = "github.com/google/martian/v3",
        sum = "h1:wCKgOCHuUEVfsaQLpPSJb7VdYCdTVZQAuOdYm1yc/60=",
        version = "v3.1.0",
    )
    go_repository(
        name = "com_github_google_pprof",
        importpath = "github.com/google/pprof",
        sum = "h1:LR89qFljJ48s990kEKGsk213yIJDPI4205OKOzbURK8=",
        version = "v0.0.0-20201218002935-b9804c9f04c2",
    )
    go_repository(
        name = "com_github_google_renameio",
        importpath = "github.com/google/renameio",
        sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_google_s2a_go",
        importpath = "github.com/google/s2a-go",
        sum = "h1:FAgZmpLl/SXurPEZyCMPBIiiYeTbqfjlbdnCNTAkbGE=",
        version = "v0.1.3",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:yk9/cqRKtT9wXZSsRH9aurXEpJX+U6FLtpYTdC3R06k=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:UBtEZqx1bjXtOQ5BVTkuYghXrr3N4V123VKJK67vJZc=",
        version = "v2.8.0",
    )
    go_repository(
        name = "com_github_googleapis_google_cloud_go_testing",
        importpath = "github.com/googleapis/google-cloud-go-testing",
        sum = "h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=",
        version = "v0.0.0-20200911160855-bcd43fbb19e8",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        importpath = "github.com/grpc-ecosystem/grpc-gateway/v2",
        sum = "h1:YBftPWNWd4WwGqtY2yeZL2ef8rHAxPBD8KFhJpmcqms=",
        version = "v2.16.0",
    )

    go_repository(
        name = "com_github_hashicorp_consul_api",
        importpath = "github.com/hashicorp/consul/api",
        sum = "h1:9IHTjNVSZ7MIwjlW3N3a7iGiykCMDpxZu8jsxFJh0yc=",
        version = "v1.20.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_cleanhttp",
        importpath = "github.com/hashicorp/go-cleanhttp",
        sum = "h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_hclog",
        importpath = "github.com/hashicorp/go-hclog",
        sum = "h1:La19f8d7WIlm4ogzNHB0JGqs5AUDAZ2UfCY4sJXcJdM=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_immutable_radix",
        importpath = "github.com/hashicorp/go-immutable-radix",
        sum = "h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_rootcerts",
        importpath = "github.com/hashicorp/go-rootcerts",
        sum = "h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru",
        importpath = "github.com/hashicorp/golang-lru",
        sum = "h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=",
        version = "v0.5.4",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_serf",
        importpath = "github.com/hashicorp/serf",
        sum = "h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=",
        version = "v0.10.1",
    )
    go_repository(
        name = "com_github_ianlancetaylor_demangle",
        importpath = "github.com/ianlancetaylor/demangle",
        sum = "h1:mV02weKRL81bEnm8A0HT1/CAelMQDBuQIfLw8n+d6xI=",
        version = "v0.0.0-20200824232613-28f6c0f3b639",
    )
    go_repository(
        name = "com_github_jackc_pgpassfile",
        importpath = "github.com/jackc/pgpassfile",
        sum = "h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_jackc_pgservicefile",
        importpath = "github.com/jackc/pgservicefile",
        sum = "h1:bbPeKD0xmW/Y25WS6cokEszi5g+S0QxI/d45PkRi7Nk=",
        version = "v0.0.0-20221227161230-091c0ba34f0a",
    )
    go_repository(
        name = "com_github_jackc_pgx_v5",
        importpath = "github.com/jackc/pgx/v5",
        sum = "h1:Fcr8QJ1ZeLi5zsPZqQeUZhNhxfkkKBOgJuYkJHoBOtU=",
        version = "v5.3.1",
    )
    go_repository(
        name = "com_github_jackc_puddle_v2",
        importpath = "github.com/jackc/puddle/v2",
        sum = "h1:RdcDk92EJBuBS55nQMMYFXTxwstHug4jkhT5pq8VxPk=",
        version = "v2.2.0",
    )
    go_repository(
        name = "com_github_jinzhu_copier",
        importpath = "github.com/jinzhu/copier",
        sum = "h1:GlvfUwHk62RokgqVNvYsku0TATCF7bAHVwEXoBh3iJg=",
        version = "v0.3.5",
    )

    go_repository(
        name = "com_github_jinzhu_inflection",
        importpath = "github.com/jinzhu/inflection",
        sum = "h1:K317FqzuhWc8YvSVlFMCCUb36O/S9MCKRDI7QkRKD/E=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_jinzhu_now",
        importpath = "github.com/jinzhu/now",
        sum = "h1:/o9tlHleP7gOFmsnYNz3RGnqzefHA47wQpKrrdTIwXQ=",
        version = "v1.1.5",
    )
    go_repository(
        name = "com_github_josharian_intern",
        importpath = "github.com/josharian/intern",
        sum = "h1:vlS4z54oSdjm0bgjRigI+G1HpF+tI+9rE5LLzOg8HmY=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_jstemmer_go_junit_report",
        importpath = "github.com/jstemmer/go-junit-report",
        sum = "h1:6QPYqodiu3GuPL+7mfx+NwDdp2eTkp9IfEUpgAwUN0o=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_kisielk_gotool",
        importpath = "github.com/kisielk/gotool",
        sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_klauspost_cpuid_v2",
        importpath = "github.com/klauspost/cpuid/v2",
        sum = "h1:acbojRNwl3o09bUq+yDCtZFc1aiwaAAxtcn8YkZXnvk=",
        version = "v2.2.4",
    )
    go_repository(
        name = "com_github_kr_fs",
        importpath = "github.com/kr/fs",
        sum = "h1:Jskdu9ieNAYnjxsi0LbQp1ulIKZV1LAFgK1tWhpZgl8=",
        version = "v0.1.0",
    )

    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_pty",
        importpath = "github.com/kr/pty",
        sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_ledongthuc_pdf",
        importpath = "github.com/ledongthuc/pdf",
        sum = "h1:6Yzfa6GP0rIo/kULo2bwGEkFvCePZ3qHDDTC3/J9Swo=",
        version = "v0.0.0-20220302134840-0c2507a12d80",
    )

    go_repository(
        name = "com_github_leodido_go_urn",
        importpath = "github.com/leodido/go-urn",
        sum = "h1:XlAE/cm/ms7TE/VMVoduSpNBoyc2dOxHs5MZSwAN63Q=",
        version = "v1.2.4",
    )
    go_repository(
        name = "com_github_magiconair_properties",
        importpath = "github.com/magiconair/properties",
        sum = "h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_github_mailru_easyjson",
        importpath = "github.com/mailru/easyjson",
        sum = "h1:UGYAvKxe3sBsEDzO8ZeWOSlIQfWFlxbzLZe7hwFURr0=",
        version = "v0.7.7",
    )

    go_repository(
        name = "com_github_mattn_go_colorable",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:jF+Du6AlPIjs2BiUiQlKOX0rt3SujHxPnksPKZbaA40=",
        version = "v0.1.12",
    )

    go_repository(
        name = "com_github_mattn_go_isatty",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:JITubQf0MOLdlGRuRq+jtsDlekdYPia9ZFsB8h/APPA=",
        version = "v0.0.19",
    )
    go_repository(
        name = "com_github_mattn_go_sqlite3",
        importpath = "github.com/mattn/go-sqlite3",
        sum = "h1:vfoHhTN1af61xCRSWzFIWzx2YskyMTwHLrExkBOjvxI=",
        version = "v1.14.15",
    )

    go_repository(
        name = "com_github_mcuadros_go_defaults",
        importpath = "github.com/mcuadros/go-defaults",
        sum = "h1:FODb8WSf0uGaY8elWJAkoLL0Ri6AlZ1bFlenk56oZtc=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_mitchellh_go_homedir",
        importpath = "github.com/mitchellh/go-homedir",
        sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_mitchellh_mapstructure",
        importpath = "github.com/mitchellh/mapstructure",
        sum = "h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=",
        version = "v1.5.0",
    )

    go_repository(
        name = "com_github_modern_go_concurrent",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_niemeyer_pretty",
        importpath = "github.com/niemeyer/pretty",
        sum = "h1:fD57ERR4JtEqsWbfPhv4DMiApHyliiK5xCTNVSPiaAs=",
        version = "v0.0.0-20200227124842-a10e7caefd8e",
    )
    go_repository(
        name = "com_github_orisano_pixelmatch",
        importpath = "github.com/orisano/pixelmatch",
        sum = "h1:x0TT0RDC7UhAVbbWWBzr41ElhJx5tXPWkIHA2HWPRuw=",
        version = "v0.0.0-20220722002657-fb0b55479cde",
    )

    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:0ctb6s9mE31h0/lhu+J6OPmVeDxJn+kYnJc2jZR9tGQ=",
        version = "v2.0.8",
    )
    go_repository(
        name = "com_github_pkg_diff",
        importpath = "github.com/pkg/diff",
        sum = "h1:aoZm08cpOy4WuID//EZDgcC4zIxODThtZNPirFr42+A=",
        version = "v0.0.0-20210226163009-20ebb0f2a09e",
    )
    go_repository(
        name = "com_github_pkg_errors",
        importpath = "github.com/pkg/errors",
        sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_pkg_sftp",
        importpath = "github.com/pkg/sftp",
        sum = "h1:I2qBYMChEhIjOgazfJmV3/mZM256btk6wkCDRmW7JYs=",
        version = "v1.13.1",
    )

    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
        version = "v0.0.0-20190812154241-14fe0d1b01d4",
    )
    go_repository(
        name = "com_github_rabbitmq_amqp091_go",
        importpath = "github.com/rabbitmq/amqp091-go",
        sum = "h1:qrQtyzB4H8BQgEuJwhmVQqVHB9O4+MNDJCCAcpc3Aoo=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_redis_go_redis_extra_rediscmd_v9",
        importpath = "github.com/redis/go-redis/extra/rediscmd/v9",
        sum = "h1:EaDatTxkdHG+U3Bk4EUr+DZ7fOGwTfezUiUJMaIcaho=",
        version = "v9.0.5",
    )
    go_repository(
        name = "com_github_redis_go_redis_extra_redisotel_v9",
        importpath = "github.com/redis/go-redis/extra/redisotel/v9",
        sum = "h1:EfpWLLCyXw8PSM2/XNJLjI3Pb27yVE+gIAfeqp8LUCc=",
        version = "v9.0.5",
    )

    go_repository(
        name = "com_github_redis_go_redis_v9",
        importpath = "github.com/redis/go-redis/v9",
        sum = "h1:WlYJg71ODF0dVspZZCpYmoF1+U1Jjk9Rwd7pq6QmlCg=",
        version = "v9.2.1",
    )
    go_repository(
        name = "com_github_rogpeppe_fastuuid",
        importpath = "github.com/rogpeppe/fastuuid",
        sum = "h1:Ppwyp6VYCF1nvBTXL3trRso7mXMlRrw9ooo375wvi2s=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:TMyTOH3F/DB16zRVcYyreMH6GnZZrwQVAoYjRBZyWFQ=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_crypt",
        importpath = "github.com/sagikazarmark/crypt",
        sum = "h1:96E1qrToLBU6fGzo+PRRz7KGOc9FkYFiPnR3/zf8Smg=",
        version = "v0.10.0",
    )
    go_repository(
        name = "com_github_sirupsen_logrus",
        importpath = "github.com/sirupsen/logrus",
        sum = "h1:oxx1eChJGI6Uks2ZC4W1zpLlVgqB8ner4EuQwV4Ik1Y=",
        version = "v1.9.2",
    )

    go_repository(
        name = "com_github_spf13_afero",
        importpath = "github.com/spf13/afero",
        sum = "h1:stMpOSZFs//0Lv29HduCmli3GUfpFoF3Y1Q/aXj/wVM=",
        version = "v1.9.5",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:R+kOtfhWQE6TVQzY+4D7wJLBgkdVasCEFxSUBYBYIlA=",
        version = "v1.5.1",
    )
    go_repository(
        name = "com_github_spf13_jwalterweatherman",
        importpath = "github.com/spf13/jwalterweatherman",
        sum = "h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:rGGH0XDZhdUOryiDWjmIvUSWpbNqisK8Wk0Vyefw8hc=",
        version = "v1.16.0",
    )

    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:CcVxjf3Q8PM0mHUKJCdn+eZZtm5yQwehR5yeSVQQcUk=",
        version = "v1.8.4",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:X1TuBLAMDFbaTAChgCBLu3DU3UPyELpnF2jjJ2cz/S8=",
        version = "v1.4.2",
    )

    go_repository(
        name = "com_github_twitchyliquid64_golang_asm",
        importpath = "github.com/twitchyliquid64/golang-asm",
        sum = "h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=",
        version = "v0.15.1",
    )
    go_repository(
        name = "com_github_ugorji_go",
        importpath = "github.com/ugorji/go",
        sum = "h1:qYhyWUUd6WbiM+C6JZAUkIJt/1WrjzNHY9+KCIjVqTo=",
        version = "v1.2.7",
    )
    go_repository(
        name = "com_github_ugorji_go_codec",
        importpath = "github.com/ugorji/go/codec",
        sum = "h1:BMaWp1Bb6fHwEtbplGBGJ498wD+LKlNSl25MjdZY4dU=",
        version = "v1.2.11",
    )

    go_repository(
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:dPmz1Snjq0kmkz159iL7S6WzdahUTHnHB5M56WFVifs=",
        version = "v1.3.5",
    )
    go_repository(
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:Zc8gqp3+a9/Eyph2KDmcGaPtbKRIoqq4YTlL4NMD0Ys=",
        version = "v0.110.0",
    )
    go_repository(
        name = "com_google_cloud_go_bigquery",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:PQcPefKFdaIzjQFbiyOgAqyx8q5djaE7x9Sqe712DPA=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:JNBsyXVoOoNJtTQcnEY5uYpZIbeCTYIeDe0Xh1bySMk=",
        version = "v1.21.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:mg4jlk7mCAj6xXp9UJ4fjI9VUI5rubuGBW5aJ7UnBMY=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_google_cloud_go_datastore",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:/May9ojXjRkPBNVrq+oWLqmWCkr4OU5uRY29bu0mRyQ=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:IBlRyxgGySXu5VuW0RgGFlTtLukSnNkpDiEOMkQkmpA=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:v+yFJOfKC3yZdY6ZUI933pIYdhyhV8S3NpWrXWmg7jM=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:ukjixP1wl0LpnZ6LWtZJ0mX5tBmjp1f8Sqer8Z2OMUU=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:6RRlFMv1omScs6iq2hfE3IvgE+l6RfJPampq8UZc5TU=",
        version = "v1.14.0",
    )
    go_repository(
        name = "com_shuralyov_dmitri_gpu_mtl",
        importpath = "dmitri.shuralyov.com/gpu/mtl",
        sum = "h1:VpgP7xuJadIUuKccphEpTJnWhS2jkQyMt6Y7pJCD7fY=",
        version = "v0.0.0-20190408044501-666a987793e9",
    )
    go_repository(
        name = "in_gopkg_alexcesaro_quotedprintable_v3",
        importpath = "gopkg.in/alexcesaro/quotedprintable.v3",
        sum = "h1:2gGKlE2+asNV9m7xrywl36YYNnBG5ZQ0r/BOOxqPpmk=",
        version = "v3.0.0-20150716171945-2caba252f4dc",
    )

    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_errgo_v2",
        importpath = "gopkg.in/errgo.v2",
        sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
        version = "v2.1.0",
    )
    go_repository(
        name = "in_gopkg_ini_v1",
        importpath = "gopkg.in/ini.v1",
        sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
        version = "v1.67.0",
    )
    go_repository(
        name = "in_gopkg_mail_v2",
        importpath = "gopkg.in/mail.v2",
        sum = "h1:WYFn/oANrAGP2C0dcV6/pbkPzv8yGzqTjPmTeO7qoXk=",
        version = "v2.3.1",
    )

    go_repository(
        name = "in_gopkg_yaml_v2",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_etcd_go_etcd_api_v3",
        importpath = "go.etcd.io/etcd/api/v3",
        sum = "h1:4wSsluwyTbGGmyjJktOf3wFQoTBIURXHnq9n/G/JQHs=",
        version = "v3.5.9",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_pkg_v3",
        importpath = "go.etcd.io/etcd/client/pkg/v3",
        sum = "h1:oidDC4+YEuSIQbsR94rY9gur91UPL6DnxDCIYd2IGsE=",
        version = "v3.5.9",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v2",
        importpath = "go.etcd.io/etcd/client/v2",
        sum = "h1:AELPkjNR3/igjbO7CjyF1fPuVPjrblliiKj+Y6xSGOU=",
        version = "v2.305.7",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v3",
        importpath = "go.etcd.io/etcd/client/v3",
        sum = "h1:r5xghnU7CwbUxD/fbUtRyJGaYNfDun8sp/gTr1hew6E=",
        version = "v3.5.9",
    )
    go_repository(
        name = "io_gorm_driver_postgres",
        importpath = "gorm.io/driver/postgres",
        sum = "h1:ytTDxxEv+MplXOfFe3Lzm7SjG09fcdb3Z/c056DTBx0=",
        version = "v1.5.2",
    )
    go_repository(
        name = "io_gorm_driver_sqlite",
        importpath = "gorm.io/driver/sqlite",
        sum = "h1:zKYbzRCpBrT1bNijRnxLDJWPjVfImGEn0lSnUY5gZ+c=",
        version = "v1.5.0",
    )

    go_repository(
        name = "io_gorm_gorm",
        importpath = "gorm.io/gorm",
        sum = "h1:iyNd8fNAe8W9dvtlgeRI5zSVZPsq3OpcTu37cYcpCmw=",
        version = "v1.25.4",
    )
    go_repository(
        name = "io_gorm_plugin_opentelemetry",
        importpath = "gorm.io/plugin/opentelemetry",
        sum = "h1:7p0ocWELjSSRI7NCKPW2mVe6h43YPini99sNJcbsTuc=",
        version = "v0.1.4",
    )

    go_repository(
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:y73uSU6J157QMP2kn2r30vwW1A2W2WFwSCGnAVxeaD0=",
        version = "v0.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_github_com_gin_gonic_gin_otelgin",
        importpath = "go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin",
        sum = "h1:0KYeVr81ogcVRLXVcXFuPQMNZngplnP8MqrE8CqvHeg=",
        version = "v0.45.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_runtime",
        importpath = "go.opentelemetry.io/contrib/instrumentation/runtime",
        sum = "h1:EbmAUG9hEAMXyfWEasIt2kmh/WmXUznUksChApTgBGc=",
        version = "v0.42.0",
    )

    go_repository(
        name = "io_opentelemetry_go_contrib_propagators_b3",
        importpath = "go.opentelemetry.io/contrib/propagators/b3",
        sum = "h1:Yty9Vs4F3D6/liF1o6FNt0PvN85h/BJJ6DQKJ3nrcM0=",
        version = "v1.20.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_propagators_jaeger",
        importpath = "go.opentelemetry.io/contrib/propagators/jaeger",
        sum = "h1:Zbpbmwav32Ea5jSotpmkWEl3a6Xvd4tw/3xxGO1i05Y=",
        version = "v1.17.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_propagators_opencensus",
        importpath = "go.opentelemetry.io/contrib/propagators/opencensus",
        sum = "h1:wQBFvTXNs/UcCiWEi0cadOGJ8LMxDSOWbNtldFS0VI4=",
        version = "v0.42.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_propagators_ot",
        importpath = "go.opentelemetry.io/contrib/propagators/ot",
        sum = "h1:ufo2Vsz8l76eI47jFjuVyjyB3Ae2DmfiCV/o6Vc8ii0=",
        version = "v1.17.0",
    )

    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:MuS/TNf4/j4IXsZuJegVzI1cwut7Qc00344rgH7p8bs=",
        version = "v1.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_bridge_opencensus",
        importpath = "go.opentelemetry.io/otel/bridge/opencensus",
        sum = "h1:YHivttTaDhbZIHuPlg1sWsy2P5gj57vzqPfkHItgbwQ=",
        version = "v0.39.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_internal_retry",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/internal/retry",
        sum = "h1:t4ZwRPU+emrcvM2e9DHd0Fsf0JTPVcbfa/BhTDF03d0=",
        version = "v1.16.0",
    )

    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlpmetric",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlpmetric",
        sum = "h1:ZtfnDL+tUrs1F0Pzfwbg2d59Gru9NCH3bgSHBM6LDwU=",
        version = "v0.42.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlpmetric_otlpmetricgrpc",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
        sum = "h1:NmnYCiR0qNufkldjVvyQfZTHSdzeHoZ41zggMsdMcLM=",
        version = "v0.42.0",
    )

    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace",
        sum = "h1:Mne5On7VWdx7omSrSSZvM4Kw7cS7NQkOOmLcgscI51U=",
        version = "v1.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracegrpc",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
        sum = "h1:3d+S281UTjM+AbF31XSOYn1qXn3BgIdWl8HNEpx08Jk=",
        version = "v1.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_stdout_stdouttrace",
        importpath = "go.opentelemetry.io/otel/exporters/stdout/stdouttrace",
        sum = "h1:2PunuO5SbkN5MhCbuHCd3tC6qrcaj+uDAkX/qBU5BAs=",
        version = "v1.15.1",
    )

    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:aTzpGtV0ar9wlV4Sna9sdJyII5jTVJEvKETPiOKwvpE=",
        version = "v1.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:6USY6zH+L8uMH8L3t1enZPR3WFEmSTADlqldyHtJi3o=",
        version = "v1.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:EJoTO5qysMsYCa+w4UghwFV/ptQgqSL/8Ni+hx+8i1k=",
        version = "v1.19.0",
    )

    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:DFVQmlVbfVeOuBRrwdtaehRrWiL1JoVs9CPIQ1Dzxpg=",
        version = "v1.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_proto_otlp",
        importpath = "go.opentelemetry.io/proto/otlp",
        sum = "h1:T0TX0tmXU8a3CbNXzEKGeU5mIVOdf0oykP+u2lIVU/I=",
        version = "v1.0.0",
    )

    go_repository(
        name = "io_rsc_binaryregexp",
        importpath = "rsc.io/binaryregexp",
        sum = "h1:HfqmD5MEmC0zvwBuF187nq9mdnXjXsSivRiXN7SmRkE=",
        version = "v0.2.0",
    )

    go_repository(
        name = "io_rsc_pdf",
        importpath = "rsc.io/pdf",
        sum = "h1:k1MczvYDUvJBe93bYd7wrZLLUEcLZAuF824/I4e5Xr4=",
        version = "v0.1.1",
    )
    go_repository(
        name = "io_rsc_quote_v3",
        importpath = "rsc.io/quote/v3",
        sum = "h1:9JKUTTIUgS6kzR9mK1YuGKv6Nl+DijDNIc0ghT58FaY=",
        version = "v3.1.0",
    )
    go_repository(
        name = "io_rsc_sampler",
        importpath = "rsc.io/sampler",
        sum = "h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=",
        version = "v1.3.0",
    )
    go_repository(
        name = "org_golang_google_api",
        importpath = "google.golang.org/api",
        sum = "h1:zDobeejm3E7pEG1mNHvdxvjs5XJoCMzyNH+CmwL94Es=",
        version = "v0.122.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
        version = "v1.6.7",
    )
    go_repository(
        name = "org_golang_google_genproto",
        importpath = "google.golang.org/genproto",
        sum = "h1:Z0hjGZePRE0ZBWotvtrwxFNrNE9CUAGtplaDK5NNI/g=",
        version = "v0.0.0-20230711160842-782d3b101e98",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:FmF5cCW94Ij59cfpoLiwTgodWmm60eEV0CjlsVg2fuw=",
        version = "v0.0.0-20230711160842-782d3b101e98",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:bVf09lpb+OJbByTj913DRJioFFAjf/ZGxEz7MajTp2U=",
        version = "v0.0.0-20230711160842-782d3b101e98",
    )

    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:SXUpjxeVF3FKrTYQI4f4KvbGD5u2xccdYdurwowix5I=",
        version = "v1.58.2",
    )

    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:g0LDEJHgrBl9N9r17Ru3sqWhkIx2NB67okBHPwC7hs8=",
        version = "v1.31.0",
    )
    go_repository(
        name = "org_golang_x_arch",
        importpath = "golang.org/x/arch",
        sum = "h1:02VY4/ZcO/gBOH6PUaoiptASxtXU10jazRCP865E97k=",
        version = "v0.3.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:mvySKfSWJ+UKUii46M40LOvyWfN0s2U+46/jDd0e6Ck=",
        version = "v0.13.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:QE6XYQK6naiK1EPAe1g/ILLxN5RBoH5xkJk3CqlMI/Y=",
        version = "v0.0.0-20200224162631-6cc2880d07d6",
    )
    go_repository(
        name = "org_golang_x_image",
        importpath = "golang.org/x/image",
        sum = "h1:+qEpEAPhDZ1o0x3tHzZTQDArnOixOzGD9HUJfcg0mb4=",
        version = "v0.0.0-20190802002840-cff245a6509b",
    )

    go_repository(
        name = "org_golang_x_lint",
        importpath = "golang.org/x/lint",
        sum = "h1:2M3HP5CCK1Si9FQhwnzYhXdG6DXeebvUHFpre8QvbyI=",
        version = "v0.0.0-20201208152925-83fdc39ff7b5",
    )
    go_repository(
        name = "org_golang_x_mobile",
        importpath = "golang.org/x/mobile",
        sum = "h1:4+4C/Iv2U4fMZBiMCc98MG1In4gJY5YRhtpDNeDeHWs=",
        version = "v0.0.0-20190719004257-d2bd2a29d028",
    )

    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:LUYupSeNrTNCGzR/hVBk2NHZO4hXcVaW1k4Qx7rjPx8=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:ugBLEUaxABaB5AJqW9enI0ACdci2RUd4eP51NTBvuJ8=",
        version = "v0.15.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:zHCpF2Khkwy4mMB4bv0U37YtJdTGW8jI0glAApi0Kh8=",
        version = "v0.10.0",
    )

    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:ftCYgMx6zT/asHUrPw8BLLscYtGznsLAnjq5RH9P66E=",
        version = "v0.3.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:CM0HF96J0hcLAwsHPJZjfdNzs0gftsLfgKt57wWHJ0o=",
        version = "v0.12.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:/ZfYdc3zq+q02Rv9vGqTeSItdzZTSNDmfTi0mBAuidU=",
        version = "v0.12.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:ablQoSUd0tRdKxZewP80B+BaqeKJuVhuRxj/dkrun3k=",
        version = "v0.13.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:xYY+Bajn2a7VBmTM5GikTmnK8ZuX8YgnQCqZpbBNtmA=",
        version = "v0.1.0",
    )

    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:BOw41kyTf3PuCW1pVQf8+Cyg8pMlkYB1oo9iJ6D/lKM=",
        version = "v0.6.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:H2TDz8ibqkAF6YGhCdN3jS9O0/s90v0rJh3X/OLHEUk=",
        version = "v0.0.0-20220907171357-04be3eba64a2",
    )
    go_repository(
        name = "org_uber_go_atomic",
        importpath = "go.uber.org/atomic",
        sum = "h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=",
        version = "v1.9.0",
    )
    go_repository(
        name = "org_uber_go_dig",
        importpath = "go.uber.org/dig",
        sum = "h1:5Chju+tUvcC+N7N6EV08BJz41UZuO3BmHcN4A287ZLI=",
        version = "v1.17.0",
    )
    go_repository(
        name = "org_uber_go_fx",
        importpath = "go.uber.org/fx",
        sum = "h1:ZMC/pnRvhsthOZh9MZjMq5U8Or3mA9zBSPaLnzs3ihQ=",
        version = "v1.20.0",
    )
    go_repository(
        name = "org_uber_go_goleak",
        importpath = "go.uber.org/goleak",
        sum = "h1:NBol2c7O1ZokfZ0LEU9K6Whx/KnwvepVetCUhtKja4A=",
        version = "v1.2.1",
    )
    go_repository(
        name = "org_uber_go_multierr",
        importpath = "go.uber.org/multierr",
        sum = "h1:blXXJkSxSSfBVBlC76pxqeO+LN3aDfLQo+309xJstO0=",
        version = "v1.11.0",
    )
    go_repository(
        name = "org_uber_go_zap",
        importpath = "go.uber.org/zap",
        sum = "h1:FiJd5l1UOLj0wCgbSE0rwwXHzEdAZS6hiiSnxJN/D60=",
        version = "v1.24.0",
    )
