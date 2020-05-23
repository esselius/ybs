<template>
    <a :href="auth_url">login</a>
</template>

<script lang="ts">
    import Vue from "vue";

    export default Vue.extend({
        data: function() {
            let client_id = '107d8176c8d51180d2304f193f006888676cdb3eb6d77febe209ba20d35087d6'
            let redirect_uri = 'http://localhost:1234'

            return {
                auth_url: "https://app.youneedabudget.com/oauth/authorize?client_id=" + client_id + "&redirect_uri="+ redirect_uri + "&response_type=token"
            }
        },
    });
</script>


