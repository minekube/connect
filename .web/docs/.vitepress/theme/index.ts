import Theme from 'vitepress/theme'
import './styles/vars.css'

export default {
  ...Theme,
  enhanceApp({ app }) {
    // app.component('TextAndImage', SvgImage)
  }
}
