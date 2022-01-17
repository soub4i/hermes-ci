<template>
  <div>

    <div
      v-if="repo"
      class="sm:flex sm:items-center px-2 py-4"
    >
      <div class="flex-grow">

        <div class="w-full">

          <div class="text-center">

            <h1 class="px-2 py-3 leading my-4"> {{ repo.name }} </h1>
            <p class="text-gray-500 text-sm ">
              {{ repo.description }}
            </p>

            <div v-if="!hasHermesHook">
              <button
                class="text-white bg-green-500 py-2 px-4 rounded"
                @click.prevent="createHook"
              >

                Activate Hermes

              </button>
            </div>
            <div v-else>
              <button
                disabled
                class="text-white bg-green-500 py-2 px-2 rounded"
              >

                Already configured

              </button>

            </div>

          </div>
        </div>
      </div>
    </div>

    <div class="overflow-x-auto">
      <div class="min-w-screen min-h-screen bg-gray-100 flex justify-center bg-gray-100 font-sans overflow-hidden">
        <div class="w-full lg:w-5/6">
          <div class="bg-white shadow-md rounded  mt-6">
            <table class="min-w-max w-full table-auto">
              <thead>
                <tr class="bg-gray-200 text-gray-600 uppercase text-sm leading-normal">
                  <th class="py-3 px-6 text-left">Workflow</th>
                  <th class="py-3 px-6 text-left">Commit</th>
                  <th class="py-3 px-6 text-left">Sender</th>
                  <th class="py-3 px-6 text-center">Status</th>
                  <th class="py-3 px-6 text-center">Actions</th>
                </tr>
              </thead>
              <tbody class="text-gray-600 text-sm font-light">
                <tr
                  v-for="row in jobs"
                  :key="row.id"
                  class="border-b border-gray-200 bg-gray-50 hover:bg-gray-100"
                >
                  <td class="py-3 px-6 text-left">
                    <div class="flex items-center">

                      <span class="font-medium">{{ row.workflow }}</span>
                    </div>
                  </td>
                  <td class="py-3 px-6 text-left">
                    <div class="flex items-center">

                      <span class="font-medium">{{ row.name }}</span>
                    </div>
                  </td>
                  <td class="py-3 px-6 text-left">
                    <div class="flex items-center">
                      <div class="mr-2">
                        <img
                          class="w-6 h-6 rounded-full"
                          :src="row.owner.avatar"
                        />
                      </div>
                      <span>{{ row.owner.name }}</span>
                    </div>
                  </td>
                  <td class="py-3 px-6 text-center">
                    <span
                      :class="row.finished ? 'bg-green-200 text-green-600' : 'bg-red-200 text-red-600' "
                      class=" py-1 px-3 rounded-full text-xs"
                    >{{ row.finished ? "Finished" : "Not finished" }}</span>
                  </td>
                  <td class="py-3 px-6 text-center">
                    <div class="flex item-center justify-center">
                      <nuxt-link
                        :to="{
                          name: 'repos-id-job-job',
                          params: {
                            id: id,
                            job: row._id
                          }
                        }"
                        class="text-blue-500 hover:text-blue-700"
                      >
                        <div class="w-4 mr-2 transform hover:text-purple-500 hover:scale-110">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                            />
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                            />
                          </svg>
                        </div>
                      </nuxt-link>

                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="false"
      class="w-full p-6"
    >
      <div class="relative rounded overflow-hidden border  mb-8 bg-white">
        <div
          v-if="repo"
          class=" p-4 py-8"
        >
          <div
            v-if="false"
            class="bg-white mx-auto max-w-lg shadow-lg rounded-lg overflow-hidden"
          >
            <div class="sm:flex sm:items-center px-2 py-4">
              <div class="flex-grow">

                <div class="w-full">

                  <div class="text-center">

                    <h1 class="px-2 py-3 leading my-4"> {{ repo.name }} </h1>
                    <p class="text-gray-500 text-sm ">
                      {{ repo.description }}
                    </p>

                    <div v-if="!hasHermesHook">
                      <button
                        class="text-white bg-green-500 py-2 px-4 rounded"
                        @click.prevent="createHook"
                      >

                        Activate Hermes

                      </button>
                    </div>
                    <div v-else>
                      <button
                        disabled
                        class="text-white bg-green-500 py-2 px-2 rounded"
                      >

                        Already configured

                      </button>

                      <button
                        class="text-white bg-red-500 py-2 px-2 rounded"
                        @click.prevent="getLogs"
                      >

                        Check for logs
                      </button>

                    </div>

                  </div>
                </div>
              </div>
            </div>

          </div>
          <div class=" bg-gray-300 mx-auto my-4 max-w-lg shadow-lg rounded-lg overflow-hidden">
            <div class="sm:flex sm:items-center px-2 py-4">
              <div class="flex-grow">

                <div class="w-full">

                  <div>

                    {{ jobs }}

                  </div>

                </div>
              </div>
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
  data() {
    return {
      hooks: [],
      jobs: [],
      repo: null,
    }
  },
  computed: {
    id() {
      return this.$route.params.id
    },
    user() {
      return this.$auth.user
    },
    hasHermesHook() {
      return !!this.hooks.find((h: any) => h?.config?.digest === 'Hermes')
    },
  },
  async mounted() {
    this.getRepo()
    await this.getHooks()
    this.getJobs()
  },
  methods: {
    async getHooks() {
      try {
        const { data } = await this.$axios.get(
          `${this.process.env.gitHubUrl}repos/${this?.user?.login}/${this?.id}/hooks`
        )
        this.hooks = data
      } catch (error) {
        console.log(error)
      }
    },
    async getRepo() {
      try {
        const { data } = await this.$axios.get(
          `${this.process.env.gitHubUrl}repos/${this?.user?.login}/${this?.id}`
        )
        this.repo = data
      } catch (error) {
        console.log(error)
      }
    },
    createHook() {
      try {
        this.$axios.post(
          `${this.process.env.gitHubUrl}repos/${this?.user?.login}/${this?.id}/hooks`,
          {
            name: 'web',
            active: true,
            events: ['push', 'pull_request'],
            config: {
              url: `http://hermes.soubai.me/github/${this?.repo?.id}`,
              content_type: 'json',
              insecure_ssl: '0',
              digest: 'Hermes',
            },
          },
          {
            headers: {
              Accept: 'application/vnd.github.v3+json',
            },
          }
        )
      } catch (error) {
        console.log(error)
      }
    },
    async getJobs() {
      const { data } = await this.$axios.get(
        `${this.process.env.baseUrl}/github/${this?.repo?.id}`,
        {}
      )
      this.jobs = data || []
    },
  },
})
</script>