<#import "template.ftl" as layout>
<@layout.registrationLayout; section>
    <#if section = "header">
        <h2 class="h3 text-dark">${msg("logoutConfirmTitle")}</h2>
        <p class="text-secondary">${msg("logoutConfirmHeader")}</p>
    </#if>

    <#if section = "form">
        <form action="${url.logoutConfirmAction}" method="POST">
            <input type="hidden" name="session_code" value="${logoutConfirm.code}">
            <div class="mt-4 d-grid">
                <input tabindex="4" class="btn btn-primary btn-block" name="confirmLogout" id="kc-logout" type="submit" value="${msg("doLogout")}"/>
            </div>
        </form>

        <div class="mt-4">
            <#if logoutConfirm.skipLink>
            <#else>
                <#if (client.baseUrl)?has_content>
                    <a href="${client.baseUrl}">${kcSanitize(msg("backToApplication"))?no_esc}</a>
                </#if>
            </#if>
        </div>
    </#if>
</@layout.registrationLayout>
