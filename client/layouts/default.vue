<template>

  <div class="
      antialiased
    ">

    <header class="text-gray-700 body-font border-b border-gray-200">
      <div class="container mx-auto flex flex-wrap p-5 flex-col md:flex-row items-center">
        <a
          class="flex title-font font-medium items-center text-gray-900 mb-4 md:mb-0"
          target="_blank"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            stroke="currentColor"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            class="w-10 h-10 text-white p-2 bg-indigo-500 rounded-full"
            viewBox="0 0 24 24"
          >
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"></path>
          </svg>
          <span class="ml-3 text-xl">HermesCI</span>
        </a>
        <nav class="md:ml-auto flex flex-wrap items-center text-base justify-center">
          <NuxtLink
            class="mr-5 hover:text-gray-900"
            to="/"
          >
            Home
          </NuxtLink>

          <NuxtLink
            class="mr-5 hover:text-gray-900"
            to="/#features"
          >
            Features
          </NuxtLink>
          <NuxtLink
            class="mr-5 hover:text-gray-900"
            to="/#docs"
          >
            Docs
          </NuxtLink>
          <NuxtLink
            class="mr-5 hover:text-gray-900"
            to="/#about"
          >
            About
          </NuxtLink>

          <template v-if="!isAuthenticated">

            <nuxt-link
              class="nav-link active"
              aria-current="page"
              to="/login"
            >
              <a class="mr-5 hover:text-gray-900">Login</a>

            </nuxt-link>

          </template>
          <template v-else>
            <NuxtLink
              class="mr-5 hover:text-gray-900"
              to="/repos"
            >

              <button class="text-white bg-indigo-500 border-0 py-1 px-4 rounded focus:outline-none hover:bg-indigo-600">
                Projects
              </button>

            </NuxtLink>

            <nuxt-link
              class="nav-link active"
              aria-current="page"
              to="#"
            >
              <a
                class="mr-5 hover:text-gray-900"
                @click="logout"
              >Logout
                <span>
                  <a
                    class="mr-5 hover:text-gray-900 text-sm"
                    href="#"
                  >
                    (@{{ user.login }})
                  </a>
                </span>
              </a>

            </nuxt-link>

          </template>
        </nav>

      </div>
    </header>

    <div class="">

      <Nuxt />

    </div>

    <footer class="text-gray-700 body-font">

      <div class="border-t border-gray-200">
        <div class="container px-5 py-8 flex flex-wrap mx-auto items-center">
          <div class="flex md:flex-no-wrap flex-wrap justify-center md:justify-start">
            <p class="text-gray-500 text-sm text-center sm:text-left">© 2022 —
              <a
                href="#"
                class="text-gray-600 ml-1"
                target="_blank"
                rel="noopener noreferrer"
              >Hermes CI</a>
            </p>
            <span class="sm:ml-auto sm:mt-0 mt-2 ml-2 sm:w-auto w-full  text-center text-gray-500 text-sm">
              <a href="#">
                Github
              </a>
            </span>
          </div>
          <span class="inline-flex lg:ml-auto lg:mt-0 mt-6 w-full justify-center md:justify-start md:w-auto">
            <a
              href="#"
              class="ml-3 text-gray-500"
            >
              <svg
                fill="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                class="w-5 h-5"
                viewBox="0 0 28 28"
              >
                <path d="M18,3C8.373,3,3,8.373,3,15c0,5.623,3.872,10.328,9.092,11.63C12.036,26.468,12,26.28,12,26.047v-2.051 c-0.487,0-1.303,0-1.508,0c-0.821,0-1.551-0.353-1.905-1.009c-0.393-0.729-0.461-1.844-1.435-2.526 c-0.289-0.227-0.069-0.486,0.264-0.451c0.615,0.174,1.125,0.596,1.605,1.222c0.478,0.627,0.703,0.769,1.596,0.769 c0.433,0,1.081-0.025,1.691-0.121c0.328-0.833,0.895-1.6,1.588-1.962c-3.996-0.411-5.903-2.399-5.903-5.098 c0-1.162,0.495-2.286,1.336-3.233C9.053,10.647,8.706,8.73,9.435,8c1.798,0,2.885,1.166,3.146,1.481C13.477,9.174,14.461,9,15.495,9 c1.036,0,2.024,0.174,2.922,0.483C18.675,9.17,19.763,8,21.565,8c0.732,0.731,0.381,2.656,0.102,3.594 c0.836,0.945,1.328,2.066,1.328,3.226c0,2.697-1.904,4.684-5.894,5.097C18.199,20.49,19,22.1,19,23.313v2.734 c0,0.104-0.023,0.179-0.035,0.268C23.641,24.676,27,20.236,27,15C27,8.373,21.627,3,15,3z"></path>
              </svg>
            </a>

            <a class="ml-3 text-gray-500">
              <svg
                fill="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                class="w-5 h-5"
                viewBox="0 0 24 24"
              >
                <path d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z"></path>
              </svg>
            </a>

            <a class="ml-3 text-gray-500">
              <svg
                fill="currentColor"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="0"
                class="w-5 h-5"
                viewBox="0 0 24 24"
              >
                <path
                  stroke="none"
                  d="M16 8a6 6 0 016 6v7h-4v-7a2 2 0 00-2-2 2 2 0 00-2 2v7h-4v-7a6 6 0 016-6zM2 9h4v12H2z"
                ></path>
                <circle
                  cx="4"
                  cy="4"
                  r="2"
                  stroke="none"
                ></circle>
              </svg>
            </a>
          </span>
        </div>
      </div>

    </footer>

  </div>

</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
  computed: {
    isAuthenticated() {
      return this.$store.state.auth.loggedIn // it check if user isAuthenticated
    },
    user() {
      return this.$auth.user // it get the user from the store
    },
  },
  methods: {
    async logout() {
      await this.$auth.logout() // this method will logout the user and make token to false on the local storage of the user browser
    },
  },
})
</script>