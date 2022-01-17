import { Middleware } from '@nuxt/types'

const auth: Middleware = ({ store, redirect, route }) => {
  if (!store.state.auth.loggedIn) {
    console.log('not logged in')
    const path = encodeURIComponent(route.path)
    return redirect(`/login?r=${path}`)
  }
  // Use context
}

export default auth
