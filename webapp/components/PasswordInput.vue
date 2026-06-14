<template>
  <div class="password-input">
    <input
      :value="modelValue"
      type="password"
      class="input"
      :required="required"
      autocomplete="new-password"
      @input="onInput"
    />
    <div class="password-strength">
      <span
        v-for="n in 4"
        :key="n"
        class="password-strength-bar"
        :class="barClass(n)"
      ></span>
    </div>
  </div>
</template>

<script>
// Replaces vue-password (Vue 2 only). Keeps the same interface the auth pages
// rely on: v-model for the value, a `strength` prop (0-4) driving the meter,
// and an `input` event carrying the raw value.
export default {
  name: 'PasswordInput',
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    strength: {
      type: Number,
      default: 0,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue', 'input'],
  methods: {
    onInput(event) {
      const value = event.target.value
      this.$emit('update:modelValue', value)
      this.$emit('input', value)
    },
    barClass(n) {
      if (n > this.strength) {
        return ''
      }
      if (this.strength <= 1) {
        return 'is-danger'
      }
      if (this.strength === 2) {
        return 'is-warning'
      }
      if (this.strength === 3) {
        return 'is-info'
      }
      return 'is-success'
    },
  },
}
</script>

<style scoped>
.password-strength {
  display: flex;
  gap: 4px;
  margin-top: 6px;
}
.password-strength-bar {
  height: 4px;
  flex: 1;
  border-radius: 2px;
  background: #dbdbdb;
}
.password-strength-bar.is-danger {
  background: #f14668;
}
.password-strength-bar.is-warning {
  background: #ffe08a;
}
.password-strength-bar.is-info {
  background: #3e8ed0;
}
.password-strength-bar.is-success {
  background: #48c78e;
}
</style>
