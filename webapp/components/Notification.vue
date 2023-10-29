<template>
<article class="message" :class=level>
  <div class="message-header">
    <p>{{ title }}</p>
    <button class="delete" aria-label="delete" @click="close"></button>
  </div>
  <div class="message-body">
     {{ message }}
  </div>
</article>
</template>

<script>
  export default {
    name: 'Notification',
    props: {
      message: {
        type: String,
        default: ''
      },
      level: {
        type: String,
        default: 'is-danger'
      }
    },
    computed: {
      title: {
        get () {
          if (this.level === "is-danger") {
            return "Error"
          }
          if (this.level === "is-info") {
            return "Information"
          }
          if (this.level === "is-warn") {
            return "Warning"
          }
          return "Message"
        },
      },
    },
    methods: {
      close () {
        this.$emit('close')
      },
    },
  };
</script>
