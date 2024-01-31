<#macro registrationLayout bodyClass="" displayInfo=false displayMessage=true displayRequiredFields=false>
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="robots" content="noindex, nofollow">

    <#if properties.meta?has_content>
        <#list properties.meta?split(' ') as meta>
            <meta name="${meta?split('==')[0]}" content="${meta?split('==')[1]}"/>
        </#list>
    </#if>
    <title>${msg("loginTitle",(realm.displayName!''))}</title>
    <link rel="icon" href="${url.resourcesPath}/img/favicon.ico" />
    <#if properties.stylesCommon?has_content>
        <#list properties.stylesCommon?split(' ') as style>
            <link href="${url.resourcesCommonPath}/${style}" rel="stylesheet" />
        </#list>
    </#if>
    <#if properties.styles?has_content>
        <#list properties.styles?split(' ') as style>
            <link href="${url.resourcesPath}/${style}" rel="stylesheet" />
        </#list>
    </#if>
    <#if properties.scripts?has_content>
        <#list properties.scripts?split(' ') as script>
            <script src="${url.resourcesPath}/${script}" type="text/javascript"></script>
        </#list>
    </#if>
    <#if scripts??>
        <#list scripts as script>
            <script src="${script}" type="text/javascript"></script>
        </#list>
    </#if>
</head>

<body>
<div class="container-fluid h-100">
    <div class="row h-100">
        <div class="col-6 d-none d-md-block login-right">
            <div class="row h-100">
                <div class="col-12 pt-4 ps-4">
                    <img src="${url.resourcesPath}/img/dupman.png" class="logo"
                         alt="${kcSanitize(msg("loginTitleHtml",(realm.displayNameHtml!'')))?no_esc}">
                </div>

                <#nested "rightContent">
            </div>
        </div>
        <div class="col-12 col-md-6">
            <div class="row h-100">
                <div class="col-12 pt-4 ps-4 d-block d-md-none">
                    <div class="d-flex justify-content-center">
                        <img src="${url.resourcesPath}/img/dupman.png" class="logo"
                             alt="${kcSanitize(msg("loginTitleHtml",(realm.displayNameHtml!'')))?no_esc}">
                    </div>
                </div>

                <div class="col-12 pt-4 pe-5 d-flex justify-content-end mb-5">
                    <#nested "leftTopContent">
                </div>

                <div class="d-flex justify-content-center">
                    <div class="col-12 w-75 w-xl-50">
                        <#if !(auth?has_content && auth.showUsername() && !auth.showResetCredentials())>
                            <#nested "header">
                        <#else>
                            <#nested "show-username">
                            <div id="kc-username" class="${properties.kcFormGroupClass!}">
                                <label id="kc-attempted-username">${auth.attemptedUsername}</label>
                                <a id="reset-login" href="${url.loginRestartFlowUrl}" aria-label="${msg("restartLoginTooltip")}">
                                    <div class="kc-login-tooltip">
                                        <i class="${properties.kcResetFlowIcon!}"></i>
                                        <span class="kc-tooltip-text">${msg("restartLoginTooltip")}</span>
                                    </div>
                                </a>
                            </div>
                        </#if>

                        <#if displayMessage && message?has_content && (message.type != 'warning' || !isAppInitiatedAction??)>
                            <div class="alert alert-<#if message.type = 'error'>danger<#else>${message.type}</#if>" role="alert">
                                ${kcSanitize(message.summary)?no_esc}
                            </div>
                        </#if>

                        <#nested "form">

                        <#if auth?has_content && auth.showTryAnotherWayLink()>
                            <form id="kc-select-try-another-way-form" action="${url.loginAction}" method="post">
                                <div class="${properties.kcFormGroupClass!}">
                                    <input type="hidden" name="tryAnotherWay" value="on"/>
                                    <a href="#" id="try-another-way"
                                       onclick="document.forms['kc-select-try-another-way-form'].submit();return false;">${msg("doTryAnotherWay")}</a>
                                </div>
                            </form>
                        </#if>

                        <#nested "socialProviders">

                        <#if displayInfo>
                            <#nested "info">
                        </#if>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>
</#macro>
