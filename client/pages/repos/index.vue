<template>
  <div>
    <div class="  relative rounded overflow-hidden bg-white">
      <div class="py-8">
        <div class="bg-white mx-auto max-w-2xl	  shadow-lg rounded-lg overflow-hidden">
          <div class="sm:flex sm:items-center px-2 py-4">
            <div class="flex-grow">
              <h3 class="font-normal px-2 py-3 leading-tight">Repos</h3>
              <input
                type="text"
                placeholder="Search repo"
                class="my-2 w-full text-sm bg-gray-200 text-gray-500 rounded h-10 p-3 focus:outline-none"
                v-model="filter"
              />
              <div class="w-full max-h-96 overflow-auto		">
                <ul class="list text-sm font-medium text-gray-900 bg-white rounded-lg border border-gray-200 dark:bg-gray-700 dark:border-gray-600 dark:text-white">

                  <li
                    v-for="repo in publicRepos"
                    :key="repo.id"
                    class="flex py-2 px-4 w-full rounded-t-lg border-b border-gray-200 dark:border-gray-600"
                  >

                    <div class="w-4/5 h-12 py-4 px-1">
                      <p class="hover:text-blue-dark"> {{ repo.name }} </p>
                    </div>
                    <div class="w-1/5 h-10 text-right p-4">
                      <p class="text-sm">
                        <a :href="repo.url"></a>
                      </p>
                      <NuxtLink :to="{
                        name: 'repos-id',
                        params: {
                          id: repo.name
                        }
                      }">

                        <button class="text-white bg-indigo-500 border-0 py-1 px-4 rounded focus:outline-none hover:bg-indigo-600">

                          Setup

                        </button>
                      </NuxtLink>

                    </div>
                  </li>
                </ul>
              </div>
            </div>
          </div>
          <div class="sm:flex bg-gray-light sm:items-center px-2 py-4">
            <div class="flex-grow text-right">
              <button class="text-gray-500 hover:text-gray-700 py-2 px-4 rounded">
                Cancel
              </button>

            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  middleware: 'authenticated',
  data() {
    return {
      repos: null,
      filter: null,
    }
  },

  computed: {
    publicRepos() {
      return ((this.repos as any) || []).filter((repo: any) => {
        return this.filter
          ? repo.private === false &&
              repo.name.toLowerCase().includes(this?.filter.toLowerCase())
          : repo.private === false
      })
    },
    user() {
      return this.$auth.user
    },
  },

  mounted() {
    this.getRepos()
  },
  methods: {
    getRepos() {
      this.$axios
        .get(`${process.env.gitHubUrl}users/${this?.user?.login}/repos`)
        .then((response) => {
          this.repos = response.data
        })
    },
  },
})
</script>