<template>
  <v-app>
    <v-main>
      <HelloWorld />
      <v-btn @click="onClick">{{ count }}</v-btn>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import HelloWorld from "./components/HelloWorld.vue";
// import { Vue, setup, Options } from "vue-class-component";
import { defineComponent, reactive, toRefs } from "vue";

// import {} from ''

export default defineComponent({
  name: "hello",
  components: {
    HelloWorld,
  },
  setup() {
    const state = reactive({
      count: 0,
    });
    return {
      ...toRefs(state),
      onClick : () =>{
        state.count += 1
      }
    }
  },
});
//  {
// handler() {}
// cnt = setup(() => {
//   const count = ref(0);
//   const increase = () => {
//     console.log(count);
//     count.value += 1;
//   };
//   return {
//     count,
//     increase,
//   };
// });
</script>
