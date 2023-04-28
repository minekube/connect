import Theme from 'vitepress/theme'
import './styles/vars.css'
import VPBadge from 'vitepress/dist/client/theme-default/components/VPBadge.vue'

export default {
  ...Theme,
  enhanceApp({ app }) {
    // app.component('TextAndImage', SvgImage)
    app.component('VPBadge', VPBadge)
  }
}
