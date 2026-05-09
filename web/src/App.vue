<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <RouterView />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router'
import type { GlobalThemeOverrides } from 'naive-ui'

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#8b5cf6',
    primaryColorHover: '#a78bfa',
    primaryColorPressed: '#6d28d9',
    primaryColorSuppl: '#8b5cf6',
    bodyColor: '#0b0716',
    cardColor: '#1a1130',
    modalColor: '#1a1130',
    popoverColor: '#1a1130',
    textColorBase: '#eae3ff',
    textColor1: '#eae3ff',
    textColor2: '#b8add8',
    textColor3: '#7b6ea0',
    borderColor: 'rgba(196,181,253,0.12)',
    dividerColor: 'rgba(196,181,253,0.12)',
    inputColor: 'rgba(26,17,48,0.72)',
    tableColor: '#120b22',
    tableHeaderColor: '#17102a',
    codeColor: '#0e0820'
  },
  Button: {
    textColorPrimary: '#0b0716',
    textColorHoverPrimary: '#0b0716'
  },
  Input: {
    colorFocus: 'rgba(26,17,48,0.85)'
  },
  Card: {
    color: 'linear-gradient(160deg, rgba(26,17,48,0.82), rgba(18,11,34,0.95))'
  }
}
</script>
