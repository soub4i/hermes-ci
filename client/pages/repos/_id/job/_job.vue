
<template>
  <div v-if="job">
    <div class="page-header">

      <h1 class="mb-2 pb-2 border-b border-white">{{ job.name }}</h1>

      <div
        v-for="l in logs"
        :key="l.id"
        class="w-full flex text-xs hover:bg-grey-darkest py-2"
      >
        <div class="flex-no-shrink w-12 text-yellow-light mr-2">INFO</div>
        <div class="flex-no-shrink w-56 mr-2">{{ new Date(l.timestamp * 1000)   }} </div>
        <div class="flex-1 break-words-all">{{ l.data }} </div>
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
      job: null,
    }
  },
  computed: {
    id() {
      return this.$route.params.job
    },
    user() {
      return this.$auth.user
    },
    logs() {
      return (this.job.logs || [])
        .map((l: string) => {
          try {
            return JSON.parse(l)
          } catch (error) {
            return null
          }
        })
        .filter((l: string | null) => l)
    },
  },
  mounted() {
    this.getJob()
  },
  methods: {
    async getJob() {
      try {
        const { data } = await this.$axios.get(
          `http://127.0.0.1:8080/jobs/${this.id}`
        )
        this.job = data
      } catch (error) {
        console.log(error)
      }
    },
  },
})
</script>