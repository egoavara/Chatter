declare module "*.vue" {
  import { defineComponent } from "vue";
  // import {  } from "vuetify/components";
  const Component: ReturnType<typeof defineComponent>;
  export default Component;
}
// declare module "vuetify/lib"
// declare module "vuetify/components"
// declare module "vuetify" {

//   import {  } from "vuetify";
//   export default Component;
// }
