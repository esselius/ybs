<template>
  <a :href="auth_url">Perform login</a>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  data() {
    const clientId =
      '107d8176c8d51180d2304f193f006888676cdb3eb6d77febe209ba20d35087d6'
    const redirectUri = 'http://localhost:1234'
    return {
      auth_url: `https://app.youneedabudget.com/oauth/authorize?clientId=${clientId}&redirect_uri=${redirectUri}&response_type=token`
    }
  }
})
</script>
