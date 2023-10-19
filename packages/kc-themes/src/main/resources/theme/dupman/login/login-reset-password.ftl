<#import "template.ftl" as layout>
<@layout.registrationLayout displayInfo=true displayMessage=!messagesPerField.existsError('username'); section>
    <#if section = "header">
        <h2 class="h3 text-dark">${msg("emailForgotTitle")}</h2>
        <p class="text-secondary">${msg("emailForgotSubtitle")}</p>
    </#if>

    <#if section = "rightContent">
        <div class="d-flex">
            <div class="w-100 me-5">
                <div class="col-12 d-flex justify-content-end">
                    <img src="${url.resourcesPath}/img/forgot-password.svg" class="w-100" alt="Forgot Password">
                </div>
            </div>
        </div>
    </#if>

    <#if section = "leftTopContent">
        <#if realm.password && realm.registrationAllowed && !registrationDisabled??>
            <div class="col-12 pt-4 pe-5 d-flex justify-content-end mb-5">
                <span>
                    <span class="me-2">${msg("alreadyHaveAccount")}</span>
                    <a class="btn btn-primary" role="button" href="${url.loginUrl}">${msg("doLogIn")}</a>
                </span>
            </div>
        </#if>
    </#if>

    <#if section = "form">
        <form action="${url.loginAction}" method="post">
            <div class="mb-4">
                <label class="form-label" for="username"><#if !realm.loginWithEmailAllowed>${msg("username")}<#elseif !realm.registrationEmailAsUsername>${msg("usernameOrEmail")}<#else>${msg("email")}</#if></label>
                <input type="text" id="username" name="username" autofocus
                       class="form-control shadow <#if messagesPerField.existsError('username')>is-invalid</#if>"
                       value="${(auth.attemptedUsername!'')}" />

                <#if messagesPerField.existsError('username')>
                    <div class="invalid-feedback">
                        ${kcSanitize(messagesPerField.get('username'))?no_esc}
                    </div>
                </#if>
            </div>

            <div class="mt-4 d-grid">
                <input class="btn btn-primary btn-block" type="submit" value="${msg("doSubmit")}"/>
            </div>
        </form>
    </#if>

    <#if section = "info">
        <p class="mt-4">
            <#if realm.duplicateEmailsAllowed>
                ${msg("emailInstructionUsername")}
            <#else>
                ${msg("emailInstruction")}
            </#if>
        </p>
    </#if>
</@layout.registrationLayout>
