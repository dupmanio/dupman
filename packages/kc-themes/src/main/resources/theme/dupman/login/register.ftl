<#import "template.ftl" as layout>
<@layout.registrationLayout displayMessage=!messagesPerField.existsError('firstName','lastName','email','username','password','password-confirm'); section>
    <#if section = "header">
        <h2 class="h3 text-dark">${msg("registerTitle")}</h2>
        <p class="text-secondary">${msg("registerSubtitle")}</p>
    </#if>

    <#if section = "rightContent">
        <div class="d-flex">
            <div class="w-100 me-5">
                <div class="col-12 d-flex justify-content-end">
                    <img src="${url.resourcesPath}/img/sign-in.svg" class="w-100" alt="Regisster">
                </div>
            </div>
        </div>
    </#if>

    <#if section = "leftTopContent">
        <#if realm.password && realm.registrationAllowed && !registrationDisabled??>
            <span>
                <span class="me-2">${msg("alreadyHaveAccount")}</span>
                <a class="btn btn-primary" role="button" href="${url.loginUrl}">${msg("doLogIn")}</a>
            </span>
        </#if>
    </#if>

    <#if section = "form">
        <form action="${url.registrationAction}" method="post">
            <div class="mb-2">
                <label class="form-label" for="firstName">${msg("firstName")}</label>
                <input type="text" id="firstName" name="firstName"
                       class="form-control shadow <#if messagesPerField.existsError('firstName')>is-invalid</#if>"
                       value="${(register.formData.firstName!'')}" />

                <#if messagesPerField.existsError('firstName')>
                    <div class="invalid-feedback">
                        ${kcSanitize(messagesPerField.get('firstName'))?no_esc}
                    </div>
                </#if>
            </div>

            <div class="mb-2">
                <label class="form-label" for="lastName">${msg("lastName")}</label>
                <input type="text" id="lastName" name="lastName"
                       class="form-control shadow <#if messagesPerField.existsError('lastName')>is-invalid</#if>"
                       value="${(register.formData.lastName!'')}" />

                <#if messagesPerField.existsError('lastName')>
                    <div class="invalid-feedback">
                        ${kcSanitize(messagesPerField.get('lastName'))?no_esc}
                    </div>
                </#if>
            </div>

            <div class="mb-2">
                <label class="form-label" for="email">${msg("email")}</label>
                <input type="text" id="email" name="email"
                       class="form-control shadow <#if messagesPerField.existsError('email')>is-invalid</#if>"
                       value="${(register.formData.email!'')}"
                       autocomplete="email" />

                <#if messagesPerField.existsError('email')>
                    <div class="invalid-feedback">
                        ${kcSanitize(messagesPerField.get('email'))?no_esc}
                    </div>
                </#if>
            </div>

            <#if !realm.registrationEmailAsUsername>
                <div class="mb-2">
                    <label class="form-label" for="username">${msg("username")}</label>
                    <input type="text" id="username" name="username"
                           class="form-control shadow <#if messagesPerField.existsError('username')>is-invalid</#if>"
                           value="${(register.formData.username!'')}"
                           autocomplete="username" />

                    <#if messagesPerField.existsError('username')>
                        <div class="invalid-feedback">
                            ${kcSanitize(messagesPerField.get('username'))?no_esc}
                        </div>
                    </#if>
                </div>
            </#if>

            <#if passwordRequired??>
                <div class="mb-2">
                    <label class="form-label" for="password">${msg("password")}</label>
                    <input type="password" id="password" name="password"
                           class="form-control shadow <#if messagesPerField.existsError('password', 'password-confirm')>is-invalid</#if>"
                           value="${(register.formData.password!'')}"
                           autocomplete="new-password" />

                    <#if messagesPerField.existsError('password')>
                        <div class="invalid-feedback">
                            ${kcSanitize(messagesPerField.get('password'))?no_esc}
                        </div>
                    </#if>
                </div>

                <div class="mb-2">
                    <label class="form-label" for="password-confirm">${msg("passwordConfirm")}</label>
                    <input type="password" id="password-confirm" name="password-confirm"
                           class="form-control shadow <#if messagesPerField.existsError('password-confirm')>is-invalid</#if>"
                           value="${(register.formData.password!'')}"
                           autocomplete="new-password" />

                    <#if messagesPerField.existsError('password-confirm')>
                        <div class="invalid-feedback">
                            ${kcSanitize(messagesPerField.get('password-confirm'))?no_esc}
                        </div>
                    </#if>
                </div>
            </#if>

            <#if recaptchaRequired??>
                <div class="form-group">
                    <div class="${properties.kcInputWrapperClass!}">
                        <div class="g-recaptcha" data-size="compact" data-sitekey="${recaptchaSiteKey}"></div>
                    </div>
                </div>
            </#if>

            <div class="mt-4 d-grid">
                <input class="btn btn-primary btn-block" type="submit" value="${msg("doRegister")}"/>
            </div>
        </form>
    </#if>
</@layout.registrationLayout>
