<#import "template.ftl" as layout>
<@layout.registrationLayout displayMessage=!messagesPerField.existsError('username','password') displayInfo=realm.password && realm.registrationAllowed && !registrationDisabled??; section>
    <#if section = "header">
        <h2 class="h3 text-dark">${msg("loginAccountTitle")}</h2>
        <p class="text-secondary">${msg("loginAccountSubtitle")}</p>
    </#if>

    <#if section = "rightContent">
        <div class="d-flex">
            <div class="w-100 me-5">
                <div class="col-12 d-flex justify-content-end">
                    <img src="${url.resourcesPath}/img/login.svg" alt="Login">
                </div>

                <div class="col-12 text-end">
                    <h1 class="display-2">${msg("welcomeText")}</h1>
                </div>
            </div>
        </div>
    </#if>

    <#if section = "leftTopContent">
        <#if realm.password && realm.registrationAllowed && !registrationDisabled??>
            <div class="col-12 pt-4 pe-5 d-flex justify-content-end mb-5">
                <span>
                    <span class="me-2">${msg("noAccount")}</span>
                    <a class="btn btn-primary" role="button" href="${url.registrationUrl}">${msg("doRegister")}</a>
                </span>
            </div>
        </#if>
    </#if>

    <#if section = "form">
        <#if realm.password>
            <form onsubmit="login.disabled = true; return true;" action="${url.loginAction}" method="post">
                <#if !usernameHidden??>
                    <div class="mb-2">
                        <label class="form-label" for="username"><#if !realm.loginWithEmailAllowed>${msg("username")}<#elseif !realm.registrationEmailAsUsername>${msg("usernameOrEmail")}<#else>${msg("email")}</#if></label>
                        <input tabindex="1" type="text" id="username" name="username"
                               class="form-control shadow <#if messagesPerField.existsError('username','password')>is-invalid</#if>"
                               value="${(login.username!'')}" />

                        <#if messagesPerField.existsError('username','password')>
                            <div class="invalid-feedback">
                                ${kcSanitize(messagesPerField.getFirstError('username','password'))?no_esc}
                            </div>
                        </#if>
                    </div>
                </#if>

                <div class="mb-4">
                    <div class="d-flex justify-content-between">
                        <label class="form-label" for="password">${msg("password")}</label>
                        <#if realm.resetPasswordAllowed>
                            <span><a tabindex="3" href="${url.loginResetCredentialsUrl}">${msg("doForgotPassword")}</a></span>
                        </#if>
                    </div>
                    <input tabindex="2" type="password" id="password" name="password"
                           class="form-control shadow <#if messagesPerField.existsError('username','password')>is-invalid</#if>"
                    />

                    <#if usernameHidden?? && messagesPerField.existsError('username','password')>
                        <div class="invalid-feedback">
                            ${kcSanitize(messagesPerField.getFirstError('username','password'))?no_esc}
                        </div>
                    </#if>
                </div>

                <#if realm.rememberMe && !usernameHidden??>
                    <div class="form-check">
                        <input tabindex="4" id="rememberMe" name="rememberMe" type="checkbox" class="form-check-input" <#if login.rememberMe??>checked</#if>>
                        <label class="form-check-label" for="rememberMe">${msg("rememberMe")} </label>
                    </div>
                </#if>

                <div class="mt-4 d-grid">
                    <input type="hidden" id="id-hidden-input" name="credentialId" <#if auth.selectedCredential?has_content>value="${auth.selectedCredential}"</#if>/>
                    <input tabindex="5" class="btn btn-primary btn-block" name="login" type="submit" value="${msg("doLogIn")}"/>
                </div>
            </form>
        </#if>
    </#if>

    <#if section = "socialProviders">
        <#if realm.password && social.providers??>
            <div id="kc-social-providers">
                <div class="divider d-flex align-items-center my-4">
                    <p class="text-center fw-bold mx-3 mb-0 text-muted">${msg("identity-provider-login-label")}</p>
                </div>

                <div class="d-grid gap-2">
                    <#list social.providers as p>
                        <a class="btn btn-outline-secondary" role="button" href="${p.loginUrl}">
                            <i class="fa-brands fa-${p.alias} fa-lg"></i> ${msg("signInWith", p.displayName)}
                        </a>
                    </#list>
                </div>
            </div>
        </#if>
    </#if>
</@layout.registrationLayout>
