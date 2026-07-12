<template>
  <!-- KubeAdmin logo：蓝色圆角六边形 + 白色舵轮（Kubernetes helm wheel）。
       与 docs-site/public/logo.svg 同源；class="icon" 让 el-icon 用 font-size 控制尺寸。 -->
  <svg class="icon" viewBox="0 0 128 128" xmlns="http://www.w3.org/2000/svg">
    <polygon
      points="64,8 112.5,36 112.5,92 64,120 15.5,92 15.5,36"
      fill="#2563EB"
      stroke="#2563EB"
      stroke-width="11"
      stroke-linejoin="round"
    />
    <g stroke="#FFFFFF" stroke-width="5" stroke-linecap="round" fill="none">
      <circle cx="64" cy="64" r="34" />
      <line x1="78" y1="64" x2="94" y2="64" />
      <line x1="73.9" y1="73.9" x2="85.2" y2="85.2" />
      <line x1="64" y1="78" x2="64" y2="94" />
      <line x1="54.1" y1="73.9" x2="42.8" y2="85.2" />
      <line x1="50" y1="64" x2="34" y2="64" />
      <line x1="54.1" y1="54.1" x2="42.8" y2="42.8" />
      <line x1="64" y1="50" x2="64" y2="34" />
      <line x1="73.9" y1="54.1" x2="85.2" y2="42.8" />
    </g>
    <g fill="#FFFFFF">
      <circle cx="64" cy="64" r="11" />
      <circle cx="98" cy="64" r="3.5" />
      <circle cx="88.04" cy="88.04" r="3.5" />
      <circle cx="64" cy="98" r="3.5" />
      <circle cx="39.96" cy="88.04" r="3.5" />
      <circle cx="30" cy="64" r="3.5" />
      <circle cx="39.96" cy="39.96" r="3.5" />
      <circle cx="64" cy="30" r="3.5" />
      <circle cx="88.04" cy="39.96" r="3.5" />
    </g>
  </svg>
</template>
