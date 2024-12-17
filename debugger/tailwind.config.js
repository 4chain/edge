/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        echogy: {
          text: {
            primary: '#0a0a0a',     // 亮色主题 - 深色文本
            secondary: '#666666',   // 亮色主题 - 次要文本
            'primary-dark': '#ffffff',     // 暗色主题 - 白色文本
            'secondary-dark': '#888888',   // 暗色主题 - 次要文本
          },
          bg: {
            primary: '#ffffff',     // 亮色主题 - 白色背景
            secondary: '#f5f5f5',   // 亮色主题 - 浅灰背景
            'primary-dark': '#0a0a0a',     // 暗色主题 - 接近纯黑
            'secondary-dark': '#1a1a1a',   // 暗色主题 - 深灰
          },
          border: {
            DEFAULT: '#e5e5e5',     // 亮色主题 - 浅灰边框
            dark: '#242424',        // 暗色主题 - 深色边框
          },
          hover: {
            primary: '#f5f5f5',     // 亮色主题 - 主要悬停
            secondary: '#eeeeee',   // 亮色主题 - 次要悬停
            'primary-dark': '#242424',    // 暗色主题 - 主要悬停
            'secondary-dark': '#2a2a2a',  // 暗色主题 - 次要悬停
          },
          green: '#10b981',         // 成功状态
          red: '#ef4444',          // 错误状态
          gray: '#6b7280',         // 中性状态
        },
      },
    },
  },
  plugins: [],
}
